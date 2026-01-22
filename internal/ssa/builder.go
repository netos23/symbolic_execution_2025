// Package ssa предоставляет функции для построения SSA представления
package ssa

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"symbolic-execution-course/internal/util"

	"golang.org/x/tools/go/ast/astutil"
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
	file, err := parser.ParseFile(b.fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return b.ParseAndBuildSSAFiles([]*ast.File{file}, funcName)
}

func (b *Builder) ParseAndBuildSSAPackage(dir, funcName string) (*ssa.Function, error) {
	files := make([]*ast.File, 0)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		f, err := parser.ParseFile(b.fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		files = append(files, f)
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return b.ParseAndBuildSSAFiles(files, funcName)
}

func (b *Builder) ParseAndBuildSSAFiles(files []*ast.File, funcName string) (*ssa.Function, error) {

	// Здесь по хорошему нужен fix point или более хорошая рекурсивная реализация но в текущем случае достаточно и двух раз
	unrolledFiles := util.Convert(files, func(file *ast.File) *ast.File {
		return astutil.Apply(file, UnrollLoops, nil).(*ast.File)
	})
	unrolledFiles = util.Convert(unrolledFiles, func(file *ast.File) *ast.File {
		return astutil.Apply(file, UnrollLoops, nil).(*ast.File)
	})

	// Печать файлов с раскручеными циклами
	//b.DumpFiles(unrolledFiles)

	file := util.FirstOrNil(files)

	pkg := types.NewPackage("homework1/main.go", file.Name.Name)
	tc := &types.Config{Importer: importer.Default()}

	ssa, _, err := ssautil.BuildPackage(
		tc, b.fset, pkg, unrolledFiles, ssa.SanityCheckFunctions,
	)
	if err != nil {
		return nil, err
	}

	if fun := ssa.Func(funcName); fun != nil {
		return fun, nil
	}

	return nil, fmt.Errorf("missing function: %s", funcName)
}

func (b *Builder) DumpFiles(unrolledFiles []*ast.File) {
	for i, fl := range unrolledFiles {
		name := fmt.Sprintf("opt%d.go", i)
		f, err := os.Create(name)
		if err != nil {
			panic(err)
		}
		err = printer.Fprint(f, b.fset, fl)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func (b *Builder) PrintBlocksAndInstructions(fun *ssa.Function) {
	fmt.Println("---------------")
	fmt.Printf("Function: %s\n", fun.Name())
	for i, block := range fun.Blocks {
		b.PrintBlock(i, block)
	}
	fmt.Println("---------------")
}

func (b *Builder) PrintBlock(i int, block *ssa.BasicBlock) {
	fmt.Printf("  Block %d (%s):\n", i, block.Comment)
	for j, instr := range block.Instrs {

		if named, ok := instr.(interface{ Name() string }); ok {
			fmt.Printf("    %d: %s = %s\n", j, named.Name(), instr.String())
		} else {
			fmt.Printf("    %d: %s\n", j, instr.String())
		}

	}
}
