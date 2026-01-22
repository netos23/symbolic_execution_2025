package main

import (
	"fmt"
	"os"
	"symbolic-execution-course/internal"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <filename> <func> [-v]")
		os.Exit(1)
	}

	filename := os.Args[1]
	function := os.Args[2]
	verbose := false
	if len(os.Args) > 3 {
		verbose = true
	}

	result := internal.AnalyseSource(filename, function, verbose)
	if verbose {
		for _, interpreter := range result {
			fmt.Println(interpreter)
		}
	}

}
