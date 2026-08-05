package method

type Method string

const (
	Device   Method = "DEVICE"
	Browser  Method = "BROWSER"
	Workload Method = "WORKLOAD"
)

func (method Method) Valid() bool {
	switch method {
	case Browser, Device, Workload:
		return true
	default:
		return false
	}
}

func (method Method) String() string {
	return string(method)
}
