package publication

import (
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
)

type Options struct {
	Vault   Vault
	Ledger  Ledger
	Random  login.Random
	Session session.Session
}
