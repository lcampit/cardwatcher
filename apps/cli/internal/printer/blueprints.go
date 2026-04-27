package printer

import (
	"os"

	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintBlueprintsTable(blueprints []*apiv1.Blueprint) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Blueprint", "Expansion ID", "Name", "ID"})
	for _, blueprint := range blueprints {
		t.AppendRow(table.Row{"", blueprint.ExpansionId, blueprint.Name, blueprint.Id})
		t.AppendSeparator()
	}
	t.Render()
}
