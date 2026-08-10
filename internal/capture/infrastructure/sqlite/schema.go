package sqlite

import (
	"strings"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

// column binds a journal column to how its value is written from an entry and read back into a row.
// The insert, the select, and the scan all derive from the one ordered layout below, so their column order can never drift apart.
type column struct {
	name  string
	read  func(*row) any
	write func(journal.Entry) any
}

var layout = []column{
	{
		name:  "trace",
		read:  func(row *row) any { return &row.trace },
		write: func(entry journal.Entry) any { return entry.Correlation().Trace().String() },
	},
	{
		name:  "span",
		read:  func(row *row) any { return &row.span },
		write: func(entry journal.Entry) any { return entry.Correlation().Span().String() },
	},
	{
		name:  "parent",
		read:  func(row *row) any { return &row.parent },
		write: func(entry journal.Entry) any { return entry.Correlation().Parent().String() },
	},
	{
		name:  "sequence",
		read:  func(row *row) any { return &row.sequence },
		write: func(entry journal.Entry) any { return entry.Correlation().Sequence() },
	},
	{
		name:  "mark",
		read:  func(row *row) any { return &row.mark },
		write: func(entry journal.Entry) any { return entry.Correlation().Mark().String() },
	},
	{
		name:  "origin",
		read:  func(row *row) any { return &row.origin },
		write: func(entry journal.Entry) any { return entry.Origin().String() },
	},
	{
		name:  "fidelity",
		read:  func(row *row) any { return &row.fidelity },
		write: func(entry journal.Entry) any { return entry.Fidelity().String() },
	},
	{
		name:  "category",
		read:  func(row *row) any { return &row.category },
		write: func(entry journal.Entry) any { return entry.Category().String() },
	},
	{
		name:  "payload",
		read:  func(row *row) any { return &row.payload },
		write: func(entry journal.Entry) any { return entry.Payload() },
	},
	{
		name:  "artifact",
		read:  func(row *row) any { return &row.artifact },
		write: func(entry journal.Entry) any { return entry.Artifact().String() },
	},
	{
		name:  "state",
		read:  func(row *row) any { return &row.state },
		write: func(entry journal.Entry) any { return entry.State().String() },
	},
}

var columns = func() string {
	names := make([]string, len(layout))
	for index, column := range layout {
		names[index] = column.name
	}
	return strings.Join(names, ", ")
}()

var placeholders = strings.TrimSuffix(strings.Repeat("?, ", len(layout)), ", ")

func (Store) values(entry journal.Entry) []any {
	values := make([]any, len(layout))
	for index, column := range layout {
		values[index] = column.write(entry)
	}
	return values
}

func (row *row) targets() []any {
	targets := make([]any, len(layout))
	for index, column := range layout {
		targets[index] = column.read(row)
	}
	return targets
}
