package server

import (
	"context"
	"log/slog"

	"card-watcher/internal/models"
)

func (s *server) ListBlueprints(ctx context.Context, in *models.ListBlueprintsRequest) (*models.ListBlueprintsResponse, error) {
	s.logger.Info("received a ListExpansions request",
		slog.Int("expansionId", int(in.ExpansionId)),
		slog.String("name", in.Name))
	response, err := s.service.ListBlueprints(ctx, int(in.ExpansionId), in.Name)
	if err != nil {
		s.logger.Error("error in list blueprints", slog.Any("error", err))
		return nil, err
	}
	return response, nil
}
