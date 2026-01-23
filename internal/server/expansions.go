package server

import (
	"context"
	"log/slog"

	"card-watcher/internal/models"

	"github.com/rs/zerolog/log"
)

func (s *server) ListExpansions(ctx context.Context, in *models.ListExpansionsRequest) (*models.ListExpansionsResponse, error) {
	s.logger.Info("Received a ListExpansions request",
		slog.String("name", in.Name),
		slog.String("code", in.Code))
	response, err := s.service.ListExpansions(ctx, in.Name, in.Code)
	if err != nil {
		log.Error().Err(err).Msg("error in list expansions")
		return nil, err
	}
	return &response, nil
}
