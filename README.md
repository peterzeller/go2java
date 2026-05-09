# go2java

`go2java` is a project template for a tool that translates Go code to Java code.

The first milestone is intentionally small: load a Go source file and build its
typed syntax tree. Later translation passes can use that typed tree to map Go
declarations, expressions, and types into Java code.

## Current Features

- CLI command for loading one Go file.
- Typed syntax tree loading with:
  - `golang.org/x/tools/go/packages`
  - `go/ast`
  - `go/types`
- Basic diagnostics for package loading and type checking errors.

## Project Structure

```text
.
├── cmd/go2java/          # CLI entry point
├── internal/syntaxtree/  # Typed Go syntax tree loading
├── go.mod
├── .gitignore
└── README.md
```

## Usage

Install dependencies:

```sh
go mod tidy
```

Load and summarize the typed syntax tree for a Go file:

```sh
go run ./cmd/go2java ./path/to/file.go
```

Example output:

```text
Package: main
File: ./path/to/file.go
Declarations: 2
Definitions: 4
Uses: 7
Implicits: 0
Scopes: 3
Selections: 0
```

## Next Steps

- Add an intermediate representation for translation.
- Map Go package/type/function declarations to Java equivalents.
- Implement expression and statement translation.
- Add golden-file tests for translated Java output.
