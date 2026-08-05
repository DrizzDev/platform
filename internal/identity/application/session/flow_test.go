package session_test

import (
	"context"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/marking"
	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/application/session"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	handle "github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

type vault func(context.Context) (credential.Record, error)

func (vault vault) Active(scope context.Context) (credential.Record, error) { return vault(scope) }

type refresh func(context.Context, credential.Record) (session.Renewal, error)

func (refresh refresh) Renew(scope context.Context, record credential.Record) (session.Renewal, error) {
	return refresh(scope, record)
}

type reader func(context.Context) (epoch.Epoch, error)

func (reader reader) Read(scope context.Context) (epoch.Epoch, error) { return reader(scope) }

type clock struct{}

func (clock) Now() time.Time { return time.Unix(3000, 0) }

// ledger is a fake session Publication that records what it was asked to fence
// and publish, and can be told to deny either transition.
type ledger struct {
	mark    marking.Marking
	settled session.Candidate
	deny    error
	reject  error
}

func (ledger *ledger) Attempt(_ context.Context, mark marking.Marking) error {
	ledger.mark = mark
	return ledger.deny
}

func (ledger *ledger) Publish(_ context.Context, candidate session.Candidate) (session.Receipt, error) {
	ledger.settled = candidate
	if ledger.reject != nil {
		return session.Receipt{}, ledger.reject
	}
	return session.Receipt{
		Subject: candidate.Prior.Subject(),
		Session: candidate.Prior.Session(),
		Method:  candidate.Prior.Method(),
		Expiry:  candidate.Renewal.Expiry,
	}, nil
}

type harness struct {
	test *testing.T
}

func (harness harness) active() credential.Record {
	harness.test.Helper()
	who, _ := subject.New("google-oauth2|first")
	owner, _ := handle.New("LOCAL")
	record, failure := credential.New(credential.Input{
		Issuer: "https://issuer.example/", Client: "native", Handle: "prior",
		Subject: who, Session: owner, Method: method.Browser,
		Refresh: []byte("old-refresh"), Issued: time.Unix(1000, 0), Expiry: time.Unix(2000, 0),
		Revision: 3, Schema: 1,
	})
	if failure != nil {
		harness.test.Fatal(failure)
	}
	return record
}

func (harness harness) build(options session.Options) session.Flow {
	harness.test.Helper()
	made, failure := session.New(options)
	if failure != nil {
		harness.test.Fatal(failure)
	}
	return made
}

func (harness harness) renewal() session.Renewal {
	return session.Renewal{Refresh: []byte("new-refresh"), Access: []byte("access-token"), Expiry: time.Unix(4000, 0)}
}

func TestAccess(test *testing.T) {
	test.Parallel()

	fixture := harness{test: test}
	record := fixture.active()
	credential, failure := fixture.build(session.Options{
		Vault:       vault(func(context.Context) (credential.Record, error) { return record, nil }),
		Refresh:     refresh(func(context.Context, credential.Record) (session.Renewal, error) { return fixture.renewal(), nil }),
		Publication: &ledger{},
		Epoch:       reader(func(context.Context) (epoch.Epoch, error) { return epoch.Epoch(7), nil }),
		Clock:       clock{},
	}).Access(context.Background(), session.Input{})
	if failure != nil {
		test.Fatalf("access failed: %v", failure)
	}
	if string(credential.Token()) != "access-token" {
		test.Fatalf("token = %q", credential.Token())
	}
	if !credential.Expiry().Equal(time.Unix(4000, 0)) {
		test.Fatalf("expiry = %v", credential.Expiry())
	}
}

func TestCurrent(test *testing.T) {
	test.Parallel()

	fixture := harness{test: test}
	record := fixture.active()
	store := &ledger{}
	result := fixture.build(session.Options{
		Vault:       vault(func(context.Context) (credential.Record, error) { return record, nil }),
		Refresh:     refresh(func(context.Context, credential.Record) (session.Renewal, error) { return fixture.renewal(), nil }),
		Publication: store,
		Epoch:       reader(func(context.Context) (epoch.Epoch, error) { return epoch.Epoch(7), nil }),
		Clock:       clock{},
	}).Current(context.Background(), session.Input{})

	if result.Failed() {
		fault, _ := result.Failure()
		test.Fatalf("renewal failed: %s", fault.Code())
	}
	if result.Subject().String() != "google-oauth2|first" || result.Standing().String() != "ACTIVE" {
		test.Fatalf("result = %+v", result)
	}
	if store.mark.Attempt().Revision() != 3 || uint64(store.mark.Attempt().Epoch()) != 7 {
		test.Fatalf("fenced marking = %+v", store.mark.Attempt())
	}
	if uint64(store.settled.Expected) != 7 || store.settled.Prior.Revision() != 3 {
		test.Fatalf("published candidate = %+v", store.settled)
	}
	if !store.settled.Moment.Equal(time.Unix(3000, 0)) {
		test.Fatalf("moment = %v", store.settled.Moment)
	}
}

func TestRequired(test *testing.T) {
	test.Parallel()

	fixture := harness{test: test}
	cases := map[string]session.Options{
		"missing": {
			Vault: vault(func(context.Context) (credential.Record, error) { return credential.Record{}, session.Missing{} }),
		},
		"fenced": {
			Vault:       vault(func(context.Context) (credential.Record, error) { return fixture.active(), nil }),
			Publication: &ledger{deny: session.Uncertain{}},
		},
		"unrenewable": {
			Vault: vault(func(context.Context) (credential.Record, error) { return fixture.active(), nil }),
			Refresh: refresh(func(context.Context, credential.Record) (session.Renewal, error) {
				return session.Renewal{}, session.Uncertain{}
			}),
		},
	}
	for name, options := range cases {
		options := options
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			fixture := harness{test: test}
			options.Clock = clock{}
			options.Epoch = reader(func(context.Context) (epoch.Epoch, error) { return 1, nil })
			if options.Refresh == nil {
				options.Refresh = refresh(func(context.Context, credential.Record) (session.Renewal, error) { return fixture.renewal(), nil })
			}
			if options.Publication == nil {
				options.Publication = &ledger{}
			}
			result := fixture.build(options).Current(context.Background(), session.Input{})
			fault, present := result.Failure()
			if !present || fault.Code() != code.Required {
				test.Fatalf("code = %v (present %v)", fault.Code(), present)
			}
		})
	}
}

func TestStorage(test *testing.T) {
	test.Parallel()

	fixture := harness{test: test}
	result := fixture.build(session.Options{
		Vault:       vault(func(context.Context) (credential.Record, error) { return credential.Record{}, session.Storage{} }),
		Refresh:     refresh(func(context.Context, credential.Record) (session.Renewal, error) { return fixture.renewal(), nil }),
		Publication: &ledger{},
		Epoch:       reader(func(context.Context) (epoch.Epoch, error) { return 1, nil }),
		Clock:       clock{},
	}).Current(context.Background(), session.Input{})

	fault, present := result.Failure()
	if !present || fault.Code() != code.Storage {
		test.Fatalf("code = %v (present %v)", fault.Code(), present)
	}
}
