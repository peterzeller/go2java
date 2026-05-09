package translate

import (
	"fmt"
	"strings"

	"go2java/internal/intermediate"
)

func IntermediateToJava(p *intermediate.Program) (string, error) {
	var b strings.Builder
	b.WriteString("public class Main {\n")
	for _, fn := range p.Functions {
		if err := writeFunction(&b, fn); err != nil {
			return "", err
		}
	}
	b.WriteString("}\n")
	return b.String(), nil
}

func writeFunction(b *strings.Builder, fn *intermediate.Function) error {
	static := "static "
	b.WriteString("    public " + static + string(fn.ReturnType) + " " + fn.Name + "(")
	if fn.Name == "main" && fn.ReturnType == intermediate.TypeVoid && len(fn.Parameters) == 0 {
		b.WriteString("String[] args")
	} else {
		for i, p := range fn.Parameters {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(string(p.Type) + " " + p.Name)
		}
	}
	b.WriteString(") {\n")
	for _, st := range fn.Body {
		s, err := stmtToJava(st)
		if err != nil {
			return err
		}
		for _, line := range strings.Split(strings.TrimSuffix(s, "\n"), "\n") {
			b.WriteString("        " + line + "\n")
		}
	}
	b.WriteString("    }\n")
	return nil
}

func stmtToJava(st intermediate.Stmt) (string, error) {
	switch s := st.(type) {
	case *intermediate.ExprStmt:
		e, err := exprToJava(s.Expr)
		if err != nil {
			return "", err
		}
		return e + ";", nil
	case *intermediate.AssignStmt:
		e, err := exprToJava(s.Value)
		if err != nil {
			return "", err
		}
		if s.Type != "" {
			return fmt.Sprintf("%s %s = %s;", s.Type, s.Name, e), nil
		}
		return fmt.Sprintf("%s = %s;", s.Name, e), nil
	case *intermediate.ReturnStmt:
		if s.Value == nil {
			return "return;", nil
		}
		e, err := exprToJava(s.Value)
		if err != nil {
			return "", err
		}
		return "return " + e + ";", nil
	case *intermediate.IfStmt:
		c, err := exprToJava(s.Cond)
		if err != nil {
			return "", err
		}
		var b strings.Builder
		b.WriteString("if (" + c + ") {\n")
		for _, st := range s.Then {
			line, err := stmtToJava(st)
			if err != nil {
				return "", err
			}
			b.WriteString("    " + line + "\n")
		}
		b.WriteString("}")
		if len(s.Else) > 0 {
			b.WriteString(" else {\n")
			for _, st := range s.Else {
				line, err := stmtToJava(st)
				if err != nil {
					return "", err
				}
				b.WriteString("    " + line + "\n")
			}
			b.WriteString("}")
		}
		return b.String(), nil
	case *intermediate.ForStmt:
		c, err := exprToJava(s.Cond)
		if err != nil {
			return "", err
		}
		var b strings.Builder
		b.WriteString("while (" + c + ") {\n")
		for _, st := range s.Body {
			line, err := stmtToJava(st)
			if err != nil {
				return "", err
			}
			b.WriteString("    " + line + "\n")
		}
		b.WriteString("}")
		return b.String(), nil
	default:
		return "", fmt.Errorf("unsupported stmt %T", st)
	}
}

func exprToJava(ex intermediate.Expr) (string, error) {
	switch e := ex.(type) {
	case *intermediate.IntLiteral:
		return e.Value, nil
	case *intermediate.StringLiteral:
		return e.Value, nil
	case *intermediate.BoolLiteral:
		if e.Value {
			return "true", nil
		}
		return "false", nil
	case *intermediate.IdentExpr:
		return e.Name, nil
	case *intermediate.BinaryExpr:
		l, _ := exprToJava(e.Left)
		r, _ := exprToJava(e.Right)
		return fmt.Sprintf("(%s %s %s)", l, e.Op, r), nil
	case *intermediate.CallExpr:
		name := e.Func
		if name == "fmt.Println" {
			name = "System.out.println"
		}
		args := make([]string, 0, len(e.Args))
		for _, a := range e.Args {
			aa, err := exprToJava(a)
			if err != nil {
				return "", err
			}
			args = append(args, aa)
		}
		return fmt.Sprintf("%s(%s)", name, strings.Join(args, ", ")), nil
	default:
		return "", fmt.Errorf("unsupported expr %T", ex)
	}
}
