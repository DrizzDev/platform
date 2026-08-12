package mcp

import (
	"context"

	protocol "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
)

type installing struct {
	Serial string `json:"serial" jsonschema:"the serial of the device"`
	Path   string `json:"path" jsonschema:"the path to the application package file"`
}

type naming struct {
	Serial string `json:"serial" jsonschema:"the serial of the device"`
	App    string `json:"app" jsonschema:"the application's package or bundle identifier"`
}

// application is one application in a listing tool's structured output.
type application struct {
	Id   string `json:"id" jsonschema:"the application identifier"`
	Name string `json:"name" jsonschema:"the application name"`
	Note string `json:"note" jsonschema:"the application version or process id"`
}

type catalogue struct {
	Apps []application `json:"apps" jsonschema:"the applications"`
}

type space struct {
	Megabytes int `json:"megabytes" jsonschema:"the free disk space in megabytes"`
}

// applications registers the application-action capabilities as MCP tools.
func (server Server) applications(perform Perform) {
	shelf := catalog.New()

	add, _ := shelf.Lookup(catalog.Install)
	protocol.AddTool(server.server, &protocol.Tool{Name: add.Title(), Description: add.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input installing) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Install(scope, operator.Package{Serial: input.Serial, Path: input.Path}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Installed."), nil, nil
		})

	start, _ := shelf.Lookup(catalog.Launch)
	protocol.AddTool(server.server, &protocol.Tool{Name: start.Title(), Description: start.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input naming) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Launch(scope, operator.Application{Serial: input.Serial, App: input.App}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Launched."), nil, nil
		})

	stop, _ := shelf.Lookup(catalog.Terminate)
	protocol.AddTool(server.server, &protocol.Tool{Name: stop.Title(), Description: stop.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input naming) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Terminate(scope, operator.Application{Serial: input.Serial, App: input.App}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Terminated."), nil, nil
		})

	wipe, _ := shelf.Lookup(catalog.Wipe)
	protocol.AddTool(server.server, &protocol.Tool{Name: wipe.Title(), Description: wipe.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input naming) (*protocol.CallToolResult, any, error) {
			if _, failure := perform.Wipe(scope, operator.Application{Serial: input.Serial, App: input.App}); failure != nil {
				return server.refuse(failure), nil, nil
			}
			return server.done("Cleared."), nil, nil
		})
}

// inspections registers the application and device read capabilities as MCP tools.
func (server Server) inspections(perform Perform) {
	shelf := catalog.New()

	have, _ := shelf.Lookup(catalog.Installed)
	protocol.AddTool(server.server, &protocol.Tool{Name: have.Title(), Description: have.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, catalogue, error) {
			listing, failure := perform.Installed(scope, operator.Target{Serial: input.Serial})
			if failure != nil {
				return server.refuse(failure), catalogue{}, nil
			}
			return nil, server.applist(listing), nil
		})

	live, _ := shelf.Lookup(catalog.Running)
	protocol.AddTool(server.server, &protocol.Tool{Name: live.Title(), Description: live.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, catalogue, error) {
			listing, failure := perform.Running(scope, operator.Target{Serial: input.Serial})
			if failure != nil {
				return server.refuse(failure), catalogue{}, nil
			}
			return nil, server.applist(listing), nil
		})

	front, _ := shelf.Lookup(catalog.Foreground)
	protocol.AddTool(server.server, &protocol.Tool{Name: front.Title(), Description: front.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, any, error) {
			report, failure := perform.Foreground(scope, operator.Target{Serial: input.Serial})
			if failure != nil {
				return server.refuse(failure), nil, nil
			}
			return &protocol.CallToolResult{Content: []protocol.Content{&protocol.TextContent{Text: report.Text}}}, nil, nil
		})

	link, _ := shelf.Lookup(catalog.Url)
	protocol.AddTool(server.server, &protocol.Tool{Name: link.Title(), Description: link.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, any, error) {
			report, failure := perform.Url(scope, operator.Target{Serial: input.Serial})
			if failure != nil {
				return server.refuse(failure), nil, nil
			}
			return &protocol.CallToolResult{Content: []protocol.Content{&protocol.TextContent{Text: report.Text}}}, nil, nil
		})

	free, _ := shelf.Lookup(catalog.Disk)
	protocol.AddTool(server.server, &protocol.Tool{Name: free.Title(), Description: free.Summary()},
		func(scope context.Context, _ *protocol.CallToolRequest, input handle) (*protocol.CallToolResult, space, error) {
			measure, failure := perform.Disk(scope, operator.Target{Serial: input.Serial})
			if failure != nil {
				return server.refuse(failure), space{}, nil
			}
			return nil, space{Megabytes: measure.Value}, nil
		})
}

func (Server) applist(listing operator.Listing) catalogue {
	apps := make([]application, 0, len(listing.Apps))
	for _, item := range listing.Apps {
		apps = append(apps, application{Id: item.Id, Name: item.Name, Note: item.Note})
	}
	return catalogue{Apps: apps}
}
