package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lcampit/cardwatcher/apps/server/internal/errors"
	"github.com/lcampit/cardwatcher/apps/server/internal/mongo"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

func (s *service) ListWatches(ctx context.Context) (*apiv1.ListWatchesResponse, error) {
	watches, err := s.mongoAdapter.GetWatches(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting watches from mongo adapter: %w", err)
	}

	var result []*apiv1.Watch

	for _, entity := range watches {
		result = append(result, convertEntityWatchToModelWatch(entity))
	}

	s.logger.Debug("returning watches", slog.Int("watch_count", len(result)))
	return &apiv1.ListWatchesResponse{
		Watches: result,
	}, nil
}

func (s *service) CreateWatch(ctx context.Context, request *apiv1.CreateWatchRequest) (string, error) {
	nameOrCodeRequested := normalizeString(request.ExpansionNameOrCode)
	expansion, ok := s.getExpansionFromMaps(nameOrCodeRequested)
	if !ok {
		expansions, err := s.cardtraderAdapter.GetExpansions(ctx)
		if err != nil {
			return "", errors.NewInternal("cardtrader expansion endpoint failed", "", err)
		}
		for _, exp := range expansions {
			if exp.GetNormalizedName() == nameOrCodeRequested || exp.GetNormalizedCode() == nameOrCodeRequested {
				expansion = exp
				break
			}
		}
		if expansion == nil {
			return "", errors.NewNotFound("expansion", request.ExpansionNameOrCode, "")
		}
	}

	blueprints, err := s.cardtraderAdapter.GetBlueprints(ctx, expansion.ID)
	if err != nil {
		return "", errors.NewInternal("cardtrader blueprint endpoint failed", "", err)
	}

	condition := mongo.WatchConditionAny
	if request.Condition != apiv1.Condition_CONDITION_UNSPECIFIED {
		condition = convertModelConditionToEntityCondition(request.Condition)
	}

	language := mongo.WatchLanguageAny
	if request.Language != apiv1.Language_LANGUAGE_UNSPECIFIED {
		language = convertModelLanguageToEntityLanguage(request.Language)
	}

	cardNameRequested := normalizeString(request.CardName)
	for _, blueprint := range blueprints {
		if blueprint.ExpansionID == expansion.ID && cardNameRequested == normalizeString(blueprint.Name) {
			watchID, err := s.mongoAdapter.SaveWatch(ctx, &mongo.Watch{
				Name:          blueprint.Name,
				ExpansionID:   expansion.ID,
				ExpansionName: expansion.Name,
				BlueprintID:   blueprint.ID,
				Condition:     condition,
				Language:      language,
				Foil:          request.Foil,
			})
			if err != nil {
				return "", errors.NewInternal("creating watch failed", "", err)
			}
			return watchID, nil
		}
	}

	return "", errors.NewNotFound("blueprint", cardNameRequested, "")
}

func (s *service) DeleteWatchByID(ctx context.Context, watchID string) error {
	return s.mongoAdapter.DeleteWatchByID(ctx, watchID)
}
