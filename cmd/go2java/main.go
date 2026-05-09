package main

import (
	"fmt"
	"os"

	"go2java/internal/syntaxtree"
	"go2java/internal/translate"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <go-file>\n", os.Args[0])
		os.Exit(2)
	}

	tree, err := syntaxtree.LoadTypedFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	program, err := translate.GoToIntermediate(tree)
	if err != nil {
		fmt.Fprintf(os.Stderr, "translate: %v\n", err)
		os.Exit(1)
	}

	javaCode, err := translate.IntermediateToJava(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "java print: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(javaCode)
}
