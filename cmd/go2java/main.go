package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go2java/internal/syntaxtree"
	"go2java/internal/translate"
)

func main() {
	write := flag.Bool("write", false, "write translated Java for examples/*.go into translated/")
	flag.Parse()

	if *write {
		if err := writeExamples(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [-write] <go-file>\n", os.Args[0])
		os.Exit(2)
	}

	javaCode, err := translateOne(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(javaCode)
}

func writeExamples() error {
	examples, err := filepath.Glob(filepath.Join("examples", "*.go"))
	if err != nil {
		return err
	}
	if len(examples) == 0 {
		return fmt.Errorf("no examples found")
	}
	if err := os.MkdirAll("translated", 0o755); err != nil {
		return err
	}

	for _, ex := range examples {
		javaCode, err := translateOne(ex)
		if err != nil {
			return err
		}
		base := strings.TrimSuffix(filepath.Base(ex), ".go")
		out := filepath.Join("translated", base+".java")
		if err := os.WriteFile(out, []byte(javaCode), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func translateOne(path string) (string, error) {
	tree, err := syntaxtree.LoadTypedFile(path)
	if err != nil {
		return "", err
	}
	program, err := translate.GoToIntermediate(tree)
	if err != nil {
		return "", err
	}
	return translate.IntermediateToJava(program)
}
