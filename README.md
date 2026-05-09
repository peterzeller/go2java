# go2java

`go2java` translates a subset of Go into Java through a multi-pass pipeline.

## Current Features

- Load typed Go syntax trees (`go/ast`, `go/types`, `go/packages`).
- Translate Go AST into an intermediate AST (Go+Java friendly).
- Print intermediate AST as Java source code.
- End-to-end example tests that run both Go and Java and compare outputs.

## Supported Subset (current)

- Basic types: `int`, `string`, `bool`
- Functions with at most one return value
- Statements: expression statements, assignments, return, `if`, `for` (condition-only)
- Calls: basic calls including `fmt.Println` -> `System.out.println`

## Project Structure

```text
.
├── cmd/go2java/            # CLI entry point
├── examples/               # Go example files used for e2e translation tests
├── internal/intermediate/  # Intermediate AST nodes
├── internal/syntaxtree/    # Typed Go syntax tree loading
├── internal/translate/     # Go->IR and IR->Java translation
└── README.md
```

## Usage

```sh
go run ./cmd/go2java ./examples/hello.go
```

## Tests

```sh
go test ./...
```
