package artifact

import "context"

// Footprint reports the total bytes of every published object,
// so a reclaim pass can weigh the artifacts against the disk budget.
func (store Store) Footprint(scope context.Context) (int64, error) {
	var total int64

	failure := store.observe(scope, probe{name: "footprint", work: func(context.Context) error {
		stored, failure := store.walk()
		if failure != nil {
			return failure
		}
		for _, object := range stored {
			total += object.size
		}
		return nil
	}})
	return total, failure
}
