package server

import (
	"card-watcher/internal/models"
	"context"

	"github.com/rs/zerolog/log"
)

func (s *server) ListBlueprints(ctx context.Context, in *models.ListBlueprintsRequest) (*models.ListBlueprintsResponse, error) {
	log.Info().Msgf("Received a ListExpansions request for expansion id '%d' and blueprint name '%s'", in.ExpansionId, in.Name)
	response, err := s.service.ListBlueprints(ctx, int(in.ExpansionId), in.Name)
	if err != nil {
		log.Error().Err(err).Msg("error in list expansions")
		return nil, err
	}
	return &response, nil
}
