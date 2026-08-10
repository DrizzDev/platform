package sqlite

import (
	"context"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/correlation"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
	"github.com/DrizzDev/platform/internal/capture/domain/trace"
)

// Read returns every entry for a trace in recorded order, reconstructing each into a validated value. A row that cannot
// be reconstructed is a corruption and fails the read rather than yielding a partial value (REL-021).
func (store Store) Read(scope context.Context, subject trace.Trace) ([]journal.Entry, error) {
	var entries []journal.Entry

	failure := store.observe(scope, probe{name: "read", work: func(scope context.Context) error {
		rows, failure := store.handle.QueryContext(scope,
			"SELECT span, parent, sequence, mark, origin, fidelity, category, payload, state "+
				"FROM journal WHERE trace = ? ORDER BY id", subject.String())

		if failure != nil {
			return failure
		}
		defer func() { _ = rows.Close() }()

		for rows.Next() {
			var line row
			if failure := rows.Scan(&line.span, &line.parent, &line.sequence, &line.mark,
				&line.origin, &line.fidelity, &line.category, &line.payload, &line.state); failure != nil {
				return failure
			}
			entry, failure := line.entry(subject)
			if failure != nil {
				return failure
			}
			entries = append(entries, entry)
		}
		return rows.Err()
	}})
	return entries, failure
}

type row struct {
	sequence int64
	payload  []byte
	span     string
	parent   string
	mark     string
	origin   string
	fidelity string
	category string
	state    string
}

func (row row) entry(subject trace.Trace) (journal.Entry, error) {
	here, failure := span.New(row.span)

	if failure != nil {
		return journal.Entry{}, Corrupt{}
	}

	above := span.Span{}
	if row.parent != "" {
		above, failure = span.New(row.parent)
		if failure != nil {
			return journal.Entry{}, Corrupt{}
		}
	}

	link, failure := correlation.New(correlation.Input{
		Trace: subject, Span: here, Parent: above, Sequence: row.sequence, Mark: mark.Mark(row.mark),
	})
	if failure != nil {
		return journal.Entry{}, Corrupt{}
	}

	entry, failure := journal.New(journal.Input{
		Correlation: link,
		Payload:     row.payload,
		State:       journal.State(row.state),
		Origin:      origin.Origin(row.origin),
		Fidelity:    fidelity.Fidelity(row.fidelity),
		Category:    category.Category(row.category),
	})
	if failure != nil {
		return journal.Entry{}, Corrupt{}
	}
	return entry, nil
}
