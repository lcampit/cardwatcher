package printer

import (
	"os"

	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintWatchesTable(watches []*apiv1.Watch) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Watches", "Watch ID", "Expansion ID", "Blueprint ID", "Expansion Name", "Name", "Condition", "Language", "Foil"})
	for _, watch := range watches {
		t.AppendRow(table.Row{"", watch.WatchId, watch.ExpansionId, watch.BlueprintId, watch.ExpansionName, watch.Name, normalizeWatchCondition(watch.Condition), normalizeWatchLanguage(watch.Language), watch.Foil})
		t.AppendSeparator()
	}
	t.Render()
}

func normalizeWatchCondition(condition apiv1.Condition) string {
	switch condition {
	case apiv1.Condition_CONDITION_NEAR_MINT:
		return "NEAR MINT"
	case apiv1.Condition_CONDITION_MODERATELY_PLAYED:
		return "MODERATELY PLAYED"
	case apiv1.Condition_CONDITION_SLIGHTLY_PLAYED:
		return "SLIGHTLY PLAYED"
	case apiv1.Condition_CONDITION_PLAYED:
		return "PLAYED"
	case apiv1.Condition_CONDITION_POOR:
		return "POOR"
	}
	return "ANY"
}

func normalizeWatchLanguage(language apiv1.Language) string {
	switch language {
	case apiv1.Language_LANGUAGE_EN:
		return "EN"
	case apiv1.Language_LANGUAGE_DE:
		return "DE"
	case apiv1.Language_LANGUAGE_FR:
		return "FR"
	case apiv1.Language_LANGUAGE_IT:
		return "IT"
	case apiv1.Language_LANGUAGE_JP:
		return "JP"
	case apiv1.Language_LANGUAGE_PT:
		return "PT"
	case apiv1.Language_LANGUAGE_ES:
		return "ES"
	}
	return "ANY"
}
