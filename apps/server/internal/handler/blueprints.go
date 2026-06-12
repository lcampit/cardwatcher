package handler

import (
	"context"

	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

func (s *handler) ListBlueprints(ctx context.Context, in *apiv1.ListBlueprintsRequest) (*apiv1.ListBlueprintsResponse, error) {
	response, err := s.service.ListBlueprints(ctx, in.ExpansionId, in.Name)
	if err != nil {
		return nil, err
	}
	return response, nil
}
