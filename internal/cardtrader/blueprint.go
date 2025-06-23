package cardtrader

import (
	"context"
	"fmt"
	"strconv"

	"github.com/carlmjohnson/requests"
	"github.com/rs/zerolog/log"
)

type blueprint struct {
	Id                 int                  `json:"id"`
	Name               string               `json:"name"`
	Version            string               `json:"version"`
	GameId             int                  `json:"game_id"`
	CategoryId         int                  `json:"category_id"`
	ExpansionId        int                  `json:"expansion_id"`
	ImageUrl           string               `json:"image_url"`
	EditableProperties []EditableProperties `json:"editable_properties"`
	ScryfallId         string               `json:"scryfall_id"`
	CardMarketIds      []int                `json:"card_market_ids"`
	TcgPlayerId        int                  `json:"tcg_player_id"`
}

func (a *cardtraderAdapter) GetBlueprints(ctx context.Context, expansionId int) ([]*blueprint, error) {
	var response []*blueprint
	endpoint := fmt.Sprintf("%s/%s/%s", a.baseUrl, "blueprints", "export")
	expansionIdString := strconv.Itoa(expansionId)
	err := requests.URL(endpoint).Bearer(a.accessToken).
		Param("expansion_id", expansionIdString).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in cardtrader blueprint endpoint for expansion id %d: %w", expansionId, err)
	}
	log.Debug().Msgf("received %d blueprints for expansion id %d", len(response), expansionId)
	return response, nil
}

func (a *cardtraderAdapter) GetBlueprintNameByExpansionId(ctx context.Context, expansionId, blueprintId int) (string, error) {
	var response []blueprint
	endpoint := fmt.Sprintf("%s/%s/%s", a.baseUrl, "blueprints", "export")
	expansionIdString := strconv.Itoa(expansionId)
	err := requests.URL(endpoint).Bearer(a.accessToken).
		Param("expansion_id", expansionIdString).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return "", fmt.Errorf("error in cardtrader blueprint endpoint for expansion id %d, blueprint id %d: %w", expansionId, blueprintId, err)
	}
	log.Debug().Msgf("received %d blueprints for expansion id %d", len(response), expansionId)

	for _, blueprint := range response {
		if blueprint.Id == blueprintId {
			log.Debug().Msgf("found name for expansion id %d blueprint id %d: %s", expansionId, blueprintId, blueprint.Name)
			return blueprint.Name, nil
		}
	}
	return "", fmt.Errorf("no name found for expansion id %d and blueprint id %d", expansionId, blueprintId)
}
