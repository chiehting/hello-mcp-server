# Hello MCP Server

This project is an example from the Google Codelabs "[Gemini CLI Hands-on](https://codelabs.developers.google.com/gemini-cli-hands-on)" tutorial.

## Overview

## Building and Running

To build all the applications into a `bin` directory, navigate to the project root and run:

```bash
mkdir -p bin
go build -o bin/godoctor-cli ./cmd/godoctor-cli
go build -o bin/godoctor-client ./cmd/godoctor-client
go build -o bin/godoctor-http ./cmd/godoctor-http
go build -o bin/godoctor-sse ./cmd/godoctor-sse
```

To run a specific application, for example, `godoctor-cli`:

```bash
./bin/godoctor-cli
```
