package journal

// Transition moves one entry to a new synchronization state. The target State is a field so a third state can be
// applied later without changing this shape.
type Transition struct {
	Entry Entry
	State State
}
