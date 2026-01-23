package mongo

import (
	"context"
	"log/slog"
	"testing"

	"card-watcher/internal/entities"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCRUDWatches(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.DiscardHandler)
	mongoAdapter := NewMongoAdapter(logger, testHost, testPort, testDatabase)

	watchID1 := bson.NewObjectID()
	watchID2 := bson.NewObjectID()
	watchID3 := bson.NewObjectID()

	watch1 := &entities.Watch{
		WatchID:     watchID1,
		ExpansionID: 1,
		BlueprintID: 1,
		Condition:   entities.WatchConditionNM,
		Foil:        false,
	}
	watch2 := &entities.Watch{
		WatchID:     watchID2,
		ExpansionID: 2,
		BlueprintID: 2,
		Condition:   entities.WatchConditionNM,
		Foil:        false,
	}

	watch3 := &entities.Watch{
		WatchID:     watchID3,
		ExpansionID: 2,
		BlueprintID: 3,
		Condition:   entities.WatchConditionNM,
		Foil:        false,
	}

	insertedWatchID1, err := mongoAdapter.SaveWatch(ctx, watch1)
	assert.Nil(t, err, "saving watch 1 failed")
	insertedWatchID2, err := mongoAdapter.SaveWatch(ctx, watch2)
	assert.Nil(t, err, "saving watch 2 failed")
	insertedWatchID3, err := mongoAdapter.SaveWatch(ctx, watch3)
	assert.Nil(t, err, "saving watch 3 failed")

	watchFromDB1, err := mongoAdapter.GetWatchByWatchID(ctx, insertedWatchID1)
	assert.Nil(t, err, "getting watch 1 failed")
	assert.NotNil(t, watchFromDB1, "watcher 1 not found by get")

	assert.Equal(t, watch1.WatchID, watchFromDB1.WatchID)
	assert.Equal(t, watch1.ExpansionID, watchFromDB1.ExpansionID)
	assert.Equal(t, watch1.BlueprintID, watchFromDB1.BlueprintID)
	assert.Equal(t, watch1.Condition, watchFromDB1.Condition)
	assert.Equal(t, watch1.Foil, watchFromDB1.Foil)

	err = mongoAdapter.DeleteWatchByID(ctx, insertedWatchID1)
	assert.Nil(t, err)
	err = mongoAdapter.DeleteWatchByID(ctx, insertedWatchID2)
	assert.Nil(t, err)

	absentWatches, err := mongoAdapter.GetWatchByWatchID(ctx, insertedWatchID1)
	assert.Nil(t, absentWatches)
	assert.Nil(t, err)

	absentWatches, err = mongoAdapter.GetWatchByWatchID(ctx, insertedWatchID2)
	assert.Nil(t, absentWatches)
	assert.Nil(t, err)

	watchFromDB3, err := mongoAdapter.GetWatchByWatchID(ctx, insertedWatchID3)
	assert.Nil(t, err, "getting watch 3 failed")
	assert.NotNil(t, watchFromDB3, "watcher 3 not found by get")

	assert.Equal(t, watch3.WatchID, watchFromDB3.WatchID)
	assert.Equal(t, watch3.ExpansionID, watchFromDB3.ExpansionID)
	assert.Equal(t, watch3.BlueprintID, watchFromDB3.BlueprintID)
	assert.Equal(t, watch3.Condition, watchFromDB3.Condition)
	assert.Equal(t, watch3.Foil, watchFromDB3.Foil)

	err = mongoAdapter.DeleteWatchByID(ctx, watchID3.Hex())
	assert.Nil(t, err)
}
