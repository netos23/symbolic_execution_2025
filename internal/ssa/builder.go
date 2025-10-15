// Package ssa предоставляет функции для построения SSA представления
package ssa

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// Builder отвечает за построение SSA из исходного кода Go
type Builder struct {
	fset *token.FileSet
}

// NewBuilder создаёт новый экземпляр Builder
func NewBuilder() *Builder {
	return &Builder{
		fset: token.NewFileSet(),
	}
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// ParseAndBuildSSA парсит исходный код Go и создаёт SSA представление
// Возвращает SSA программу и функцию по имени
func (b *Builder) ParseAndBuildSSA(source string, funcName string) (*ssa.Function, error) {
	file, err := parser.ParseFile(b.fset, "main.go", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	files := []*ast.File{file}

	pkg := types.NewPackage("homework1/main.go", "main")
	tc := &types.Config{Importer: importer.Default()}

	ssa, _, err := ssautil.BuildPackage(
		tc, b.fset, pkg, files, ssa.SanityCheckFunctions,
	)
	if err != nil {
		return nil, err
	}

	if fun := ssa.Func(funcName); fun != nil {
		return fun, nil
	}

	return nil, fmt.Errorf("missing function: %s", funcName)
}
