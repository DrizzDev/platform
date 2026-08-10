package sqlite

import "github.com/DrizzDev/platform/internal/capture/application/journal"

// entries is the ordered layout of the journal table.
var entries = layout[journal.Entry, row]{
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
