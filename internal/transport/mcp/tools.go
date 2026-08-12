package mcp

import (
	"context"
	"errors"
	"strings"

	protocol "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
	"github.com/DrizzDev/platform/internal/platform/filesystem"
)

// mime is the media-type prefix for the captured image content returned to the agent.
const mime = "image/"

// artifact is a captured file to persist so the agent receives a path to work with, not only inline bytes.
type artifact struct {
	extension string
	content   []byte
}

// keep writes a captured artifact to a file and returns a note with its path, so the agent can reference the file
// rather than reaching for the device's native tooling.
func (Server) keep(item artifact) (*protocol.TextContent, error) {
	path, failure := filesystem.New().Save(filesystem.File{Extension: strings.ToLower(item.extension), Content: item.content})
	if failure != nil {
		return nil, failure
	}
	return &protocol.TextContent{Text: "Saved to " + path}, nil
}

// snapshot is the input to the screenshot tool. Its schema is inferred from these fields.
type snapshot struct {
	Serial string `json:"serial" jsonschema:"the serial of the device to capture"`
}

// roster is the structured output of the devices tool.
type roster struct {
	Serials []string `json:"serials" jsonschema:"the connected device serials"`
}

// contact is the input to the tap tool.
type contact struct {
	Serial string `json:"serial" jsonschema:"the serial of the device to tap"`
	X      int    `json:"x" jsonschema:"the horizontal coordinate to tap, in pixels"`
	Y      int    `json:"y" jsonschema:"the vertical coordinate to tap, in pixels"`
}

// register exposes the catalogued capabilities as MCP tools. Their names and descriptions come from the catalog, so
// the agent connection and the command line always present the same capabilities described the same way.
func (server Server) register(perform Perform) {
	shelf := catalog.New()
	screen, _ := shelf.Lookup(catalog.Screenshot)
	list, _ := shelf.Lookup(catalog.Devices)

	protocol.AddTool(server.server, &protocol.Tool{Name: screen.Title(), Description: screen.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input snapshot) (*protocol.CallToolResult, any, error) {
			shot, failure := perform.Screenshot(scope, operator.Target{Serial: input.Serial})
			if failure != nil {
				return server.refuse(failure), nil, nil
			}
			image := &protocol.ImageContent{Data: shot.Image, MIMEType: mime + strings.ToLower(shot.Format)}
			content := []protocol.Content{image}
			if note, failure := server.keep(artifact{extension: shot.Format, content: shot.Image}); failure == nil {
				content = append(content, note)
			}
			return &protocol.CallToolResult{Content: content}, nil, nil
		})

	protocol.AddTool(server.server, &protocol.Tool{Name: list.Title(), Description: list.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, _ struct{}) (*protocol.CallToolResult, roster, error) {
			devices, failure := perform.Devices(scope)
			if failure != nil {
				return server.refuse(failure), roster{}, nil
			}
			return nil, roster{Serials: devices.Serials}, nil
		})

	press, _ := shelf.Lookup(catalog.Tap)
	protocol.AddTool(server.server, &protocol.Tool{Name: press.Title(), Description: press.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input contact) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Tap(scope, operator.Contact{Serial: input.Serial, X: input.X, Y: input.Y}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Tapped."), nil, nil
		})

	server.interactions(perform)
	server.observations(perform)
	server.applications(perform)
	server.inspections(perform)
	server.emulators(perform)
}

// done renders a performed action that returns no data of its own as a short confirmation the agent can read.
func (Server) done(message string) *protocol.CallToolResult {
	return &protocol.CallToolResult{Content: []protocol.Content{&protocol.TextContent{Text: message}}}
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
