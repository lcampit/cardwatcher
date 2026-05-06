package cardtrader

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
)

type blueprint struct {
	ID                 uint64               `json:"id"`
	Name               string               `json:"name"`
	Version            string               `json:"version"`
	GameID             uint64               `json:"game_id"`
	CategoryID         uint64               `json:"category_id"`
	ExpansionID        uint64               `json:"expansion_id"`
	ImageURL           string               `json:"image_url"`
	EditableProperties []EditableProperties `json:"editable_properties"`
	ScryfallID         string               `json:"scryfall_id"`
	CardMarketIDs      []uint64             `json:"card_market_ids"`
	TcgPlayerID        uint64               `json:"tcg_player_id"`
}

func (a *cardtraderAdapter) GetBlueprints(ctx context.Context, expansionID uint64) ([]*blueprint, error) {
	var response []*blueprint
	endpoint := fmt.Sprintf("%s/%s", "blueprints", "export")
	expansionIDString := strconv.FormatUint(expansionID, 10)
	err := a.client.Path(endpoint).
		Param("expansion_id", expansionIDString).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("cardtrader get blueprint request for expansion id %d: %w", expansionID, err)
	}
	a.logger.Debug("received blueprints for expansion id",
		slog.Int("blueprintCount", len(response)),
		slog.Uint64("expansionId", expansionID))
	return response, nil
}

func (a *cardtraderAdapter) GetBlueprintNameByExpansionID(ctx context.Context, expansionID, blueprintID uint64) (string, error) {
	var response []blueprint
	endpoint := fmt.Sprintf("%s/%s", "blueprints", "export")
	expansionIDString := strconv.FormatUint(expansionID, 10)
	err := a.client.Path(endpoint).
		Param("expansion_id", expansionIDString).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return "", fmt.Errorf("cardtrader get blueprint endpoint for expansion id %d, blueprint id %d: %w", expansionID, blueprintID, err)
	}
	a.logger.Debug("received blueprints for expansion id",
		slog.Int("blueprintCount", len(response)),
		slog.Uint64("expansionId", expansionID))

	for _, blueprint := range response {
		if blueprint.ID == blueprintID {
			a.logger.Debug("found blueprint name",
				slog.Uint64("expansionId", expansionID),
				slog.Uint64("blueprintId", blueprintID),
				slog.String("blueprintName", blueprint.Name))
			return blueprint.Name, nil
		}
	}
	return "", fmt.Errorf("no name found for expansion id %d and blueprint id %d", expansionID, blueprintID)
}
