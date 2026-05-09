package main

import (
	"fmt"
	"os"

	"go2java/internal/syntaxtree"
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

	fmt.Printf("Package: %s\n", tree.Package.Name())
	fmt.Printf("File: %s\n", tree.Path)
	fmt.Printf("Declarations: %d\n", len(tree.File.Decls))
	fmt.Printf("Definitions: %d\n", len(tree.TypesInfo.Defs))
	fmt.Printf("Uses: %d\n", len(tree.TypesInfo.Uses))
	fmt.Printf("Implicits: %d\n", len(tree.TypesInfo.Implicits))
	fmt.Printf("Scopes: %d\n", len(tree.TypesInfo.Scopes))
	fmt.Printf("Selections: %d\n", len(tree.TypesInfo.Selections))
}
