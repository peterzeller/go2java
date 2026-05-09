package syntaxtree

import (
	"go/ast"
	"go/types"
	"path/filepath"
	"testing"
)

func TestLoadTypedFileIncludesVariableType(t *testing.T) {
	tree, err := LoadTypedFile(filepath.Join("testdata", "example.go"))
	if err != nil {
		t.Fatalf("LoadTypedFile() error = %v", err)
	}

	obj := objectForDefinition(t, tree, "answer")

	if obj.Type().String() != "int" {
		t.Fatalf("answer type = %q, want %q", obj.Type().String(), "int")
	}
}

func objectForDefinition(t *testing.T, tree *TypedFile, name string) types.Object {
	t.Helper()

	var found types.Object
	ast.Inspect(tree.File, func(node ast.Node) bool {
		ident, ok := node.(*ast.Ident)
		if !ok || ident.Name != name {
			return true
		}

		if obj := tree.TypesInfo.Defs[ident]; obj != nil {
			found = obj
			return false
		}

		return true
	})

	if found == nil {
		t.Fatalf("definition for %q not found", name)
	}

	return found
}
