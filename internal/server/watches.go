package server

import (
	"card-watcher/internal/models"
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) SaveWatch(ctx context.Context, in *models.SaveWatchRequest) (*models.SaveWatchResponse, error) {
	log.Info().Msg("Received a SaveWatch request")
	log.Debug().Interface("request", in).Msg("")
	watchId, err := s.service.SaveWatch(ctx, int(in.ExpansionId), int(in.BlueprintId), in.Condition, in.Foil)
	if err != nil {
		log.Error().Err(err).Msg("error in save watch")
		return nil, err
	}
	return &models.SaveWatchResponse{
		WatchId: watchId,
	}, nil
}

func (s *server) ListWatches(ctx context.Context, in *emptypb.Empty) (*models.ListWatchesResponse, error) {
	log.Info().Msg("Received a ListWatches request")
	response, err := s.service.ListWatches(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error in list watches")
		return nil, err
	}
	return &response, nil
}
