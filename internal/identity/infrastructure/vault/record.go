package vault

import (
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

type wire struct {
	Issuer   string    `json:"issuer"`
	Client   string    `json:"client"`
	Handle   string    `json:"handle"`
	Subject  string    `json:"subject"`
	Session  string    `json:"session"`
	Method   string    `json:"method"`
	Refresh  []byte    `json:"refresh"`
	Issued   time.Time `json:"issued"`
	Expiry   time.Time `json:"expiry"`
	Revision uint64    `json:"revision"`
	Schema   uint32    `json:"schema"`
}

func (wire wire) record() (credential.Record, error) {
	account, failure := subject.New(wire.Subject)
	if failure != nil {
		return credential.Record{}, failure
	}
	handle, failure := session.New(wire.Session)
	if failure != nil {
		return credential.Record{}, failure
	}
	return credential.New(credential.Input{
		Session:  handle,
		Subject:  account,
		Issuer:   wire.Issuer,
		Client:   wire.Client,
		Handle:   wire.Handle,
		Issued:   wire.Issued,
		Expiry:   wire.Expiry,
		Schema:   wire.Schema,
		Refresh:  wire.Refresh,
		Revision: wire.Revision,

		Method: method.Method(wire.Method),
	})
}
