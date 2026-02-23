package printer

import (
	"os"

	api "github.com/lcampit/card-watcher-server/internal/api/v1"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintExpansionTable(expansions []*api.Expansion) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Expansions", "Name", "Code", "ID"})
	for _, expansion := range expansions {
		t.AppendRow(table.Row{"", expansion.Name, expansion.Code, expansion.Id})
		t.AppendSeparator()
	}
	t.Render()
}
