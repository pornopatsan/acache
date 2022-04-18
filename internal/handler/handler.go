package handler

import (
	"context"
	"log"

	"acache/pkg/util/errors"

	pb "acache/pkg/gen/proto"
)

func New(cache CacheRecorder) *Handler {
	return &Handler{cache: cache}
}

type Handler struct {
	pb.UnimplementedACacheServer
	cache CacheRecorder
}

func (h *Handler) Set(_ context.Context, item *pb.Item) (*pb.Response, error) {
	log.Printf("Set %v\n", item.Key)

	if err := h.cache.Set(item.Key, item.Value); err != nil {
		return &pb.Response{Status: pb.Status_UNKNOWN_ERROR}, errors.HandleError(err)
	}

	return &pb.Response{Status: pb.Status_OK}, nil
}

func (h *Handler) Get(_ context.Context, key *pb.Key) (*pb.ItemResponse, error) {
	log.Printf("Get %v\n", key.Key)

	k, v, err := h.cache.Get(key.Key)
	if err != nil {
		return &pb.ItemResponse{Status: pb.Status_KEY_NOT_FOUND}, errors.HandleError(err)
	}

	return &pb.ItemResponse{
		Item: &pb.Item{
			Key:   k,
			Value: v,
		},
		Status: pb.Status_OK,
	}, nil
}

func (h *Handler) Delete(_ context.Context, key *pb.Key) (*pb.Response, error) {
	log.Printf("Delete %v\n", key.Key)

	if err := h.cache.Delete(key.Key); err != nil {
		return &pb.Response{Status: pb.Status_KEY_NOT_FOUND}, errors.HandleError(err)
	}

	return &pb.Response{Status: pb.Status_OK}, nil
}

func (h *Handler) Size(_ context.Context, _ *pb.Empty) (*pb.SizeResponse, error) {
	return &pb.SizeResponse{Size: h.cache.Size()}, nil
}
