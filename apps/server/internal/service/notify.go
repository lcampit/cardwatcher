package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lcampit/cardwatcher/apps/server/internal/cardtrader"
	"github.com/lcampit/cardwatcher/apps/server/internal/mongo"
)

func (s *service) watchAndNotify() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	watches, err := s.mongoAdapter.GetWatches(ctx)
	if err != nil {
		s.logger.Error("getting watches", slog.Any("error", err))
		return
	}

	s.logger.Debug("handling notifications", slog.Int("blueprintCount", len(watches)))
	for _, watch := range watches {
		blueprintProducts, err := s.cardtraderAdapter.GetProducts(ctx, watch.BlueprintID, watch.Foil)
		if err != nil {
			s.logger.Error("error getting current pricing",
				slog.Uint64("blueprintId", watch.BlueprintID),
				slog.String("blueprintName", watch.Name),
				slog.Any("error", err))
			continue
		}

		var msg string
		if len(blueprintProducts) == 0 {
			msg = fmt.Sprintf("Looks like no one is selling %s today", watch.Name)
			s.logger.Info("no products found for blueprint",
				slog.Uint64("blueprintId", watch.BlueprintID),
				slog.String("blueprintName", watch.Name))
		} else {
			for _, product := range blueprintProducts {
				if watchConditionsMatchProduct(watch, product) {
					msg = fmt.Sprintf("%s price is %d today", watch.Name, product.Price.Cents)
					s.logger.Info("found price for blueprint",
						slog.Uint64("blueprintPricing", product.Price.Cents),
						slog.Uint64("blueprintId", watch.BlueprintID),
						slog.String("blueprintName", watch.Name))
					break
				}
			}
		}

		err = s.ntfyAdapter.Notify(ctx, msg)
		if err != nil {
			s.logger.Error("creating notification", slog.Any("error", err))
		}
	}
}

func watchConditionsMatchProduct(watch *mongo.Watch, product cardtrader.Product) bool {
	languageMatch := watch.Language == product.Properties.Language
	sellsViaHub := product.User.SellsViaHub
	if watch.Condition == mongo.WatchConditionAny {
		return sellsViaHub && languageMatch
	}

	conditionMatch := watch.Condition == product.Properties.Condition
	return languageMatch && conditionMatch && sellsViaHub
}
