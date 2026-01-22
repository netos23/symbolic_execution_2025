package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type TestFile struct {
	File      string   `json:"file"`
	Functions []string `json:"functions"`
}

type TestSuite struct {
	Tests []TestFile `json:"tests"`
}

func main() {
	root := "final_tests"
	suite := testSuite(root)

	clean()
	compile()
	suc, fail := 0, 0
	for _, t := range suite.Tests {
		sucI, failI := test(t)
		suc += sucI
		fail += failI
	}

	fmt.Println()
	if fail == 0 {
		fmt.Printf("🎉 %d tests passed.\n", suc)
	} else {
		fmt.Printf("::error::%d test passed, %d failed.\n", suc, fail)
		os.Exit(1)
	}

	defer clean()
}

func testSuite(root string) TestSuite {
	var suite TestSuite

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}

		var funcs []string
		for _, decl := range node.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil {
				funcs = append(funcs, fn.Name.Name)
			}
		}
		if len(funcs) > 0 {
			suite.Tests = append(suite.Tests, TestFile{
				File:      path,
				Functions: funcs,
			})
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return suite
}

func compile() {
	cmd := exec.Command("go", "build", "-buildvcs=false")
	cmd.Dir = "homework5"
	if runtime.GOOS == "darwin" {
		cmd.Env = append(os.Environ(),
			"CGO_CFLAGS=-I/opt/homebrew/Cellar/z3/4.15.4/include -O0 -g",
			"CGO_LDFLAGS=-L/opt/homebrew/Cellar/z3/4.15.4/lib -lz3",
		)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func test(file TestFile) (int, int) {
	fmt.Println("------------------------")
	fmt.Printf("Run go test %s:\n", file.File)

	suc, fail := 0, 0

	for _, test := range file.Functions {
		//fmt.Println("------------------------")
		//fmt.Printf("Run test: %s\n", test)

		//if file.File != "final_tests/structs.go" {
		//	continue
		//}

		err := runCompiledGo("homework5/homework5", "final_tests", test)

		if err == nil {
			fmt.Printf("✅ %s: %s\n", file.File, test)
			suc++
		} else {
			fmt.Printf("::group::❌ %s: %s (failed)\n", file.File, test)
			fmt.Println()
			fmt.Println(err)
			fmt.Println("::endgroup::")
			fail++
		}
	}
	return suc, fail
}

func runCompiledGo(binPath, fileArg, funcArg string) error {
	var cmd *exec.Cmd
	if len(os.Args) > 1 {
		cmd = exec.Command(binPath, fileArg, funcArg, "-v")
	} else {
		cmd = exec.Command(binPath, fileArg, funcArg)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running %s: %w", binPath, err)
	}
	return nil
}

func clean() {
	err := os.Remove("homework5/homework5")
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error removing binary: %v\n", err)
		os.Exit(1)
	}
}
