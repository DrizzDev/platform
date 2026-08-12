package mcp

import (
	"context"
	"strings"

	protocol "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
)

// extent is the structured output of the screen-size tool.
type extent struct {
	Width  int `json:"width" jsonschema:"the screen width in pixels"`
	Height int `json:"height" jsonschema:"the screen height in pixels"`
}

// observations registers the data-returning observation capabilities as MCP tools.
func (server Server) observations(perform Perform) {
	shelf := catalog.New()

	snap, _ := shelf.Lookup(catalog.Snapshot)
	protocol.AddTool(server.server, &protocol.Tool{Name: snap.Title(), Description: snap.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, any, error) {
			shot, failure := perform.Snapshot(scope, operator.Target{Serial: input.Serial})
			if failure != nil {
				return server.refuse(failure), nil, nil
			}
			image := &protocol.ImageContent{Data: shot.Image, MIMEType: mime + strings.ToLower(shot.Format)}
			content := []protocol.Content{image, &protocol.TextContent{Text: shot.Hierarchy}}
			if note, failure := server.keep(artifact{extension: shot.Format, content: shot.Image}); failure == nil {
				content = append(content, note)
			}
			return &protocol.CallToolResult{Content: content}, nil, nil
		})

	tree, _ := shelf.Lookup(catalog.Hierarchy)
	protocol.AddTool(server.server, &protocol.Tool{Name: tree.Title(), Description: tree.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, any, error) {
			read, failure := perform.Hierarchy(scope, operator.Target{Serial: input.Serial})
			if failure != nil {
				return server.refuse(failure), nil, nil
			}
			content := []protocol.Content{&protocol.TextContent{Text: read.Hierarchy}}
			if note, failure := server.keep(artifact{extension: "xml", content: []byte(read.Hierarchy)}); failure == nil {
				content = append(content, note)
			}
			return &protocol.CallToolResult{Content: content}, nil, nil
		})

	size, _ := shelf.Lookup(catalog.Dimensions)
	protocol.AddTool(server.server, &protocol.Tool{Name: size.Title(), Description: size.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, extent, error) {
			measured, failure := perform.Dimensions(scope, operator.Target{Serial: input.Serial})
			if failure != nil {
				return server.refuse(failure), extent{}, nil
			}
			return nil, extent{Width: measured.Width, Height: measured.Height}, nil
		})
}
