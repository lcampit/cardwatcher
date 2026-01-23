package server

import (
	"context"
	"log/slog"

	"card-watcher/internal/models"

	"github.com/rs/zerolog/log"
)

func (s *server) ListBlueprints(ctx context.Context, in *models.ListBlueprintsRequest) (*models.ListBlueprintsResponse, error) {
	s.logger.Info("Received a ListExpansions request",
		slog.Int("expansionID", int(in.ExpansionId)),
		slog.String("name", in.Name))
	response, err := s.service.ListBlueprints(ctx, int(in.ExpansionId), in.Name)
	if err != nil {
		log.Error().Err(err).Msg("error in list expansions")
		return nil, err
	}
	return &response, nil
}
