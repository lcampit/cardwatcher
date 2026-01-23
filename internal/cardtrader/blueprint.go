package cardtrader

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/carlmjohnson/requests"
)

type blueprint struct {
	ID                 int                  `json:"id"`
	Name               string               `json:"name"`
	Version            string               `json:"version"`
	GameID             int                  `json:"game_id"`
	CategoryID         int                  `json:"category_id"`
	ExpansionID        int                  `json:"expansion_id"`
	ImageURL           string               `json:"image_url"`
	EditableProperties []EditableProperties `json:"editable_properties"`
	ScryfallID         string               `json:"scryfall_id"`
	CardMarketIDs      []int                `json:"card_market_ids"`
	TcgPlayerID        int                  `json:"tcg_player_id"`
}

func (a *cardtraderAdapter) GetBlueprints(ctx context.Context, expansionID int) ([]*blueprint, error) {
	var response []*blueprint
	endpoint := fmt.Sprintf("%s/%s/%s", a.baseURL, "blueprints", "export")
	expansionIDString := strconv.Itoa(expansionID)
	err := requests.URL(endpoint).Bearer(a.accessToken).
		Param("expansion_id", expansionIDString).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in cardtrader blueprint endpoint for expansion id %d: %w", expansionID, err)
	}
	a.logger.Debug("received blueprints for expansion id", slog.Int("blueprintCount", len(response)), slog.Int("expansionID", expansionID))
	return response, nil
}

func (a *cardtraderAdapter) GetBlueprintNameByExpansionID(ctx context.Context, expansionID, blueprintID int) (string, error) {
	var response []blueprint
	endpoint := fmt.Sprintf("%s/%s/%s", a.baseURL, "blueprints", "export")
	expansionIDString := strconv.Itoa(expansionID)
	err := requests.URL(endpoint).Bearer(a.accessToken).
		Param("expansion_id", expansionIDString).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return "", fmt.Errorf("error in cardtrader blueprint endpoint for expansion id %d, blueprint id %d: %w", expansionID, blueprintID, err)
	}
	a.logger.Debug("received blueprints for expansion id", slog.Int("blueprintCount", len(response)), slog.Int("expansionID", expansionID))

	for _, blueprint := range response {
		if blueprint.ID == blueprintID {
			a.logger.Debug("found blueprint name",
				slog.Int("expansionID", expansionID),
				slog.Int("blueprintID", blueprintID),
				slog.String("blueprintName", blueprint.Name))
			return blueprint.Name, nil
		}
	}
	return "", fmt.Errorf("no name found for expansion id %d and blueprint id %d", expansionID, blueprintID)
}
