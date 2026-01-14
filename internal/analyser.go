package internal

import (
	"github.com/ebukreev/go-z3/z3"
	"golang.org/x/tools/go/ssa"
	"symbolic-execution-course/internal/memory"
	builder "symbolic-execution-course/internal/ssa"
	"symbolic-execution-course/internal/symbolic"
	"symbolic-execution-course/internal/translator"
	"symbolic-execution-course/internal/util"
)

type Analyser struct {
	Package      *ssa.Package
	StatesQueue  PriorityQueue
	PathSelector PathSelector
	Results      []Interpreter
	Z3Translator *translator.Z3Translator
}

func Analyse(source string, functionName string) []Interpreter {
	funBuilder := builder.NewBuilder()
	fun, err := funBuilder.ParseAndBuildSSA(source, functionName)
	if err != nil {
		panic(err)
	}

	config := z3.NewContextConfig()
	ctx := z3.NewContext(config)
	tr := translator.NewZ3Translator(ctx, config)
	defer tr.Close()

	an := Analyser{
		// no need for interproc
		Package:      nil,
		StatesQueue:  make(PriorityQueue, 0),
		Results:      make([]Interpreter, 0),
		PathSelector: &DfsPathSelector{0},
		Z3Translator: tr,
	}

	mem := memory.NewSymbolicMemory()
	locals := make(map[string]symbolic.SymbolicExpression)
	for _, p := range fun.Params {
		locals[p.Name()] = mem.MakeRef(builder.ConvertToSymbolic(p.Type()))
	}
	stack := CallStackFrame{
		Function:    fun,
		LocalMemory: locals,
		Block:       util.FirstOrNil(fun.Blocks),
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

		interpreter := item.value
		for _, instr := range util.Last(interpreter.CallStack).Block.Instrs {
			nextStates := interpreter.interpretDynamically(instr)

			for _, s := range nextStates {
				an.StatesQueue.Push(&Item{
					value:    s,
					priority: an.PathSelector.CalculatePriority(s),
				})
			}
		}

		an.Results = append(an.Results, interpreter)
	}

	return an.Results
}
