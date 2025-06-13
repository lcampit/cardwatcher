package cardtrader

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
)

type Expansion struct {
	Id     int    `json:"id"`
	GameId int    `json:"game_id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
}

func (a *cardtraderAdapter) GetExpansions(ctx context.Context) (*[]Expansion, error) {
	var response []Expansion
	endpoint := fmt.Sprintf("%s/%s", a.baseUrl, "expansions")
	err := requests.URL(endpoint).Bearer(a.accessToken).ToJSON(&response).Fetch(ctx)
	return &response, err
}
