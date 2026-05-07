package service

import (
	"testing"

	"github.com/lcampit/cardwatcher/apps/server/internal/cardtrader"
	"github.com/lcampit/cardwatcher/apps/server/internal/mongo"
)

const (
	watchName     = "watch-name-test"
	expansionName = "expansion-name-test"
	expansionID   = 1
	blueprintID   = 1
)

var allConditions = []mongo.WatchCondition{
	mongo.WatchConditionAny, mongo.WatchConditionNM, mongo.WatchConditionMP,
	mongo.WatchConditionSP, mongo.WatchConditionPL, mongo.WatchConditionPO,
}

func TestWatchProductMatchesByCondition(t *testing.T) {
	watch := mongo.Watch{
		Name:          watchName,
		ExpansionID:   expansionID,
		ExpansionName: expansionName,
		BlueprintID:   blueprintID,
		Condition:     mongo.WatchConditionNM,
		Foil:          false,
	}

	product := cardtrader.Product{
		User: cardtrader.ProductUserInfo{
			SellsViaHub: true,
		},
		Properties: cardtrader.ProductProperties{
			Condition: mongo.WatchConditionNM,
		},
	}

	matchingProduct := watchConditionsMatchProduct(&watch, product)
	if !matchingProduct {
		t.Errorf("watch does not match product by condition")
	}
}

func TestWatchWithAnyConditionMatchesAnyProduct(t *testing.T) {
	watch := mongo.Watch{
		Name:          watchName,
		ExpansionID:   expansionID,
		ExpansionName: expansionName,
		BlueprintID:   blueprintID,
		Condition:     mongo.WatchConditionAny,
		Foil:          false,
	}

	for _, condition := range allConditions {
		product := cardtrader.Product{
			User: cardtrader.ProductUserInfo{
				SellsViaHub: true,
			},
			Properties: cardtrader.ProductProperties{
				Condition: condition,
			},
		}
		matchingProduct := watchConditionsMatchProduct(&watch, product)
		if !matchingProduct {
			t.Errorf("any condition watch does not match product with condition: %s", condition)
		}
	}
}

func TestWatchDoesNotMatchProductWithoutSellsViaHub(t *testing.T) {
	watch := mongo.Watch{
		Name:          watchName,
		ExpansionID:   expansionID,
		ExpansionName: expansionName,
		BlueprintID:   blueprintID,
		Condition:     mongo.WatchConditionAny,
		Foil:          false,
	}

	for _, condition := range allConditions {
		product := cardtrader.Product{
			User: cardtrader.ProductUserInfo{
				SellsViaHub: false,
			},
			Properties: cardtrader.ProductProperties{
				Condition: condition,
			},
		}
		matchingProduct := watchConditionsMatchProduct(&watch, product)
		if matchingProduct {
			t.Errorf("watch has matched a product that is not sold via hub")
		}
	}
}

func TestWatchWithNoConditionMatchesNoProduct(t *testing.T) {
	watch := mongo.Watch{
		Name:          watchName,
		ExpansionID:   expansionID,
		ExpansionName: expansionName,
		BlueprintID:   blueprintID,
		Foil:          false,
	}

	for _, condition := range allConditions {
		product := cardtrader.Product{
			User: cardtrader.ProductUserInfo{
				SellsViaHub: true,
			},
			Properties: cardtrader.ProductProperties{
				Condition: condition,
			},
		}
		matchingProduct := watchConditionsMatchProduct(&watch, product)
		if matchingProduct {
			t.Errorf("no condition watch has matched a product")
		}
	}
}

func TestWatchProductMatchesByLanguage(t *testing.T) {
	watch := mongo.Watch{
		Name:          watchName,
		ExpansionID:   expansionID,
		ExpansionName: expansionName,
		BlueprintID:   blueprintID,
		Foil:          false,
		Language:      mongo.WatchLanguageEn,
		Condition:     mongo.WatchConditionAny,
	}

	for _, condition := range allConditions {
		product := cardtrader.Product{
			User: cardtrader.ProductUserInfo{
				SellsViaHub: true,
			},
			Properties: cardtrader.ProductProperties{
				Language:  mongo.WatchLanguageEn,
				Condition: condition,
			},
		}
		matchingProduct := watchConditionsMatchProduct(&watch, product)
		if !matchingProduct {
			t.Errorf("watch does not match a product by language")
		}
	}
	for _, condition := range allConditions {
		product := cardtrader.Product{
			User: cardtrader.ProductUserInfo{
				SellsViaHub: true,
			},
			Properties: cardtrader.ProductProperties{
				Language:  mongo.WatchLanguageDe,
				Condition: condition,
			},
		}
		matchingProduct := watchConditionsMatchProduct(&watch, product)
		if matchingProduct {
			t.Errorf("watch has matched a product with a different language")
		}
	}
}
