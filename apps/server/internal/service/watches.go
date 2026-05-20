package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lcampit/cardwatcher/apps/server/internal/cardtrader"
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

	s.logger.Debug("returning watches", slog.Int("watchCount", len(result)))
	return &apiv1.ListWatchesResponse{
		Watches: result,
	}, nil
}

func (s *service) SaveWatch(ctx context.Context, expansionID, blueprintID uint64, condition apiv1.Condition, language apiv1.Language, foil bool) (string, error) {
	blueprintName, err := s.cardtraderAdapter.GetBlueprintNameByExpansionID(ctx, expansionID, blueprintID)
	if err != nil {
		return "", fmt.Errorf("finding name for expansion %d and blueprint %d: %w", expansionID, blueprintID, err)
	}
	expansionName, err := s.cardtraderAdapter.GetExpansionNameByID(ctx, expansionID)
	if err != nil {
		return "", fmt.Errorf("finding name for expansion %d: %w", expansionID, err)
	}
	newWatchID, err := s.mongoAdapter.SaveWatch(ctx, &mongo.Watch{
		Name:          blueprintName,
		ExpansionID:   expansionID,
		ExpansionName: expansionName,
		BlueprintID:   blueprintID,
		Condition:     convertModelConditionToEntityCondition(condition),
		Language:      convertModelLanguageToEntityLanguage(language),
		Foil:          foil,
	})
	if err != nil {
		return "", err
	}
	return newWatchID, nil
}

func (s *service) CreateWatch(ctx context.Context, request *apiv1.CreateWatchRequest) (string, error) {
	nameOrCodeRequested := normalizeString(request.ExpansionNameOrCode)
	savedExpansion, ok := s.expansionCodeMap.Load(nameOrCodeRequested)

	var expansion *cardtrader.Expansion
	if ok {
		expansion, _ = savedExpansion.(*cardtrader.Expansion)
	} else {
		savedExpansion, ok = s.expansionNameMap.Load(nameOrCodeRequested)
		if ok {
			expansion, _ = savedExpansion.(*cardtrader.Expansion)
		} else {
			return "", fmt.Errorf("finding expansion in preloaded map: no expansion found for %s", request.ExpansionNameOrCode)
		}
	}

	blueprints, err := s.cardtraderAdapter.GetBlueprints(ctx, expansion.ID)
	if err != nil {
		return "", fmt.Errorf("finding blueprints for expansion %s (ID %d)", request.ExpansionNameOrCode, expansion.ID)
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
				return "", fmt.Errorf("creating watch: %v", err)
			}
			return watchID, nil
		}
	}

	return "", fmt.Errorf("no blueprint named %s found for expansion %s (ID %d)", cardNameRequested, expansion.Name, expansion.ID)
}

func (s *service) DeleteWatchByID(ctx context.Context, watchID string) error {
	return s.mongoAdapter.DeleteWatchByID(ctx, watchID)
}
