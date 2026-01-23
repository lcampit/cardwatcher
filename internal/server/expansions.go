package server

import (
	"context"
	"log/slog"

	"card-watcher/internal/models"
)

func (s *server) ListExpansions(ctx context.Context, in *models.ListExpansionsRequest) (*models.ListExpansionsResponse, error) {
	s.logger.Info("Received a ListExpansions request",
		slog.String("name", in.Name),
		slog.String("code", in.Code))
	response, err := s.service.ListExpansions(ctx, in.Name, in.Code)
	if err != nil {
		s.logger.Error("error in list expansions", slog.Any("error", err))
		return nil, err
	}
	return &response, nil
}
