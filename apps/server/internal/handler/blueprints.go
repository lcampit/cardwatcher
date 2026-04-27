package handler

import (
	"context"
	"log/slog"

	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

func (s *handler) ListBlueprints(ctx context.Context, in *apiv1.ListBlueprintsRequest) (*apiv1.ListBlueprintsResponse, error) {
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
