package cardtrader

import (
	"card-watcher/internal/entities"
	"context"
	"fmt"
	"strconv"

	"github.com/carlmjohnson/requests"
)

type Product struct {
	Id          int    `json:"id"`
	BlueprintId int    `json:"blueprint_id"`
	Name        string `json:"name_en"`
	Quantity    int    `json:"quantity"`
	Price       struct {
		Cents    int    `json:"cents"`
		Currency string `json:"currency"`
	} `json:"price"`
	Description string `json:"description"`
	Properties  struct {
		Condition entities.WatchCondition `json:"condition"`
		Signed    bool                    `json:"signed"`
		Foil      bool                    `json:"foil"`
		Language  string                  `json:"mtg_language"`
		Altered   bool                    `json:"altered"`
	} `json:"properties_hash"`
	Expansion struct {
		Id   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name_en"`
	} `json:"expansion"`
	User struct {
		Id                   int    `json:"id"`
		Username             string `json:"username"`
		SellsViaHub          bool   `json:"can_sell_via_hub"`
		CountryCode          string `json:"country_code"`
		UserType             string `json:"user_type"`
		MaxSellableIn24Hours int    `json:"max_sellable_in24h_quantity"`
	} `json:"user"`
	Graded     bool `json:"graded"`
	OnVacation bool `json:"on_vavation"`
	BundleSize int  `json:"bundle_size"`
}

func (a *cardtraderAdapter) GetCurrentPricing(ctx context.Context, watch *entities.Watch) (int, error) {
	response := map[string][]Product{}

	endpoint := fmt.Sprintf("%s/%s/%s", a.baseUrl, "marketplace", "products")
	blueprintIdString := strconv.Itoa(watch.BlueprintId)
	foilString := strconv.FormatBool(watch.Foil)
	err := requests.URL(endpoint).Bearer(a.accessToken).
		Param("language", "en").Param("blueprint_id", blueprintIdString).Param("foil", foilString).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting products from adapter: %w", err)
	}

	if products, ok := response[blueprintIdString]; ok {
		for _, product := range products {
			if product.Properties.Condition == watch.Condition && product.User.SellsViaHub {
				return product.Price.Cents, nil
			}
		}
	}

	return 0, nil
}
