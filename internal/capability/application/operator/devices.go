package operator

import "context"

// Devices lists the connected devices so a person or an agent can choose one for a capability. Listing is a read, so
// it is not recorded as an execution.
func (operator Operator) Devices(scope context.Context) (Roster, error) {
	discovery := operator.flow.Discover(scope)
	if reason, failed := discovery.Failure(); failed {
		return Roster{}, operator.refuse(reason)
	}
	serials := make([]string, 0, len(discovery.Devices()))
	for _, candidate := range discovery.Devices() {
		serials = append(serials, candidate.Serial().String())
	}
	return Roster{Serials: serials}, nil
}
