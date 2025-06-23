package server

import (
	"card-watcher/internal/models"
	"context"

	"github.com/rs/zerolog/log"
)

func (s *server) ListExpansions(ctx context.Context, in *models.ListExpansionsRequest) (*models.ListExpansionsResponse, error) {
	log.Info().Msgf("Received a ListExpansions request for name '%s' or code '%s'", in.Name, in.Code)
	response, err := s.service.ListExpansions(ctx, in.Name, in.Code)
	if err != nil {
		log.Error().Err(err).Msg("error in list expansions")
		return nil, err
	}
	return &response, nil
}
