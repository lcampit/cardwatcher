package handler

import (
	"context"
	"log/slog"

	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *handler) SaveWatch(ctx context.Context, in *apiv1.SaveWatchRequest) (*apiv1.SaveWatchResponse, error) {
	s.logger.Info("received a SaveWatch request")
	s.logger.Debug("request received", slog.Any("request", in))
	watchID, err := s.service.SaveWatch(ctx, in.ExpansionId, in.BlueprintId, in.Condition, in.Foil)
	if err != nil {
		s.logger.Error("error in save watch", slog.Any("error", err))
		return nil, err
	}
	return &apiv1.SaveWatchResponse{
		WatchId: watchID,
	}, nil
}

func (s *handler) ListWatches(ctx context.Context, in *emptypb.Empty) (*apiv1.ListWatchesResponse, error) {
	s.logger.Info("received a ListWatches request")
	response, err := s.service.ListWatches(ctx)
	if err != nil {
		s.logger.Error("error in list watches", slog.Any("error", err))
		return nil, err
	}
	return response, nil
}

func (s *handler) DeleteWatchByID(ctx context.Context, in *apiv1.DeleteWatchByIDRequest) (*emptypb.Empty, error) {
	s.logger.Info("received a DeleteWatchById request")
	s.logger.Debug("request received", slog.Any("request", in))
	return &emptypb.Empty{}, s.service.DeleteWatchByID(ctx, in.WatchId)
}
