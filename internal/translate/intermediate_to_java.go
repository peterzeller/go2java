package translate

import (
	"fmt"
	"strings"
	"unicode"

	"go2java/internal/intermediate"
)

type renderCtx struct {
	varIsArrayList map[string]bool
}

func IntermediateToJava(p *intermediate.Program) (string, error) {
	var b strings.Builder
	if len(p.Records) == 1 {
		return renderClass(p.Records[0], p.Functions)
	}
	b.WriteString("public class Main {\n")
	b.WriteString("\n    private static String goFmt(Object o) {\n")
	b.WriteString("        if (o instanceof java.util.List<?> l) return l.toString().replace(\", \", \" \");\n")
	b.WriteString("        return String.valueOf(o);\n")
	b.WriteString("    }\n")
	for _, fn := range p.Functions {
		_ = writeFunction(&b, fn, true, nil)
	}
	b.WriteString("}\n")
	return b.String(), nil
}

func renderClass(rec *intermediate.Record, extra []*intermediate.Function) (string, error) {
	var b strings.Builder
	b.WriteString("public class " + rec.Name + " {\n")
	b.WriteString("\n    private static String goFmt(Object o) {\n")
	b.WriteString("        if (o instanceof java.util.List<?> l) return l.toString().replace(\", \", \" \");\n")
	b.WriteString("        return String.valueOf(o);\n")
	b.WriteString("    }\n")
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
	ctx := &renderCtx{varIsArrayList: map[string]bool{}}
	for _, p := range fn.Parameters {
		if strings.HasPrefix(string(p.Type), "List<") {
			ctx.varIsArrayList[p.Name] = false
		}
	}
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
		line, _ := stmtToJava(st, fn.ReceiverName, rec, ctx)
		b.WriteString("        " + line + "\n")
	}
	b.WriteString("    }\n")
	return nil
}

func stmtToJava(st intermediate.Stmt, receiverName string, rec *intermediate.Record, ctx *renderCtx) (string, error) {
	switch s := st.(type) {
	case *intermediate.ExprStmt:
		e, _ := exprToJava(s.Expr, receiverName, rec, ctx)
		return e + ";", nil
	case *intermediate.AssignStmt:
		if call, ok := s.Value.(*intermediate.CallExpr); ok {
			if id, ok := call.Func.(*intermediate.IdentExpr); ok && id.Name == "append" && len(call.Args) == 2 {
				base, _ := exprToJava(call.Args[0], receiverName, rec, ctx)
				val, _ := exprToJava(call.Args[1], receiverName, rec, ctx)
				if b, ok := call.Args[0].(*intermediate.IdentExpr); ok && ctx.varIsArrayList[b.Name] {
					ctx.varIsArrayList[s.Name] = true
					if s.Name == b.Name {
						return fmt.Sprintf("%s.add(%s);", base, val), nil
					}
					prefix := ""
					if s.Type != "" {
						prefix = string(s.Type) + " "
					}
					return fmt.Sprintf("%s%s = %s; %s.add(%s);", prefix, s.Name, base, s.Name, val), nil
				}
				ctx.varIsArrayList[s.Name] = true
				prefix := ""
				if s.Type != "" {
					prefix = string(s.Type) + " "
				}
				return fmt.Sprintf("%s%s = new java.util.ArrayList<>(%s); %s.add(%s);", prefix, s.Name, base, s.Name, val), nil
			}
		}
		e, _ := exprToJava(s.Value, receiverName, rec, ctx)
		if s.KnownArrayList {
			ctx.varIsArrayList[s.Name] = true
		} else if s.MaybeNonArrayList {
			ctx.varIsArrayList[s.Name] = false
		}
		if s.Type != "" {
			return fmt.Sprintf("%s %s = %s;", s.Type, s.Name, e), nil
		}
		return fmt.Sprintf("%s = %s;", s.Name, e), nil
	case *intermediate.FieldAssignStmt:
		v, _ := exprToJava(s.Value, receiverName, rec, ctx)
		if id, ok := s.Target.(*intermediate.IdentExpr); ok && id.Name == receiverName {
			return "this." + s.Field + " = " + v + ";", nil
		}
		t, _ := exprToJava(s.Target, receiverName, rec, ctx)
		if isExported(s.Field) {
			return fmt.Sprintf("%s.set%s(%s);", t, title(s.Field), v), nil
		}
		return fmt.Sprintf("%s.%s = %s;", t, s.Field, v), nil
	case *intermediate.ReturnStmt:
		if s.Value == nil {
			return "return;", nil
		}
		e, _ := exprToJava(s.Value, receiverName, rec, ctx)
		return "return " + e + ";", nil
	default:
		return "", fmt.Errorf("unsupported")
	}
}

func exprToJava(ex intermediate.Expr, receiverName string, rec *intermediate.Record, ctx *renderCtx) (string, error) {
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
		t, _ := exprToJava(e.Target, receiverName, rec, ctx)
		if isExported(e.Field) {
			return t + ".get" + title(e.Field) + "()", nil
		}
		return t + "." + e.Field, nil
	case *intermediate.BinaryExpr:
		l, _ := exprToJava(e.Left, receiverName, rec, ctx)
		r, _ := exprToJava(e.Right, receiverName, rec, ctx)
		return "(" + l + " " + e.Op + " " + r + ")", nil
	case *intermediate.CompositeLiteral:
		args := []string{}
		for _, a := range e.Args {
			aa, _ := exprToJava(a, receiverName, rec, ctx)
			args = append(args, aa)
		}
		if strings.HasPrefix(e.TypeName, "[]") {
			return "new java.util.ArrayList<>(java.util.List.of(" + strings.Join(args, ", ") + "))", nil
		}
		return "new " + e.TypeName + "(" + strings.Join(args, ", ") + ")", nil
	case *intermediate.CallExpr:
		n := ""
		if sel, ok := e.Func.(*intermediate.SelectorExpr); ok {
			if id, ok := sel.Target.(*intermediate.IdentExpr); ok && id.Name == "fmt" && sel.Field == "Println" {
				if len(e.Args) == 1 {
					aa, _ := exprToJava(e.Args[0], receiverName, rec, ctx)
					return "System.out.println(goFmt(" + aa + "))", nil
				}
				n = "System.out.println"
			}
		}
		if n == "" {
			n, _ = exprToJava(e.Func, receiverName, rec, ctx)
		}
		if id, ok := e.Func.(*intermediate.IdentExpr); ok && id.Name == "append" {
			return "/* append only supported in assignment */", nil
		}
		args := []string{}
		for _, a := range e.Args {
			aa, _ := exprToJava(a, receiverName, rec, ctx)
			args = append(args, aa)
		}
		return n + "(" + strings.Join(args, ", ") + ")", nil
	case *intermediate.SliceExpr:
		t, _ := exprToJava(e.Target, receiverName, rec, ctx)
		low := "0"
		high := t + ".size()"
		if e.Low != nil {
			low, _ = exprToJava(e.Low, receiverName, rec, ctx)
		}
		if e.High != nil {
			high, _ = exprToJava(e.High, receiverName, rec, ctx)
		}
		return t + ".subList(" + low + ", " + high + ")", nil
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
