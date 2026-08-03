package mcp

type outcome string

const (
	success     outcome = "SUCCESS"
	cancelled   outcome = "CANCELLED"
	rejected    outcome = "REJECTED"
	interrupted outcome = "INTERRUPTED"
)
