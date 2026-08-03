package mcp

// excess marks a frame that exceeds the maximum message size.
type excess struct{}

func (excess) Error() string {
	return "MCP message exceeds the maximum size"
}

// malformed marks a frame that is truncated or not a valid JSON-RPC message.
type malformed struct{}

func (malformed) Error() string {
	return "MCP message is malformed"
}
