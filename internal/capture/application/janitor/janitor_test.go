package janitor_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/janitor"
	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/catalogue"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/policy"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
	"github.com/DrizzDev/platform/internal/capture/domain/trace"
)

var moment = time.Unix(1_000_000, 0)

type parcel struct {
	recorded time.Time
	step     span.Span
	trace    trace.Trace
	artifact digest.Digest
	class    category.Category
}

type book struct {
	parcels []parcel
	gone    map[string]bool
	claims  []journal.Claim
}

func (book *book) Settled(context.Context) ([]journal.Retained, error) {
	var out []journal.Retained
	for _, parcel := range book.parcels {
		if book.gone[parcel.step.String()] {
			continue
		}
		out = append(out, journal.Retained{
			Trace: parcel.trace, Span: parcel.step, Category: parcel.class, Recorded: parcel.recorded,
		})
	}
	return out, nil
}

func (book *book) Discard(_ context.Context, steps []span.Span) error {
	if book.gone == nil {
		book.gone = map[string]bool{}
	}
	for _, step := range steps {
		book.gone[step.String()] = true
	}
	return nil
}

func (book *book) Digests(context.Context) ([]digest.Digest, error) {
	var out []digest.Digest
	for _, parcel := range book.parcels {
		if book.gone[parcel.step.String()] || parcel.artifact.Empty() {
			continue
		}
		out = append(out, parcel.artifact)
	}
	return out, nil
}

func (book *book) Leases(context.Context) ([]journal.Claim, error) {
	return book.claims, nil
}

type item struct {
	modified time.Time
	digest   digest.Digest
}

type shelf struct {
	unit  int64
	items []item

	pruned     int
	restricted bool
	admitted   bool
}

func (shelf *shelf) Prune(_ context.Context, keep func(digest.Digest, time.Time) bool) error {
	var kept []item
	for _, object := range shelf.items {
		if keep(object.digest, object.modified) {
			kept = append(kept, object)
			continue
		}
		shelf.pruned++
	}
	shelf.items = kept
	return nil
}

func (shelf *shelf) Footprint(context.Context) (int64, error) {
	return int64(len(shelf.items)) * shelf.unit, nil
}

func (shelf *shelf) Restrict(context.Context) error {
	shelf.restricted = true
	return nil
}

func (shelf *shelf) Admit(context.Context) error {
	shelf.admitted = true
	shelf.restricted = false
	return nil
}

type bell struct {
	alerts []janitor.Pressure
}

func (bell *bell) Alert(_ context.Context, pressure janitor.Pressure) error {
	bell.alerts = append(bell.alerts, pressure)
	return nil
}

type clock struct{}

func (clock) Now() time.Time { return moment }

type fixture struct {
	test *testing.T
}

