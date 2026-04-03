package server

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/pornopatsan/acache/src/api"
	"github.com/pornopatsan/acache/src/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func startTestServer(t *testing.T, capacity uint) api.ACacheClient {
	t.Helper()

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	srv := Create(storage.Create(capacity))
	api.RegisterACacheServer(s, srv)

	go func() {
		if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			t.Logf("server exited: %v", err)
		}
	}()
	t.Cleanup(func() { s.Stop() })

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return api.NewACacheClient(conn)
}

func valueForKey(key string) []byte {
	h := sha256.Sum256([]byte(key))
	return h[:]
}

func TestServerDataIntegrity(t *testing.T) {
	client := startTestServer(t, 200)
	ctx := context.Background()

	// Save 100 items with value = sha256(key)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("item-%d", i)
		resp, err := client.Save(ctx, &api.Item{Key: key, Value: valueForKey(key)})
		require.NoError(t, err)
		assert.Equal(t, api.Status_OK, resp.Status)
	}

	// Get all back and verify values match
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("item-%d", i)
		resp, err := client.Get(ctx, &api.Key{Key: key})
		require.NoError(t, err)
		assert.Equal(t, api.Status_OK, resp.Status)
		assert.Equal(t, valueForKey(key), resp.Item.Value, "value mismatch for key %s", key)
	}
}

func TestServerConcurrentIntegrity(t *testing.T) {
	const goroutines = 16
	const opsPerGoroutine = 200
	const capacity = uint(goroutines * opsPerGoroutine) // large enough to hold all

	client := startTestServer(t, capacity)
	ctx := context.Background()

	var wg sync.WaitGroup

	// Each goroutine saves its own key range with value = sha256(key)
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < opsPerGoroutine; i++ {
				key := fmt.Sprintf("g%d-k%d", gid, i)
				_, err := client.Save(ctx, &api.Item{Key: key, Value: valueForKey(key)})
				if err != nil {
					t.Errorf("Save failed for %s: %v", key, err)
					return
				}
			}
		}(g)
	}
	wg.Wait()

	// Verify all values
	for g := 0; g < goroutines; g++ {
		for i := 0; i < opsPerGoroutine; i++ {
			key := fmt.Sprintf("g%d-k%d", g, i)
			resp, err := client.Get(ctx, &api.Key{Key: key})
			require.NoError(t, err)
			assert.Equal(t, api.Status_OK, resp.Status, "key %s not found", key)
			assert.Equal(t, valueForKey(key), resp.Item.Value, "value mismatch for key %s", key)
		}
	}
}

func TestServerCapacityEnforcement(t *testing.T) {
	const capacity = 32
	client := startTestServer(t, capacity)
	ctx := context.Background()

	// Save 64 items sequentially
	for i := 0; i < 64; i++ {
		key := fmt.Sprintf("cap-%d", i)
		_, err := client.Save(ctx, &api.Item{Key: key, Value: valueForKey(key)})
		require.NoError(t, err)
	}

	// Count how many are still present
	found := 0
	foundKeys := make([]int, 0)
	for i := 0; i < 64; i++ {
		key := fmt.Sprintf("cap-%d", i)
		resp, err := client.Get(ctx, &api.Key{Key: key})
		require.NoError(t, err)
		if resp.Status == api.Status_OK {
			found++
			foundKeys = append(foundKeys, i)
			assert.Equal(t, valueForKey(key), resp.Item.Value)
		}
	}

	assert.Equal(t, capacity, found, "expected exactly %d items, got %d", capacity, found)

	// The surviving items should be the most recently saved (32..63)
	for _, idx := range foundKeys {
		assert.GreaterOrEqual(t, idx, 32, "item %d should have been evicted", idx)
	}
}

func TestServerConcurrentCapacity(t *testing.T) {
	const capacity = 64
	const goroutines = 32
	const keysPerGoroutine = 100

	client := startTestServer(t, capacity)
	ctx := context.Background()

	for attempt := 0; attempt < 5; attempt++ {
		var wg sync.WaitGroup
		for g := 0; g < goroutines; g++ {
			wg.Add(1)
			go func(gid int) {
				defer wg.Done()
				for i := 0; i < keysPerGoroutine; i++ {
					key := fmt.Sprintf("cc-a%d-g%d-k%d", attempt, g, i)
					_, err := client.Save(ctx, &api.Item{Key: key, Value: []byte{byte(gid), byte(i)}})
					if err != nil {
						t.Errorf("Save failed: %v", err)
						return
					}
				}
			}(g)
		}
		wg.Wait()

		// Count items by probing all keys
		found := 0
		for g := 0; g < goroutines; g++ {
			for i := 0; i < keysPerGoroutine; i++ {
				key := fmt.Sprintf("cc-a%d-g%d-k%d", attempt, g, i)
				resp, err := client.Get(ctx, &api.Key{Key: key})
				if err == nil && resp.Status == api.Status_OK {
					found++
				}
			}
		}
		assert.LessOrEqual(t, found, capacity,
			"attempt %d: found %d items, exceeds capacity %d", attempt, found, capacity)
	}
}

func TestServerRaceStress(t *testing.T) {
	const capacity = 128
	const goroutines = 100
	const duration = 2 * time.Second

	client := startTestServer(t, capacity)
	ctx := context.Background()

	var wg sync.WaitGroup
	stop := make(chan struct{})

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(int64(gid)))
			for {
				select {
				case <-stop:
					return
				default:
				}

				key := fmt.Sprintf("stress-%d", rng.Intn(256))
				op := rng.Intn(3)
				switch op {
				case 0:
					client.Save(ctx, &api.Item{Key: key, Value: valueForKey(key)})
				case 1:
					client.Get(ctx, &api.Key{Key: key})
				case 2:
					client.Remove(ctx, &api.Key{Key: key})
				}
			}
		}(g)
	}

	time.Sleep(duration)
	close(stop)
	wg.Wait()
	// If we get here without panic or race detector error, the test passes
}
