package mcp

import (
	"context"

	protocol "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
)

type booting struct {
	Image string `json:"image" jsonschema:"the name of the emulator image to start"`
}

type gallery struct {
	Images []string `json:"images" jsonschema:"the available emulator image names"`
}

// emulators registers the emulator-management capabilities as MCP tools.
func (server Server) emulators(perform Perform) {
	shelf := catalog.New()

	list, _ := shelf.Lookup(catalog.Images)
	protocol.AddTool(server.server, &protocol.Tool{Name: list.Title(), Description: list.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, _ struct{}) (*protocol.CallToolResult, gallery, error) {
			images, failure := perform.Images(scope)
			if failure != nil {
				return server.refuse(failure), gallery{}, nil
			}
			return nil, gallery{Images: images.Names}, nil
		})

	start, _ := shelf.Lookup(catalog.Boot)
	protocol.AddTool(server.server, &protocol.Tool{Name: start.Title(), Description: start.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input booting) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Boot(scope, operator.Image{Name: input.Image}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Booting."), nil, nil
		})

	hold, _ := shelf.Lookup(catalog.Pause)
	protocol.AddTool(server.server, &protocol.Tool{Name: hold.Title(), Description: hold.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Pause(scope, operator.Target{Serial: input.Serial}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Paused."), nil, nil
		})

	play, _ := shelf.Lookup(catalog.Resume)
	protocol.AddTool(server.server, &protocol.Tool{Name: play.Title(), Description: play.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Resume(scope, operator.Target{Serial: input.Serial}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Resumed."), nil, nil
		})
}
