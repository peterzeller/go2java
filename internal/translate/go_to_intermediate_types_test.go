package translate

import (
	"go/parser"
	"testing"
)

func TestLowerType_AllGoIntegerTypes(t *testing.T) {
	tests := map[string]string{
		"int":     "int",
		"int8":    "byte",
		"int16":   "short",
		"int32":   "int",
		"int64":   "long",
		"uint":    "int",
		"uint8":   "byte",
		"uint16":  "short",
		"uint32":  "int",
		"uint64":  "long",
		"uintptr": "long",
		"byte":    "byte",
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
		"uint8":   "Byte",
		"uint16":  "Short",
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