func (fixture fixture) trace(value string) trace.Trace {
	fixture.test.Helper()
	made, failure := trace.New(value)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) span(value string) span.Span {
	fixture.test.Helper()

	made, failure := span.New(value)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) digest(text string) digest.Digest {
	fixture.test.Helper()
	sum := sha256.Sum256([]byte(text))

	made, failure := digest.New(hex.EncodeToString(sum[:]))
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) catalogue() catalogue.Catalogue {
	fixture.test.Helper()
	rule, failure := policy.New(policy.Input{Limit: 1024, Retention: time.Hour, Upload: true})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	made, failure := catalogue.New(catalogue.Input{Policies: map[category.Category]policy.Policy{
		category.Tool:   rule,
		category.Screen: rule,
	}})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

type rig struct {
	book    *book
	shelf   *shelf
	bell    *bell
	options janitor.Options
}

func (fixture fixture) janitor(rig rig) janitor.Janitor {
	fixture.test.Helper()
	options := rig.options
	options.Archive = rig.book
	options.Vault = rig.shelf
	options.Notifier = rig.bell
	options.Clock = clock{}
	options.Catalogue = fixture.catalogue()
	made, failure := janitor.New(options)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func TestReclaimAged(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	object := kit.digest("art-0")
	shop := &book{parcels: []parcel{{
		trace: kit.trace("01HTRACE"), step: kit.span("hop-0"), class: category.Tool,
		recorded: moment.Add(-2 * time.Hour), artifact: object,
	}}}
	disk := &shelf{items: []item{{digest: object, modified: moment.Add(-2 * time.Hour)}}, unit: 1}

	report, failure := kit.janitor(rig{book: shop, shelf: disk, bell: &bell{}}).Run(context.Background())
	if failure != nil {
		test.Fatal(failure)
	}
	if report.Reclaimed != 1 || report.Degraded {
		test.Fatalf("report = %+v, want one reclaimed and not degraded", report)
	}
	if !shop.gone["hop-0"] {
		test.Fatal("aged synced entry was not discarded")
	}
	if len(disk.items) != 0 {
		test.Fatal("artifact left unreferenced was not swept")
	}
}

func TestSparesFresh(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	object := kit.digest("art-0")
	shop := &book{parcels: []parcel{{
		trace: kit.trace("01HTRACE"), step: kit.span("hop-0"), class: category.Tool,
		recorded: moment.Add(-time.Minute), artifact: object,
	}}}
	disk := &shelf{items: []item{{digest: object, modified: moment.Add(-time.Minute)}}, unit: 1}

	report, failure := kit.janitor(rig{book: shop, shelf: disk, bell: &bell{}}).Run(context.Background())
	if failure != nil {
		test.Fatal(failure)
	}
	if report.Reclaimed != 0 || len(shop.gone) != 0 {
		test.Fatalf("a within-retention entry under budget must survive, report = %+v", report)
	}
}

func TestProtectsLease(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	subject := kit.trace("01HTRACE")
	shop := &book{
		parcels: []parcel{{
			trace: subject, step: kit.span("hop-0"), class: category.Tool, recorded: moment.Add(-2 * time.Hour),
		}},
		claims: []journal.Claim{{Trace: subject, Until: moment.Add(time.Minute)}},
	}
	disk := &shelf{unit: 1}

	report, failure := kit.janitor(rig{book: shop, shelf: disk, bell: &bell{}}).Run(context.Background())
	if failure != nil {
		test.Fatal(failure)
	}
	if report.Reclaimed != 0 || shop.gone["hop-0"] {
		test.Fatal("an aged entry of a live-leased trace must be protected")
	}
}

func TestEscalatesUnderPressure(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	object := kit.digest("art-0")
	shop := &book{parcels: []parcel{{
		trace: kit.trace("01HTRACE"), step: kit.span("hop-0"), class: category.Tool,
		recorded: moment.Add(-2 * time.Minute), artifact: object,
	}}}
	disk := &shelf{items: []item{{digest: object, modified: moment.Add(-2 * time.Hour)}}, unit: 200}

	rig := rig{book: shop, shelf: disk, bell: &bell{}, options: janitor.Options{Ceiling: 100, Relief: 50}}
	report, failure := kit.janitor(rig).Run(context.Background())
	if failure != nil {
		test.Fatal(failure)
	}
	if report.Reclaimed != 1 || report.Degraded {
		test.Fatalf("pressure must reclaim within-retention synced, report = %+v", report)
	}
	if !disk.admitted || disk.restricted {
		test.Fatal("reclaim below relief must restore writes")
	}
}

func TestDegradesAtFloor(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	object := kit.digest("art-0")
	shop := &book{}
	disk := &shelf{items: []item{{digest: object, modified: moment}}, unit: 200}

	chime := &bell{}
	rig := rig{book: shop, shelf: disk, bell: chime, options: janitor.Options{Ceiling: 100, Relief: 50}}
	report, failure := kit.janitor(rig).Run(context.Background())
	if failure != nil {
		test.Fatal(failure)
	}
	if !report.Degraded || !disk.restricted {
		test.Fatal("over budget with nothing reclaimable must degrade writes")
	}
	if len(chime.alerts) != 1 || chime.alerts[0].Used != 200 || chime.alerts[0].Ceiling != 100 {
		test.Fatalf("pressure alert = %+v", chime.alerts)
	}
}
