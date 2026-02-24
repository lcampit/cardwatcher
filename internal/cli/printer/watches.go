package printer

import (
	"os"

	api "github.com/lcampit/cardwatcher/gen/go/v1"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintWatchesTable(watches []*api.Watch) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Watches", "Watch ID", "Expansion ID", "Blueprint ID", "Expansion Name", "Name", "Condition", "Foil"})
	for _, watch := range watches {
		t.AppendRow(table.Row{"", watch.WatchId, watch.ExpansionId, watch.BlueprintId, watch.ExpansionName, watch.Name, NormalizeWatchCondition(watch.Condition), watch.Foil})
		t.AppendSeparator()
	}
	t.Render()
}

func NormalizeWatchCondition(condition api.Condition) string {
	switch condition {
	case api.Condition_CONDITION_NEAR_MINT:
		return "NEAR MINT"
	case api.Condition_CONDITION_MODERATELY_PLAYED:
		return "MODERATELY PLAYED"
	case api.Condition_CONDITION_SLIGHTLY_PLAYED:
		return "SLIGHTLY PLAYED"
	case api.Condition_CONDITION_PLAYED:
		return "PLAYED"
	case api.Condition_CONDITION_POOR:
		return "POOR"
	}
	return condition.String()
}
