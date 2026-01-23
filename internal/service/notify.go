package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func (s *service) WatchAndNotify() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	watches, err := s.mongoAdapter.GetWatches(ctx)
	if err != nil {
		s.logger.Error("error getting watches", slog.Any("error", err))
		return
	}

	s.logger.Debug("handling notifications", slog.Int("blueprintCount", len(watches)))
	for _, watch := range watches {
		blueprintPricing, err := s.cardtraderAdapter.GetCurrentPricing(ctx, watch)
		if err != nil {
			s.logger.Error("error getting current pricing",
				slog.Int("blueprintID", watch.BlueprintId),
				slog.Any("error", err))
			break
		}

		var msg string
		if blueprintPricing == 0 {
			msg = fmt.Sprintf("Looks like no one is selling %s today", watch.Name)
			s.logger.Info("no products found for blueprint",
				slog.Int("blueprintID", watch.BlueprintId),
				slog.String("blueprintName", watch.Name))
		} else {
			msg = fmt.Sprintf("%s price is %d today", watch.Name, blueprintPricing)
			s.logger.Info("found price for blueprint",
				slog.Int("blueprintPricing", blueprintPricing),
				slog.Int("blueprintID", watch.BlueprintId),
				slog.String("blueprintName", watch.Name))
		}

		err = s.ntfyAdapter.Notify(ctx, msg)
		if err != nil {
			s.logger.Error("error creting notification", slog.Any("error", err))
		}
	}
}
