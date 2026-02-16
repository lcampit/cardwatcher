package handler

import (
	"context"
	"log/slog"

	api "github.com/lcampit/card-watcher-server/internal/api/v1"
)

func (s *handler) ListBlueprints(ctx context.Context, in *api.ListBlueprintsRequest) (*api.ListBlueprintsResponse, error) {
	s.logger.Info("received a ListExpansions request",
		slog.Uint64("expansionId", in.ExpansionId),
		slog.String("name", in.Name))
	response, err := s.service.ListBlueprints(ctx, in.ExpansionId, in.Name)
	if err != nil {
		s.logger.Error("error in list blueprints", slog.Any("error", err))
		return nil, err
	}
	return response, nil
}
