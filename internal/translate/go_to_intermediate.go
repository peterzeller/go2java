package translate

import (
	"fmt"
	"go/ast"
	"go/token"

	"go2java/internal/intermediate"
	"go2java/internal/syntaxtree"
)

func GoToIntermediate(tree *syntaxtree.TypedFile) (*intermediate.Program, error) {
	program := &intermediate.Program{PackageName: tree.File.Name.Name}
	for _, decl := range tree.File.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok != token.TYPE {
				continue
			}
			recs, err := lowerTypeDecl(d)
			if err != nil {
				return nil, err
			}
			program.Records = append(program.Records, recs...)
		case *ast.FuncDecl:
			irFn, err := lowerFunc(d)
			if err != nil {
				return nil, err
			}
			program.Functions = append(program.Functions, irFn)
		}
	}
	moveMethodsIntoRecords(program)
	return program, nil
}

func lowerTypeDecl(d *ast.GenDecl) ([]*intermediate.Record, error) {
	var out []*intermediate.Record
	for _, spec := range d.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			continue
		}
		rec := &intermediate.Record{Name: ts.Name.Name}
		for _, field := range st.Fields.List {
			t, err := lowerType(field.Type)
			if err != nil {
				return nil, err
			}
			for _, name := range field.Names {
				rec.Fields = append(rec.Fields, intermediate.Parameter{Name: name.Name, Type: t})
			}
		}
		out = append(out, rec)
	}
	return out, nil
}

func lowerFunc(fn *ast.FuncDecl) (*intermediate.Function, error) {
	irFn := &intermediate.Function{Name: fn.Name.Name, ReturnType: intermediate.TypeVoid}
	if fn.Recv != nil && len(fn.Recv.List) == 1 {
		recv := fn.Recv.List[0]
		t, err := lowerType(recv.Type)
		if err != nil {
			return nil, err
		}
		irFn.ReceiverType = t
		if len(recv.Names) == 1 {
			irFn.ReceiverName = recv.Names[0].Name
		}
	}
	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			typ, err := lowerType(field.Type)
			if err != nil {
				return nil, err
			}
			for _, name := range field.Names {
				irFn.Parameters = append(irFn.Parameters, intermediate.Parameter{Name: name.Name, Type: typ})
			}
		}
	}
	if fn.Type.Results != nil {
		if len(fn.Type.Results.List) != 1 {
			return nil, fmt.Errorf("function %s: only single return value supported", fn.Name.Name)
		}
		typ, err := lowerType(fn.Type.Results.List[0].Type)
		if err != nil {
			return nil, err
		}
		irFn.ReturnType = typ
	}
	for _, st := range fn.Body.List {
		s, err := lowerStmt(st)
		if err != nil {
			return nil, err
		}
		irFn.Body = append(irFn.Body, s)
	}
	return irFn, nil
}

func lowerType(expr ast.Expr) (intermediate.Type, error) {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return "", fmt.Errorf("unsupported type %T", expr)
	}
	switch id.Name {
	case "int":
		return intermediate.TypeInt, nil
	case "string":
		return intermediate.TypeString, nil
	case "bool":
		return intermediate.TypeBool, nil
	default:
		return intermediate.Type(id.Name), nil
	}
}

