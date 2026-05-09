package syntaxtree

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

// TypedFile contains the parsed Go AST plus type information for one source file.
type TypedFile struct {
	Path      string
	FileSet   *token.FileSet
	File      *ast.File
	Package   *types.Package
	TypesInfo *types.Info
}

// LoadTypedFile loads, parses, and type-checks the Go package containing path,
// then returns the AST and type information for that specific file.
func LoadTypedFile(path string) (*TypedFile, error) {
	if path == "" {
		return nil, errors.New("go file path is required")
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes |
			packages.NeedModule,
		Dir: filepath.Dir(absolutePath),
	}

	pkgs, err := packages.Load(cfg, absolutePath)
	if err != nil {
		return nil, fmt.Errorf("load package: %w", err)
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no package found for %s", path)
	}

	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return nil, formatPackageErrors(path, pkg.Errors)
	}

	for _, file := range pkg.Syntax {
		if sameFile(pkg.Fset, file, absolutePath) {
			return &TypedFile{
				Path:      absolutePath,
				FileSet:   pkg.Fset,
				File:      file,
				Package:   pkg.Types,
				TypesInfo: pkg.TypesInfo,
			}, nil
		}
	}

	return nil, fmt.Errorf("typed syntax tree not found for %s", path)
}

func sameFile(fset *token.FileSet, file *ast.File, path string) bool {
	position := fset.Position(file.Pos())
	if position.Filename == "" {
		return false
	}

	filePath, err := filepath.Abs(position.Filename)
	if err != nil {
		return false
	}

	return filepath.Clean(filePath) == filepath.Clean(path)
}

func formatPackageErrors(path string, errs []packages.Error) error {
	message := fmt.Sprintf("type-check %s", path)
	for _, err := range errs {
		message += "\n  " + err.Error()
	}
	return errors.New(message)
}
