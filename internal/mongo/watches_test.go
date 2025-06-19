package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRUDWatches(t *testing.T) {
	ctx := context.Background()
	mongoAdapter := NewMongoAdapter(testHost, testPort, testDatabase)

	watch1 := &Watch{
		Name:        "name 1",
		WatchId:     "watchId1",
		UserId:      "userId1",
		ExpansionId: "ExpansionId1",
		BlueprintId: "BlueprintId1",
		Condition:   NEAR_MINT,
		Foil:        false,
	}
	watch2 := &Watch{
		Name:        "name 2",
		WatchId:     "watchId2",
		UserId:      "userId1",
		ExpansionId: "ExpansionId2",
		BlueprintId: "BlueprintId2",
		Condition:   NEAR_MINT,
		Foil:        false,
	}

	watch3 := &Watch{
		Name:        "name 3",
		WatchId:     "watchId3",
		UserId:      "userId3",
		ExpansionId: "ExpansionId3",
		BlueprintId: "BlueprintId3",
		Condition:   NEAR_MINT,
		Foil:        false,
	}

	err := mongoAdapter.SaveWatch(ctx, watch1)
	assert.Nil(t, err, "saving watch 1 failed")
	err = mongoAdapter.SaveWatch(ctx, watch2)
	assert.Nil(t, err, "saving watch 2 failed")
	err = mongoAdapter.SaveWatch(ctx, watch3)
	assert.Nil(t, err, "saving watch 3 failed")

	watchFromDb1, err := mongoAdapter.GetWatchByWatchId(ctx, "watchId1")
	assert.Nil(t, err, "getting watch 1 failed")
	assert.NotNil(t, watchFromDb1, "watcher 1 not found by get")

	assert.Equal(t, watch1.Name, watchFromDb1.Name)
	assert.Equal(t, watch1.WatchId, watchFromDb1.WatchId)
	assert.Equal(t, watch1.UserId, watchFromDb1.UserId)
	assert.Equal(t, watch1.ExpansionId, watchFromDb1.ExpansionId)
	assert.Equal(t, watch1.BlueprintId, watchFromDb1.BlueprintId)
	assert.Equal(t, watch1.Condition, watchFromDb1.Condition)
	assert.Equal(t, watch1.Foil, watchFromDb1.Foil)

	err = mongoAdapter.DeleteWatchesByUserId(ctx, "userId1")
	assert.Nil(t, err)

	absentWatches, err := mongoAdapter.GetWatchByWatchId(ctx, "watchId1")
	assert.Nil(t, absentWatches)

	absentWatches, err = mongoAdapter.GetWatchByWatchId(ctx, "watchId2")
	assert.Nil(t, absentWatches)

	watchFromDb3, err := mongoAdapter.GetWatchByWatchId(ctx, "watchId3")
	assert.Nil(t, err, "getting watch 3 failed")
	assert.NotNil(t, watchFromDb3, "watcher 3 not found by get")

	assert.Equal(t, watch3.Name, watchFromDb3.Name)
	assert.Equal(t, watch3.WatchId, watchFromDb3.WatchId)
	assert.Equal(t, watch3.UserId, watchFromDb3.UserId)
	assert.Equal(t, watch3.ExpansionId, watchFromDb3.ExpansionId)
	assert.Equal(t, watch3.BlueprintId, watchFromDb3.BlueprintId)
	assert.Equal(t, watch3.Condition, watchFromDb3.Condition)
	assert.Equal(t, watch3.Foil, watchFromDb3.Foil)

	err = mongoAdapter.DeleteWatch(ctx, "watchId3")
	assert.Nil(t, err)
}
