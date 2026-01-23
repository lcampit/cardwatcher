package mongo

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"card-watcher/internal/entities"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCRUDWatches(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mongoAdapter := NewMongoAdapter(logger, testHost, testPort, testDatabase)

	watchID1 := bson.NewObjectID()
	watchID2 := bson.NewObjectID()
	watchID3 := bson.NewObjectID()

	watch1 := &entities.Watch{
		WatchId:     watchID1,
		ExpansionId: 1,
		BlueprintId: 1,
		Condition:   entities.WATCH_CONDITION_NEAR_MINT,
		Foil:        false,
	}
	watch2 := &entities.Watch{
		WatchId:     watchID2,
		ExpansionId: 2,
		BlueprintId: 2,
		Condition:   entities.WATCH_CONDITION_NEAR_MINT,
		Foil:        false,
	}

	watch3 := &entities.Watch{
		WatchId:     watchID3,
		ExpansionId: 2,
		BlueprintId: 3,
		Condition:   entities.WATCH_CONDITION_NEAR_MINT,
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

	assert.Equal(t, watch1.WatchId, watchFromDB1.WatchId)
	assert.Equal(t, watch1.ExpansionId, watchFromDB1.ExpansionId)
	assert.Equal(t, watch1.BlueprintId, watchFromDB1.BlueprintId)
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

	assert.Equal(t, watch3.WatchId, watchFromDB3.WatchId)
	assert.Equal(t, watch3.ExpansionId, watchFromDB3.ExpansionId)
	assert.Equal(t, watch3.BlueprintId, watchFromDB3.BlueprintId)
	assert.Equal(t, watch3.Condition, watchFromDB3.Condition)
	assert.Equal(t, watch3.Foil, watchFromDB3.Foil)

	err = mongoAdapter.DeleteWatchByID(ctx, watchID3.Hex())
	assert.Nil(t, err)
}
