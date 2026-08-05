package vault

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
)

// limit is the platform credential-blob ceiling (Windows CRED_MAX_CREDENTIAL_BLOB_SIZE).
// A valid record always fits it by construction; the guard is a defense-in-depth backstop.
const limit = 2560

type Vault struct {
	store Store
}

func New(options Options) (Vault, error) {
	if failure := options.validate(); failure != nil {
		return Vault{}, failure
	}
	return Vault{store: options.Store}, nil
}

func (vault Vault) Write(scope context.Context, record credential.Record) error {
	stored := wire{
		Issuer:   record.Issuer(),
		Client:   record.Client(),
		Handle:   record.Handle(),
		Issued:   record.Issued(),
		Expiry:   record.Expiry(),
		Schema:   record.Schema(),
		Refresh:  record.Refresh(),
		Revision: record.Revision(),
		Method:   record.Method().String(),
		Subject:  record.Subject().String(),
		Session:  record.Session().String(),
	}
	encoded, failure := json.Marshal(stored)
	if failure != nil {
		return failure
	}
	if len(encoded) > limit {
		return errors.New("credential record exceeds the vault entry limit")
	}
	return vault.store.Write(scope, Entry{Key: string(record.Key()), Value: encoded})
}

func (vault Vault) Read(scope context.Context, key credential.Key) (credential.Record, error) {
	encoded, failure := vault.store.Read(scope, string(key))
	if failure != nil {
		return credential.Record{}, failure
	}
	if len(encoded) > limit {
		return credential.Record{}, errors.New("stored credential exceeds the vault entry limit")
	}
	var stored wire
	if failure := json.Unmarshal(encoded, &stored); failure != nil {
		return credential.Record{}, failure
	}
	return stored.record()
}

func (vault Vault) Delete(scope context.Context, key credential.Key) error {
	return vault.store.Delete(scope, string(key))
}
