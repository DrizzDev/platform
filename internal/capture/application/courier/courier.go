package courier

// Courier drains unsynchronized journal entries to the cloud. Each entry uploads its artifact, then itself, then flips
// to synced; a crash before the flip leaves the entry pending, and it is safely re-sent and deduped by identity.
type Courier struct {
	vault    vault
	ledger   ledger
	uploader uploader
}
