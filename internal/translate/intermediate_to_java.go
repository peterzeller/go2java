package translate

import (
	"fmt"
	"strings"
	"unicode"

	"go2java/internal/intermediate"
)

func IntermediateToJava(p *intermediate.Program) (string, error) {
	var b strings.Builder
	if len(p.Records) == 1 {
		return renderClass(p.Records[0], p.Functions)
	}
	b.WriteString("public class Main {\n")
	for _, fn := range p.Functions {
		_ = writeFunction(&b, fn, true, nil)
	}
	b.WriteString("}\n")
	return b.String(), nil
}

func renderClass(rec *intermediate.Record, extra []*intermediate.Function) (string, error) {
	var b strings.Builder
	b.WriteString("public class " + rec.Name + " {\n")
	for _, f := range rec.Fields {
		b.WriteString("    private " + string(f.Type) + " " + f.Name + ";\n")
	}
	b.WriteString("\n    public " + rec.Name + "(")
	for i, f := range rec.Fields {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(string(f.Type) + " " + f.Name)
	}
	b.WriteString(") {\n")
	for _, f := range rec.Fields {
		b.WriteString("        this." + f.Name + " = " + f.Name + ";\n")
	}
	b.WriteString("    }\n")
	for _, f := range rec.Fields {
		if f.Public {
			cap := title(f.Name)
			b.WriteString("\n    public " + string(f.Type) + " get" + cap + "() { return this." + f.Name + "; }\n")
			b.WriteString("    public void set" + cap + "(" + string(f.Type) + " " + f.Name + ") { this." + f.Name + " = " + f.Name + "; }\n")
		}
	}
	for _, m := range rec.Methods {
		_ = writeFunction(&b, m, false, rec)
	}
	for _, fn := range extra {
		_ = writeFunction(&b, fn, true, rec)
	}
	b.WriteString("}\n")
	return b.String(), nil
}

func writeFunction(b *strings.Builder, fn *intermediate.Function, static bool, rec *intermediate.Record) error {
	pref := ""
	if static {
		pref = "static "
	}
	b.WriteString("\n    public " + pref + string(fn.ReturnType) + " " + fn.Name + "(")
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
		line, _ := stmtToJava(st, fn.ReceiverName, rec)
		b.WriteString("        " + line + "\n")
	}
	b.WriteString("    }\n")
	return nil
}

func stmtToJava(st intermediate.Stmt, receiverName string, rec *intermediate.Record) (string, error) {
	switch s := st.(type) {
	case *intermediate.ExprStmt:
		e, _ := exprToJava(s.Expr, receiverName, rec)
		return e + ";", nil
	case *intermediate.AssignStmt:
		e, _ := exprToJava(s.Value, receiverName, rec)
		if s.Type != "" {
			return fmt.Sprintf("%s %s = %s;", s.Type, s.Name, e), nil
		}
		return fmt.Sprintf("%s = %s;", s.Name, e), nil
	case *intermediate.FieldAssignStmt:
		v, _ := exprToJava(s.Value, receiverName, rec)
		if id, ok := s.Target.(*intermediate.IdentExpr); ok && id.Name == receiverName {
			return "this." + s.Field + " = " + v + ";", nil
		}
		t, _ := exprToJava(s.Target, receiverName, rec)
		if isExported(s.Field) {
			return fmt.Sprintf("%s.set%s(%s);", t, title(s.Field), v), nil
		}
		return fmt.Sprintf("%s.%s = %s;", t, s.Field, v), nil
	case *intermediate.ReturnStmt:
		if s.Value == nil {
			return "return;", nil
		}
		e, _ := exprToJava(s.Value, receiverName, rec)
		return "return " + e + ";", nil
	default:
		return "", fmt.Errorf("unsupported")
	}
}

func exprToJava(ex intermediate.Expr, receiverName string, rec *intermediate.Record) (string, error) {
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
	case *intermediate.SelectorExpr:
		if id, ok := e.Target.(*intermediate.IdentExpr); ok && id.Name == receiverName {
			return "this." + e.Field, nil
		}
		t, _ := exprToJava(e.Target, receiverName, rec)
		if isExported(e.Field) {
			return t + ".get" + title(e.Field) + "()", nil
		}
		return t + "." + e.Field, nil
	case *intermediate.BinaryExpr:
		l, _ := exprToJava(e.Left, receiverName, rec)
		r, _ := exprToJava(e.Right, receiverName, rec)
		return "(" + l + " " + e.Op + " " + r + ")", nil
	case *intermediate.CompositeLiteral:
		args := []string{}
		for _, a := range e.Args {
			aa, _ := exprToJava(a, receiverName, rec)
			args = append(args, aa)
		}
		return "new " + e.TypeName + "(" + strings.Join(args, ", ") + ")", nil
	case *intermediate.CallExpr:
		n := ""
		if sel, ok := e.Func.(*intermediate.SelectorExpr); ok {
			if id, ok := sel.Target.(*intermediate.IdentExpr); ok && id.Name == "fmt" && sel.Field == "Println" {
				n = "System.out.println"
			}
		}
		if n == "" {
			n, _ = exprToJava(e.Func, receiverName, rec)
		}
		args := []string{}
		for _, a := range e.Args {
			aa, _ := exprToJava(a, receiverName, rec)
			args = append(args, aa)
		}
		return n + "(" + strings.Join(args, ", ") + ")", nil
	}
	return "", fmt.Errorf("unsupported")
}

func isExported(name string) bool {
	if name == "" {
		return false
	}
	return unicode.IsUpper([]rune(name)[0])
}
func title(name string) string {
	if name == "" {
		return ""
	}
	r := []rune(name)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
