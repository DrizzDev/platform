package sqlite

import (
	"context"
	"database/sql"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/correlation"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
	"github.com/DrizzDev/platform/internal/capture/domain/trace"
)

// Read returns every entry for a trace in recorded order, reconstructing each into a validated value. A row that cannot
// be reconstructed is a corruption and fails the read rather than yielding a partial value.
func (store Store) Read(scope context.Context, subject trace.Trace) ([]journal.Entry, error) {
	var entries []journal.Entry

	failure := store.observe(scope, probe{name: "read", work: func(scope context.Context) error {
		rows, failure := store.handle.QueryContext(scope,
			"SELECT "+columns+" FROM journal WHERE trace = ? ORDER BY id", subject.String())
		if failure != nil {
			return failure
		}
		defer func() { _ = rows.Close() }()
		gathered, failure := store.collect(rows)
		if failure != nil {
			return failure
		}
		entries = gathered
		return nil
	}})
	return entries, failure
}

func (Store) collect(rows *sql.Rows) ([]journal.Entry, error) {
	var entries []journal.Entry
	for rows.Next() {
		var line row
		if failure := rows.Scan(line.targets()...); failure != nil {
			return nil, failure
		}
		entry, failure := line.entry()
		if failure != nil {
			return nil, failure
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

type row struct {
	sequence int64
	payload  []byte
	trace    string
	span     string
	parent   string
	mark     string
	origin   string
	fidelity string
	category string
	artifact string
	state    string
}

func (row row) entry() (journal.Entry, error) {
	thread, failure := trace.New(row.trace)
	if failure != nil {
		return journal.Entry{}, Corrupt{}
	}
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
		Trace: thread, Span: here, Parent: above, Sequence: row.sequence, Mark: mark.Mark(row.mark),
	})
	if failure != nil {
		return journal.Entry{}, Corrupt{}
	}

	reference := digest.Digest{}
	if row.artifact != "" {
		reference, failure = digest.New(row.artifact)
		if failure != nil {
			return journal.Entry{}, Corrupt{}
		}
	}

	entry, failure := journal.New(journal.Input{
		Correlation: link,
		Payload:     row.payload,
		Artifact:    reference,
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
