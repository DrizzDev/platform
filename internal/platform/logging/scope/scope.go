package scope

type Scope struct {
	request     string
	session     string
	correlation string
}

func New(input Input) Scope {
	return Scope{
		request:     input.Request,
		session:     input.Session,
		correlation: input.Correlation,
	}
}
