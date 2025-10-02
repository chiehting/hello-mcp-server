package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godoctor-client",
	Short: "godoctor-client is a CLI client for the godoctor MCP server",
	Long:  `A CLI client that interacts with the godoctor MCP server to call various tools.`,
}

var serverAddr string

func executeTool(toolName string, args map[string]any) {
	ctx := context.Background()

	client := mcp.NewClient(
		&mcp.Implementation{Name: "godoctor-client"},
		nil,
	)

	// transport := &mcp.CommandTransport{Command: exec.Command("./bin/godoctor")}
	transport := &mcp.StreamableClientTransport{Endpoint: serverAddr}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect to MCP server: %v", err)
	}
	defer session.Close()

	params := &mcp.CallToolParams{
		Name:      toolName,
		Arguments: args,
	}

	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("Tool call failed: %v", err)
	}

	// Handle unstructured content first
	for _, content := range res.Content {
		if textContent, ok := content.(*mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			// Handle other content types if necessary
			fmt.Printf("Received content of type %T: %+v\n", content, content)
		}
	}

	// Handle structured content
	if res.StructuredContent != nil {
		// Attempt to marshal and pretty print the structured content
		prettyJSON, err := json.MarshalIndent(res.StructuredContent, "", "  ")
		if err != nil {
			log.Printf("Warning: Could not pretty print structured content: %v", err)
			fmt.Printf("Structured Content: %+v\n", res.StructuredContent)
		} else {
			fmt.Println(string(prettyJSON))
		}
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&serverAddr, "server-addr", "a", "http://localhost:8080", "Address of the MCP server")

	// helloWorld command
	helloWorldCmd := &cobra.Command{
		Use:   "helloWorld",
		Short: "Calls the helloWorld tool",
		Run: func(cmd *cobra.Command, args []string) {
			executeTool("helloWorld", nil)
		},
	}
	rootCmd.AddCommand(helloWorldCmd)

	// godoc command
	var godocPackage string
	var godocSymbol string
	godocCmd := &cobra.Command{
		Use:   "godoc",
		Short: "Calls the godoc tool",
		Run: func(cmd *cobra.Command, args []string) {
			if godocPackage == "" {
				log.Fatal("Error: --package is required for godoc tool")
			}
			executeTool("godoc", map[string]any{"package": godocPackage, "symbol": godocSymbol})
		},
	}
	godocCmd.Flags().StringVarP(&godocPackage, "package", "p", "", "The Go package to document (e.g., fmt)")
	godocCmd.Flags().StringVarP(&godocSymbol, "symbol", "s", "", "The symbol within the package to document (optional)")
	godocCmd.MarkFlagRequired("package")
	rootCmd.AddCommand(godocCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
