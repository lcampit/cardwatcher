package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

func (s *service) WatchAndNotify() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	watches, err := s.mongoAdapter.GetWatches(ctx)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	log.Debug().Msgf("handling notifications for %d blueprints", len(watches))
	for _, watch := range watches {
		blueprintPricing, err := s.cardtraderAdapter.GetCurrentPricing(ctx, watch)
		if err != nil {
			log.Error().Err(err).Msg("")
			break
		}

		var msg string
		if blueprintPricing == 0 {
			msg = fmt.Sprintf("Looks like no one is selling %s today", watch.Name)
			log.Info().Msgf("no products found for blueprint %d (%s)", watch.BlueprintId, watch.Name)
		} else {
			msg = fmt.Sprintf("%s price is %d today", watch.Name, blueprintPricing)
			log.Info().Msgf("found price %d for blueprint %d (%s)", blueprintPricing, watch.BlueprintId, watch.Name)
		}

		s.ntfyAdapter.Notify(ctx, msg)
	}
}
