package cardtrader

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
)

type InfoResponse struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	SharedSecret string `json:"shared_secret"`
}

func (a *cardtraderAdapter) Info(ctx context.Context) (*InfoResponse, error) {
	var response InfoResponse
	endpoint := fmt.Sprintf("%s/%s", a.baseUrl, "info")
	err := requests.URL(endpoint).Bearer(a.accessToken).ToJSON(&response).Fetch(ctx)
	return &response, err
}
