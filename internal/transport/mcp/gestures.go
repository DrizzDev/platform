package mcp

import (
	"context"

	protocol "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
)

// beat is the default swipe duration in milliseconds for a tool call that does not carry one.
const beat = 300

type spot struct {
	X int `json:"x" jsonschema:"the horizontal coordinate, in pixels"`
	Y int `json:"y" jsonschema:"the vertical coordinate, in pixels"`
}

type dragging struct {
	Serial string `json:"serial" jsonschema:"the serial of the device to swipe"`
	From   spot   `json:"from" jsonschema:"the point to start the swipe from"`
	To     spot   `json:"to" jsonschema:"the point to end the swipe at"`
}

type squeezing struct {
	Serial string `json:"serial" jsonschema:"the serial of the device to pinch"`
	Centre spot   `json:"centre" jsonschema:"the centre point of the pinch"`
	Inner  int    `json:"startRadius" jsonschema:"the starting radius from the centre, in pixels"`
	Outer  int    `json:"endRadius" jsonschema:"the ending radius from the centre, in pixels"`
}

type pressing struct {
	Serial string `json:"serial" jsonschema:"the serial of the device"`
	Button string `json:"button" jsonschema:"the button to press, such as up, down, select, menu, or home"`
}

type typing struct {
	Serial string `json:"serial" jsonschema:"the serial of the device"`
	Text   string `json:"text" jsonschema:"the text to type"`
}

type locating struct {
	Serial    string  `json:"serial" jsonschema:"the serial of the device"`
	Latitude  float64 `json:"latitude" jsonschema:"the latitude to report, in degrees"`
	Longitude float64 `json:"longitude" jsonschema:"the longitude to report, in degrees"`
}

type handle struct {
	Serial string `json:"serial" jsonschema:"the serial of the device"`
}

// interactions registers the interaction capabilities as MCP tools, reading each name and description from the catalog
// so the agent connection and the command line always present the same capabilities.
func (server Server) interactions(perform Perform) {
	shelf := catalog.New()

	drag, _ := shelf.Lookup(catalog.Swipe)
	protocol.AddTool(server.server, &protocol.Tool{Name: drag.Title(), Description: drag.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input dragging) (*protocol.CallToolResult, any, error) {
			shifted := operator.Drag{
				Serial:       input.Serial,
				From:         operator.Spot{X: input.From.X, Y: input.From.Y},
				To:           operator.Spot{X: input.To.X, Y: input.To.Y},
				Milliseconds: beat,
			}
			if _, failure := perform.Swipe(scope, shifted); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Swiped."), nil, nil
		})

	zoom, _ := shelf.Lookup(catalog.Pinch)
	protocol.AddTool(server.server, &protocol.Tool{Name: zoom.Title(), Description: zoom.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input squeezing) (*protocol.CallToolResult, any, error) {
			pinched := operator.Squeeze{
				Serial: input.Serial,
				Centre: operator.Spot{X: input.Centre.X, Y: input.Centre.Y},
				From:   input.Inner,
				To:     input.Outer,
			}
			if _, failure := perform.Pinch(scope, pinched); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Pinched."), nil, nil
		})

	key, _ := shelf.Lookup(catalog.Press)
	protocol.AddTool(server.server, &protocol.Tool{Name: key.Title(), Description: key.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input pressing) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Press(scope, operator.Key{Serial: input.Serial, Button: input.Button}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Pressed."), nil, nil
		})

	entry, _ := shelf.Lookup(catalog.Type)
	protocol.AddTool(server.server, &protocol.Tool{Name: entry.Title(), Description: entry.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input typing) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Type(scope, operator.Entry{Serial: input.Serial, Text: input.Text}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Typed."), nil, nil
		})

	server.navigations(perform)
}

// navigations registers the location and no-argument interaction capabilities as MCP tools.
func (server Server) navigations(perform Perform) {
	shelf := catalog.New()

	fix, _ := shelf.Lookup(catalog.Locate)
	protocol.AddTool(server.server, &protocol.Tool{Name: fix.Title(), Description: fix.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input locating) (*protocol.CallToolResult, any, error) {
			placed := operator.Fix{Serial: input.Serial, Lat: input.Latitude, Lon: input.Longitude}
			if _, failure := perform.Locate(scope, placed); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Location set."), nil, nil
		})

	wipe, _ := shelf.Lookup(catalog.Clear)
	protocol.AddTool(server.server, &protocol.Tool{Name: wipe.Title(), Description: wipe.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Clear(scope, operator.Target{Serial: input.Serial}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Cleared."), nil, nil
		})

	back, _ := shelf.Lookup(catalog.Back)
	protocol.AddTool(server.server, &protocol.Tool{Name: back.Title(), Description: back.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Back(scope, operator.Target{Serial: input.Serial}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Went back."), nil, nil
		})

	home, _ := shelf.Lookup(catalog.Home)
	protocol.AddTool(server.server, &protocol.Tool{Name: home.Title(), Description: home.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Home(scope, operator.Target{Serial: input.Serial}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Went home."), nil, nil
		})
}
