package printer

import (
	"os"

	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintExpansionTable(expansions []*apiv1.Expansion) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Expansions", "Name", "Code", "ID"})
	for _, expansion := range expansions {
		t.AppendRow(table.Row{"", expansion.Name, expansion.Code, expansion.Id})
		t.AppendSeparator()
	}
	t.Render()
}
