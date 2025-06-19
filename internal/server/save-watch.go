package server

import (
	"card-watcher/internal/models"
	"context"

	"github.com/rs/zerolog/log"
)

func (s *server) SaveWatch(ctx context.Context, in *models.SaveWatchRequest) (*models.SaveWatchResponse, error) {
	log.Info().Msg("Received a SaveWatch request")
	log.Debug().Interface("request", in).Msg("")
	watchId, err := s.service.SaveWatch(ctx, in.AccessToken, int(in.ExpansionId), int(in.BlueprintId), in.Condition, in.Foil)
	if err != nil {
		log.Error().Err(err).Msg("error in save watch")
		return nil, err
	}
	return &models.SaveWatchResponse{
		WatchId: watchId,
	}, nil
}
