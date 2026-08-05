package device

import "context"

// Display is the device-owned port that presents the challenge to the user,
// telling them where to go and which code to enter.
type Display interface {
	Show(context.Context, Instruction) error
}
