package cardtrader

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lcampit/cardwatcher/internal/server/mongo"

	"github.com/carlmjohnson/requests"
)

type Product struct {
	ID          uint64 `json:"id"`
	BlueprintID uint64 `json:"blueprint_id"`
	Name        string `json:"name_en"`
	Quantity    uint64 `json:"quantity"`
	Price       struct {
		Cents    uint64 `json:"cents"`
		Currency string `json:"currency"`
	} `json:"price"`
	Description string `json:"description"`
	Properties  struct {
		Condition mongo.WatchCondition `json:"condition"`
		Signed    bool                 `json:"signed"`
		Foil      bool                 `json:"foil"`
		Language  string               `json:"mtg_language"`
		Altered   bool                 `json:"altered"`
	} `json:"properties_hash"`
	Expansion struct {
		ID   uint64 `json:"id"`
		Code string `json:"code"`
		Name string `json:"name_en"`
	} `json:"expansion"`
	User struct {
		ID                   uint64 `json:"id"`
		Username             string `json:"username"`
		SellsViaHub          bool   `json:"can_sell_via_hub"`
		CountryCode          string `json:"country_code"`
		UserType             string `json:"user_type"`
		MaxSellableIn24Hours uint64 `json:"max_sellable_in24h_quantity"`
	} `json:"user"`
	Graded     bool   `json:"graded"`
	OnVacation bool   `json:"on_vavation"`
	BundleSize uint64 `json:"bundle_size"`
}

func (a *cardtraderAdapter) GetCurrentPricingCents(ctx context.Context, watch *mongo.Watch) (uint64, error) {
	response := map[string][]Product{}

	endpoint := fmt.Sprintf("%s/%s/%s", a.baseURL, "marketplace", "products")
	blueprintIDString := strconv.FormatUint(watch.BlueprintID, 10)
	foilString := strconv.FormatBool(watch.Foil)
	err := requests.URL(endpoint).Bearer(a.accessToken).
		Param("language", "en").Param("blueprint_id", blueprintIDString).Param("foil", foilString).
		ToJSON(&response).Fetch(ctx)
	if err != nil {
		return 0, fmt.Errorf("cardtrader get products: %w", err)
	}

	if products, ok := response[blueprintIDString]; ok {
		for _, product := range products {
			if product.Properties.Condition == watch.Condition && product.User.SellsViaHub {
				return product.Price.Cents, nil
			}
		}
	}

	return 0, nil
}
