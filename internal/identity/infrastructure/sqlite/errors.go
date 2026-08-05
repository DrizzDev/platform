package sqlite

// benign marks a failure that is an expected coordination outcome rather than a
// fault, so observation records it as declined instead of raising an error event.
type benign interface {
	benign() bool
}

// Contention reports that a lease or revision is owned by a competing holder.
type Contention struct{}

func (Contention) Error() string {
	return "the resource is held by another owner"
}

func (Contention) benign() bool {
	return true
}

// Fenced reports that a revision was already attempted and cannot be replayed.
type Fenced struct{}

func (Fenced) Error() string {
	return "the revision was already attempted"
}

func (Fenced) benign() bool {
	return true
}

// Absent reports that no head pointer exists for a session yet.
type Absent struct{}

func (Absent) Error() string {
	return "no active credential is present"
}

func (Absent) benign() bool {
	return true
}
