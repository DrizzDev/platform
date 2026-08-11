package mcp

import (
	"context"
	"errors"
	"strings"

	protocol "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
)

// snapshot is the input to the screenshot tool. Its schema is inferred from these fields.
type snapshot struct {
	Serial string `json:"serial" jsonschema:"the serial of the device to capture"`
}

// roster is the structured output of the devices tool.
type roster struct {
	Serials []string `json:"serials" jsonschema:"the connected device serials"`
}

// register exposes the catalogued capabilities as MCP tools. Their names and descriptions come from the catalog, so
// the agent connection and the command line always present the same capabilities described the same way.
func (server Server) register(perform Perform) {
	shelf := catalog.New()
	screen, _ := shelf.Lookup(catalog.Screenshot)
	list, _ := shelf.Lookup(catalog.Devices)

	protocol.AddTool(server.server, &protocol.Tool{Name: screen.Name(), Description: screen.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input snapshot) (*protocol.CallToolResult, any, error) {
			shot, failure := perform.Screenshot(scope, operator.Target{Serial: input.Serial})
			if failure != nil {
				return server.refuse(failure), nil, nil
			}
			image := &protocol.ImageContent{Data: shot.Image, MIMEType: "image/" + strings.ToLower(shot.Format)}
			return &protocol.CallToolResult{Content: []protocol.Content{image}}, nil, nil
		})

	protocol.AddTool(server.server, &protocol.Tool{Name: list.Name(), Description: list.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, _ struct{}) (*protocol.CallToolResult, roster, error) {
			devices, failure := perform.Devices(scope)
			if failure != nil {
				return server.refuse(failure), roster{}, nil
			}
			return nil, roster{Serials: devices.Serials}, nil
		})
}

// refuse maps a capability failure to a tool error the model can read, carrying only the code's safe detail.
func (Server) refuse(failure error) *protocol.CallToolResult {
	message := "The request could not be completed. Try again."
	var refusal operator.Refusal
	if errors.As(failure, &refusal) {
		message = refusal.Code.Detail()
	}
	return &protocol.CallToolResult{IsError: true, Content: []protocol.Content{&protocol.TextContent{Text: message}}}
}
