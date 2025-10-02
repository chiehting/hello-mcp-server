package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
		return nil, HelloWorldResult{Message: "Hello, World! Hi Justin!"}, nil
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

	// Create an HTTP handler for the MCP server.
	httpHandler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server { return server },
		nil,
	)

	// Determine the port to listen on.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("MCP server listening on :%s", port)
	// Start the HTTP server with logging middleware
	loggedHandler := loggingMiddleware(httpHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), loggedHandler); err != nil {
		log.Fatalf("MCP server failed: %v", err)
	}
}

// loggingMiddleware logs all incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("MCP Server: Received %s request for path: %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
