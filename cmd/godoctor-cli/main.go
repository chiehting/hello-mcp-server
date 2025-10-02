package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HelloWorldArgs represents the arguments for the helloWorld tool. Since it has no parameters, it's an empty struct.
type HelloWorldArgs struct{}

// HelloWorldResult represents the result of the helloWorld tool.
type HelloWorldResult struct {
	Message string `json:"message"`
}

// GodocArgs represents the arguments for the godoc tool.
type GodocArgs struct {
	Package string `json:"package" jsonschema:"the Go package to document"`
	Symbol  string `json:"symbol,omitempty" jsonschema:"the symbol within the package to document (optional)"`
}

// GodocResult represents the result of the godoc tool.
type GodocResult struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func main() {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "godoctor"},
		nil,
	)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "helloWorld",
		Description: "A simple tool that returns 'Hello, World!'.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args HelloWorldArgs) (*mcp.CallToolResult, HelloWorldResult, error) {
		return nil, HelloWorldResult{Message: "Hello, World!"}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "godoc",
		Description: "Invokes the 'go doc' command to retrieve documentation for a Go package or symbol.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args GodocArgs) (*mcp.CallToolResult, GodocResult, error) {
		cmdArgs := []string{"doc", args.Package}
		if args.Symbol != "" {
			cmdArgs = append(cmdArgs, args.Symbol)
		}

		cmd := exec.Command("go", cmdArgs...)
		output, err := cmd.CombinedOutput()

		if err != nil {
			return nil, GodocResult{Output: string(output), Error: err.Error()}, nil
		}

		return nil, GodocResult{Output: string(output)}, nil
	})

	// Run the server on the stdio transport.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
		os.Exit(1)
	}
}
