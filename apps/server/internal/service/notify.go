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
		availableProducts, err := s.cardtraderAdapter.GetProducts(ctx, watch.BlueprintID, watch.Foil)
		if err != nil {
			s.logger.Error("error getting current pricing",
				slog.Uint64("blueprintId", watch.BlueprintID),
				slog.String("blueprintName", watch.Name),
				slog.Any("error", err))
			continue
		}

		msg := buildNotificationMessage(watch, availableProducts)
		s.logger.Info("notification message built",
			slog.String("message", msg),
			slog.Uint64("blueprintId", watch.BlueprintID),
			slog.String("blueprintName", watch.Name))

		err = s.ntfyAdapter.Notify(ctx, msg)
		if err != nil {
			s.logger.Error("creating notification", slog.Any("error", err))
		}
	}
}

func buildNotificationMessage(watch *mongo.Watch, availableProducts []cardtrader.Product) string {
	if len(availableProducts) == 0 {
		return fmt.Sprintf("looks like no one is selling %s today", watch.Name)
	}

	for _, product := range availableProducts {
		if watchConditionsMatchProduct(watch, product) {
			return fmt.Sprintf("%s price is %d today", watch.Name, product.Price.Cents)
		}
	}

	return fmt.Sprintf("no product aviailable found for %s, condition %s, language %s, foil %b", watch.Name, watch.Condition, watch.Language, watch.Foil)
}

func watchConditionsMatchProduct(watch *mongo.Watch, product cardtrader.Product) bool {
	languageMatch := watch.Language == mongo.WatchLanguageAny || watch.Language == product.Properties.Language
	conditionMatch := watch.Condition == mongo.WatchConditionAny || watch.Condition == product.Properties.Condition
	sellsViaHub := product.User.SellsViaHub
	return languageMatch && conditionMatch && sellsViaHub
}
