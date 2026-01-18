package main

import (
	"fmt"
	"symbolic-execution-course/internal"
)

func main() {
	source := `
package main

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func CharSizeAndIndex(a []rune, x rune) byte {
	if a == nil || len(a) <= int(x) || x < 1 {
		return 255
	}
	b := make([]rune, x)
	b[0] = 5
	a[x] = x
	if b[0]+a[x] > 7 {
		return 1
	}
	return 0
}
`

	result := internal.Analyse(source, "CharSizeAndIndex")
	for _, interpreter := range result {
		fmt.Println(interpreter)
	}
}
