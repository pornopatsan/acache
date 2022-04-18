package handler

type CacheRecorder interface {
	Set(key string, value []byte) error
	Get(string) (key string, value []byte, err error)
	Delete(key string) error
	Size() int64
}
