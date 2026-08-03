package observability

// code is an approved internal fault identifier. Values are stable dotted
// lowercase event names owned by this package, not by call sites.
type code string

const command code = "command.failed"

// defect is a deliberate, reportable internal fault. It carries only an approved
// code: no cause, message, or attribute channel exists, so no client, user,
// credential, or provider content can reach diagnostics or the error sink.
type defect struct {
	name code
}

func (defect defect) event() string {
	return string(defect.name)
}
