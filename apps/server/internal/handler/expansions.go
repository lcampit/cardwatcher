package handler

import (
	"context"

	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

func (s *handler) ListExpansions(ctx context.Context, in *apiv1.ListExpansionsRequest) (*apiv1.ListExpansionsResponse, error) {
	response, err := s.service.ListExpansions(ctx, in.Game, in.Name, in.Code)
	if err != nil {
		return nil, err
	}
	return response, nil
}
