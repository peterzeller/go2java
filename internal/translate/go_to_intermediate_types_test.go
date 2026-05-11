package translate

import (
	"fmt"
	"go/parser"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go2java/internal/syntaxtree"
)

func TestLowerType_AllGoIntegerTypes(t *testing.T) {
	tests := map[string]string{
		"int":     "int",
		"int8":    "byte",
		"int16":   "short",
		"int32":   "int",
		"int64":   "long",
		"uint":    "int",
		"uint8":   "short",
		"uint16":  "int",
		"uint32":  "int",
		"uint64":  "long",
		"uintptr": "long",
		"byte":    "short",
		"rune":    "int",
	}

	for in, want := range tests {
		expr, err := parser.ParseExpr(in)
		if err != nil {
			t.Fatalf("parse %s: %v", in, err)
		}
		got, err := lowerType(expr)
		if err != nil {
			t.Fatalf("lowerType(%s): %v", in, err)
		}
		if string(got) != want {
			t.Fatalf("lowerType(%s) = %q, want %q", in, got, want)
		}
	}
}

func TestBoxedType_UnsignedUsesJavaBoxedPrimitives(t *testing.T) {
	cases := map[string]string{
		"uint":    "Integer",
		"uint8":   "Short",
		"uint16":  "Integer",
		"uint32":  "Integer",
		"uint64":  "Long",
		"uintptr": "Long",
	}

	for in, want := range cases {
		expr, _ := parser.ParseExpr(in)
		typ, err := lowerType(expr)
		if err != nil {
			t.Fatalf("lowerType(%s): %v", in, err)
		}
		if got := boxedType(typ); got != want {
			t.Fatalf("boxedType(%s) = %q, want %q", in, got, want)
		}
	}
}

func TestUnsignedOpsUseJavaHelpers_AllValueCombinations(t *testing.T) {
	ops := []string{"<", "<=", ">", ">=", "==", "!=", "/", "%"}
	values32 := []string{"1", "2", "0x80000000"}
	values64 := []string{"1", "3", "0x8000000000000000"}

	for _, op := range ops {
		for _, a := range values32 {
			for _, b := range values32 {
				if (op == "/" || op == "%") && b == "0" {
					continue
				}
				assertUnsignedHelperForCombo(t, "uint32", "Integer", op, a, b)
			}
		}
		for _, a := range values64 {
			for _, b := range values64 {
				if (op == "/" || op == "%") && b == "0" {
					continue
				}
				assertUnsignedHelperForCombo(t, "uint64", "Long", op, a, b)
			}
		}
	}
}

func assertUnsignedHelperForCombo(t *testing.T, goType, helper, op, a, b string) {
	t.Helper()
	dir := t.TempDir()
	src := fmt.Sprintf("package main\nimport \"fmt\"\nfunc f(a %s, b %s){fmt.Println(a %s b)}\nfunc main(){f(%s,%s)}\n", goType, goType, op, a, b)
	path := filepath.Join(dir, "main.go")
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	tree, err := syntaxtree.LoadTypedFile(path)
	if err != nil {
		t.Fatal(err)
	}
	ir, err := GoToIntermediate(tree)
	if err != nil {
		t.Fatal(err)
	}
	javaCode, err := IntermediateToJava(ir)
	if err != nil {
		t.Fatal(err)
	}
	want := helper + ".compareUnsigned"
	if op == "/" {
		want = helper + ".divideUnsigned"
	}
	if op == "%" {
		want = helper + ".remainderUnsigned"
	}
	if !strings.Contains(javaCode, want) {
		t.Fatalf("%s %s %s expected %s in generated Java:\n%s", goType, a, op+" "+b, want, javaCode)
	}
}
