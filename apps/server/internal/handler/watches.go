package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

func (s *handler) CreateWatch(ctx context.Context, in *apiv1.CreateWatchRequest) (*apiv1.CreateWatchResponse, error) {
	watchID, err := s.service.CreateWatch(ctx, in)
	if err != nil {
		return nil, err
	}
	return &apiv1.CreateWatchResponse{
		WatchId: watchID,
	}, nil
}

func (s *handler) ListWatches(ctx context.Context, in *emptypb.Empty) (*apiv1.ListWatchesResponse, error) {
	response, err := s.service.ListWatches(ctx)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *handler) DeleteWatchByID(ctx context.Context, in *apiv1.DeleteWatchByIDRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.service.DeleteWatchByID(ctx, in.WatchId)
}
