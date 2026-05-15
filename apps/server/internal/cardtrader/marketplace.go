package cardtrader

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lcampit/cardwatcher/apps/server/internal/mongo"
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
	Description string            `json:"description"`
	Properties  ProductProperties `json:"properties_hash"`
	Expansion   struct {
		ID   uint64 `json:"id"`
		Code string `json:"code"`
		Name string `json:"name_en"`
	} `json:"expansion"`
	User       ProductUserInfo `json:"user"`
	Graded     bool            `json:"graded"`
	OnVacation bool            `json:"on_vavation"`
	BundleSize uint64          `json:"bundle_size"`
}

type ProductUserInfo struct {
	ID                   uint64 `json:"id"`
	Username             string `json:"username"`
	SellsViaHub          bool   `json:"can_sell_via_hub"`
	CountryCode          string `json:"country_code"`
	UserType             string `json:"user_type"`
	MaxSellableIn24Hours uint64 `json:"max_sellable_in24h_quantity"`
}

type ProductProperties struct {
	Condition mongo.WatchCondition `json:"condition"`
	Signed    bool                 `json:"signed"`
	Foil      bool                 `json:"foil"`
	Language  mongo.WatchLanguage  `json:"mtg_language"`
	Altered   bool                 `json:"altered"`
}

func (a *cardtraderAdapter) GetProducts(ctx context.Context, blueprintID uint64, foil bool) ([]Product, error) {
	response := map[string][]Product{}

	// Response to this endpoint is a map blueprintID -> list of products
	// if called with a blueprintID, the map contains only one entry with the
	// cheapest 25 products for that blueprint
	// if called with an expansionID, the map contains one entry for each blueprint
	// of that expansion, each of them with its own 25 cheapest products available
	endpoint := fmt.Sprintf("%s/%s", "marketplace", "products")
	blueprintIDString := strconv.FormatUint(blueprintID, 10)
	foilString := strconv.FormatBool(foil)
	_, err := a.client.R().
		SetQueryParams(map[string]string{
			"language":     "en",
			"blueprint_id": blueprintIDString,
			"foil":         foilString,
		}).
		SetResult(&response).
		Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("cardtrader get products: %w", err)
	}

	productsList, ok := response[blueprintIDString]
	if !ok {
		return nil, fmt.Errorf("cardtrader get products: response map does not contain blueprint ID %d", blueprintID)
	}
	return productsList, nil
}
