package device

import "time"

// Instruction is the device-authorization challenge the provider issues: the
// code the user enters, where they enter it, and how the poll is paced and
// bounded. Code is the secret poll credential; User and Verification are shown.
type Instruction struct {
	Code         string
	User         string
	Verification string
	Interval     time.Duration
	Expiry       time.Time
}
