package mongo

import (
	"card-watcher/internal/entities"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCRUDWatches(t *testing.T) {
	ctx := context.Background()
	mongoAdapter := NewMongoAdapter(testHost, testPort, testDatabase)

	watchId1 := bson.NewObjectID()
	watchId2 := bson.NewObjectID()
	watchId3 := bson.NewObjectID()

	watch1 := &entities.Watch{
		WatchId:     watchId1,
		UserId:      "userId1",
		ExpansionId: "ExpansionId1",
		BlueprintId: "BlueprintId1",
		Condition:   entities.WATCH_CONDITION_NEAR_MINT,
		Foil:        false,
	}
	watch2 := &entities.Watch{
		WatchId:     watchId2,
		UserId:      "userId1",
		ExpansionId: "ExpansionId2",
		BlueprintId: "BlueprintId2",
		Condition:   entities.WATCH_CONDITION_NEAR_MINT,
		Foil:        false,
	}

	watch3 := &entities.Watch{
		WatchId:     watchId3,
		UserId:      "userId3",
		ExpansionId: "ExpansionId3",
		BlueprintId: "BlueprintId3",
		Condition:   entities.WATCH_CONDITION_NEAR_MINT,
		Foil:        false,
	}

	insertedWatchId1, err := mongoAdapter.SaveWatch(ctx, watch1)
	assert.Nil(t, err, "saving watch 1 failed")
	insertedWatchId2, err := mongoAdapter.SaveWatch(ctx, watch2)
	assert.Nil(t, err, "saving watch 2 failed")
	insertedWatchId3, err := mongoAdapter.SaveWatch(ctx, watch3)
	assert.Nil(t, err, "saving watch 3 failed")

	watchFromDb1, err := mongoAdapter.GetWatchByWatchId(ctx, insertedWatchId1)
	assert.Nil(t, err, "getting watch 1 failed")
	assert.NotNil(t, watchFromDb1, "watcher 1 not found by get")

	assert.Equal(t, watch1.WatchId, watchFromDb1.WatchId)
	assert.Equal(t, watch1.UserId, watchFromDb1.UserId)
	assert.Equal(t, watch1.ExpansionId, watchFromDb1.ExpansionId)
	assert.Equal(t, watch1.BlueprintId, watchFromDb1.BlueprintId)
	assert.Equal(t, watch1.Condition, watchFromDb1.Condition)
	assert.Equal(t, watch1.Foil, watchFromDb1.Foil)

	err = mongoAdapter.DeleteWatchesByUserId(ctx, "userId1")
	assert.Nil(t, err)

	absentWatches, err := mongoAdapter.GetWatchByWatchId(ctx, insertedWatchId1)
	assert.Nil(t, absentWatches)

	absentWatches, err = mongoAdapter.GetWatchByWatchId(ctx, insertedWatchId2)
	assert.Nil(t, absentWatches)

	watchFromDb3, err := mongoAdapter.GetWatchByWatchId(ctx, insertedWatchId3)
	assert.Nil(t, err, "getting watch 3 failed")
	assert.NotNil(t, watchFromDb3, "watcher 3 not found by get")

	assert.Equal(t, watch3.WatchId, watchFromDb3.WatchId)
	assert.Equal(t, watch3.UserId, watchFromDb3.UserId)
	assert.Equal(t, watch3.ExpansionId, watchFromDb3.ExpansionId)
	assert.Equal(t, watch3.BlueprintId, watchFromDb3.BlueprintId)
	assert.Equal(t, watch3.Condition, watchFromDb3.Condition)
	assert.Equal(t, watch3.Foil, watchFromDb3.Foil)

	err = mongoAdapter.DeleteWatchById(ctx, watchId3.Hex())
	assert.Nil(t, err)
}
