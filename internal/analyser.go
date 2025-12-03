package internal

import (
	"golang.org/x/tools/go/ssa"
	"symbolic-execution-course/internal/memory"
	builder "symbolic-execution-course/internal/ssa"
	"symbolic-execution-course/internal/symbolic"
	"symbolic-execution-course/internal/translator"
)

type Analyser struct {
	Package      *ssa.Package
	StatesQueue  PriorityQueue
	PathSelector PathSelector
	Results      []Interpreter
	Z3Translator *translator.Z3Translator
}

func Analyse(source string, functionName string) []Interpreter {
	builder := builder.NewBuilder()
	fun, err := builder.ParseAndBuildSSA(source, functionName)
	if err != nil {
		panic(err)
	}

	tr := translator.NewZ3Translator()
	defer tr.Close()

	an := Analyser{
		// no need for interproc
		Package:      nil,
		StatesQueue:  make(PriorityQueue, 0),
		Results:      make([]Interpreter, 0),
		Z3Translator: tr,
	}

	mem := memory.NewSymbolicMemory()
	locals := make(map[string]symbolic.SymbolicExpression)
	for _, p := range fun.Params {
		locals[p.Name()] = mem.MakeRef(paramToSymbolic(p))
	}
	stack := CallStackFrame{
		Function:    fun,
		LocalMemory: locals,
	}
	init := Interpreter{
		CallStack: []CallStackFrame{stack},
		Analyser:  &an,
		Heap:      mem,
	}

	an.StatesQueue.Push(&Item{
		value:    init,
		priority: an.PathSelector.CalculatePriority(init),
	})

	for an.StatesQueue.Len() != 0 {
		item := an.StatesQueue.Pop().(*Item)

	}
}

func paramToSymbolic(*ssa.Parameter) (symbolic.ExpressionType, string, *symbolic.GenericType) {

}