func lowerStmt(st ast.Stmt) (intermediate.Stmt, error) {
	// unchanged sections omitted for brevity in this file rewrite
	switch s := st.(type) {
	case *ast.ExprStmt:
		ex, err := lowerExpr(s.X)
		if err != nil {
			return nil, err
		}
		return &intermediate.ExprStmt{Expr: ex}, nil
	case *ast.AssignStmt:
		if len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			return nil, fmt.Errorf("only single assignment supported")
		}
		lhs, ok := s.Lhs[0].(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("unsupported lhs %T", s.Lhs[0])
		}
		rhs, err := lowerExpr(s.Rhs[0])
		if err != nil {
			return nil, err
		}
		typ := intermediate.Type("")
		if s.Tok == token.DEFINE {
			typ = "var"
		}
		return &intermediate.AssignStmt{Name: lhs.Name, Type: typ, Value: rhs}, nil
	case *ast.ReturnStmt:
		if len(s.Results) == 0 {
			return &intermediate.ReturnStmt{}, nil
		}
		ex, err := lowerExpr(s.Results[0])
		if err != nil {
			return nil, err
		}
		return &intermediate.ReturnStmt{Value: ex}, nil
	case *ast.IfStmt:
		cond, err := lowerExpr(s.Cond)
		if err != nil {
			return nil, err
		}
		thenStmts, err := lowerBlock(s.Body.List)
		if err != nil {
			return nil, err
		}
		var elseStmts []intermediate.Stmt
		if s.Else != nil {
			es, ok := s.Else.(*ast.BlockStmt)
			if !ok {
				return nil, fmt.Errorf("else-if unsupported")
			}
			elseStmts, err = lowerBlock(es.List)
			if err != nil {
				return nil, err
			}
		}
		return &intermediate.IfStmt{Cond: cond, Then: thenStmts, Else: elseStmts}, nil
	case *ast.ForStmt:
		if s.Init != nil || s.Post != nil {
			return nil, fmt.Errorf("for init/post unsupported")
		}
		cond, err := lowerExpr(s.Cond)
		if err != nil {
			return nil, err
		}
		body, err := lowerBlock(s.Body.List)
		if err != nil {
			return nil, err
		}
		return &intermediate.ForStmt{Cond: cond, Body: body}, nil
	default:
		return nil, fmt.Errorf("unsupported statement %T", st)
	}
}

func lowerBlock(in []ast.Stmt) ([]intermediate.Stmt, error) {
	out := make([]intermediate.Stmt, 0, len(in))
	for _, st := range in {
		s, err := lowerStmt(st)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func lowerExpr(ex ast.Expr) (intermediate.Expr, error) {
	switch e := ex.(type) {
	case *ast.BasicLit:
		switch e.Kind {
		case token.INT:
			return &intermediate.IntLiteral{Value: e.Value}, nil
		case token.STRING:
			return &intermediate.StringLiteral{Value: e.Value}, nil
		}
	case *ast.Ident:
		if e.Name == "true" {
			return &intermediate.BoolLiteral{Value: true}, nil
		}
		if e.Name == "false" {
			return &intermediate.BoolLiteral{Value: false}, nil
		}
		return &intermediate.IdentExpr{Name: e.Name}, nil
	case *ast.BinaryExpr:
		l, err := lowerExpr(e.X)
		if err != nil {
			return nil, err
		}
		r, err := lowerExpr(e.Y)
		if err != nil {
			return nil, err
		}
		return &intermediate.BinaryExpr{Op: e.Op.String(), Left: l, Right: r}, nil
	case *ast.CallExpr:
		fn, err := lowerExpr(e.Fun)
		if err != nil {
			return nil, err
		}
		args := make([]intermediate.Expr, 0, len(e.Args))
		for _, a := range e.Args {
			aa, err := lowerExpr(a)
			if err != nil {
				return nil, err
			}
			args = append(args, aa)
		}
		return &intermediate.CallExpr{Func: fn, Args: args}, nil
	case *ast.SelectorExpr:
		t, err := lowerExpr(e.X)
		if err != nil {
			return nil, err
		}
		return &intermediate.SelectorExpr{Target: t, Field: e.Sel.Name}, nil
	case *ast.CompositeLit:
		name, ok := e.Type.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("unsupported composite type %T", e.Type)
		}
		args := make([]intermediate.Expr, 0, len(e.Elts))
		for _, elt := range e.Elts {
			expr := elt
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				expr = kv.Value
			}
			v, err := lowerExpr(expr)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		}
		return &intermediate.CompositeLiteral{TypeName: name.Name, Args: args}, nil
	}
	return nil, fmt.Errorf("unsupported expression %T", ex)
}

func moveMethodsIntoRecords(p *intermediate.Program) {
	byName := map[string]*intermediate.Record{}
	for _, r := range p.Records {
		byName[r.Name] = r
	}
	var fns []*intermediate.Function
	for _, fn := range p.Functions {
		if fn.ReceiverType != "" {
			if rec := byName[string(fn.ReceiverType)]; rec != nil {
				rec.Methods = append(rec.Methods, fn)
				continue
			}
		}
		fns = append(fns, fn)
	}
	p.Functions = fns
}
