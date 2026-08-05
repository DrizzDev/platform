package publication

import (
	"encoding/base64"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/application/login"
)

const (
	schema     = 1
	entropy    = 16
	saturation = 16
)

type draft struct {
	candidate login.Candidate
	handle    string
	revision  uint64
}

func (publisher Publisher) record(draft draft) (credential.Record, error) {
	token := draft.candidate.Token
	return credential.New(credential.Input{
		Issuer:   token.Issuer,
		Client:   token.Client,
		Handle:   draft.handle,
		Subject:  token.Subject,
		Session:  publisher.session,
		Method:   token.Method,
		Refresh:  token.Refresh,
		Issued:   token.Issued,
		Expiry:   token.Expiry,
		Revision: draft.revision,
		Schema:   schema,
	})
}

// mint draws a unique per-attempt handle so a concurrent login writes a distinct
// vault candidate key and can never overwrite or supersede another attempt.
func (publisher Publisher) mint() (string, error) {
	bytes, failure := publisher.random.Bytes(entropy)
	if failure != nil {
		return "", failure
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (publisher Publisher) receipt(record credential.Record) login.Receipt {
	return login.Receipt{
		Subject: record.Subject(),
		Session: record.Session(),
		Method:  record.Method(),
		Expiry:  record.Expiry(),
	}
}
