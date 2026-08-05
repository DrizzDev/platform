package session

type Options struct {
	Vault       Vault
	Refresh     Refresh
	Publication Publication
	Epoch       Epoch
	Clock       Clock
}
