package server

import (
	"context"
	"log/slog"

	"card-watcher/internal/models"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) SaveWatch(ctx context.Context, in *models.SaveWatchRequest) (*models.SaveWatchResponse, error) {
	s.logger.Info("Received a SaveWatch request")
	s.logger.Debug("request received", slog.Any("request", in))
	watchID, err := s.service.SaveWatch(ctx, int(in.ExpansionId), int(in.BlueprintId), in.Condition, in.Foil)
	if err != nil {
		s.logger.Error("error in save watch", slog.Any("error", err))
		return nil, err
	}
	return &models.SaveWatchResponse{
		WatchId: watchID,
	}, nil
}

func (s *server) ListWatches(ctx context.Context, in *emptypb.Empty) (*models.ListWatchesResponse, error) {
	s.logger.Info("Received a ListWatches request")
	response, err := s.service.ListWatches(ctx)
	if err != nil {
		s.logger.Error("error in list watches", slog.Any("error", err))
		return nil, err
	}
	return &response, nil
}

func (s *server) DeleteWatchById(ctx context.Context, in *models.DeleteWatchByIdRequest) (*emptypb.Empty, error) {
	s.logger.Info("Received a DeleteWatchById request")
	s.logger.Debug("request received", slog.Any("request", in))
	return &emptypb.Empty{}, s.service.DeleteWatchByID(ctx, in.WatchId)
}
