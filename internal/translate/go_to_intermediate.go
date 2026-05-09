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
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		irFn, err := lowerFunc(fn)
		if err != nil {
			return nil, err
		}
		program.Functions = append(program.Functions, irFn)
	}
	return program, nil
}

func lowerFunc(fn *ast.FuncDecl) (*intermediate.Function, error) {
	irFn := &intermediate.Function{Name: fn.Name.Name, ReturnType: intermediate.TypeVoid}
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
		return "", fmt.Errorf("unsupported type %q", id.Name)
	}
}

func lowerStmt(st ast.Stmt) (intermediate.Stmt, error) {
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
			typ = intermediate.TypeInt
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
		name, err := callName(e.Fun)
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
		return &intermediate.CallExpr{Func: name, Args: args}, nil
	}
	return nil, fmt.Errorf("unsupported expression %T", ex)
}

func callName(fun ast.Expr) (string, error) {
	switch f := fun.(type) {
	case *ast.Ident:
		return f.Name, nil
	case *ast.SelectorExpr:
		x, ok := f.X.(*ast.Ident)
		if !ok {
			return "", fmt.Errorf("unsupported selector")
		}
		return x.Name + "." + f.Sel.Name, nil
	default:
		return "", fmt.Errorf("unsupported call target %T", fun)
	}
}
