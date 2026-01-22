package internal

import (
	"fmt"
	"github.com/ebukreev/go-z3/z3"
	"golang.org/x/tools/go/ssa"
	"symbolic-execution-course/internal/memory"
	builder "symbolic-execution-course/internal/ssa"
	"symbolic-execution-course/internal/symbolic"
	"symbolic-execution-course/internal/translator"
	"symbolic-execution-course/internal/util"
	"symbolic-execution-course/pkg/z3wrapper"
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

	funBuilder.PrintBlocksAndInstructions(fun)

	config := z3.NewContextConfig()
	ctx := z3.NewContext(config)
	tr := translator.NewZ3Translator(ctx, config)
	solver := z3wrapper.NewSolver(ctx)
	defer solver.Close()
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
		Function:           fun,
		LocalMemory:        locals,
		Block:              util.FirstOrNil(fun.Blocks),
		InstructionPointer: 0,
	}
	init := Interpreter{
		CallStack:      []CallStackFrame{stack},
		Analyser:       &an,
		Heap:           mem,
		RecursionLimit: 10,
	}

	an.StatesQueue.Push(&Item{
		value:    init,
		priority: an.PathSelector.CalculatePriority(init),
	})

	for an.StatesQueue.Len() != 0 {
		item := an.StatesQueue.Pop().(*Item)

		interpreter := item.value
		stackFrame := util.Last(interpreter.CallStack)
		instructions := stackFrame.Block.Instrs[stackFrame.InstructionPointer:]

		fmt.Println("Visit block:")

		reachable := true
		if interpreter.PathCondition != nil {
			solver.Push()
			expression, _ := tr.TranslateExpression(interpreter.PathCondition)
			solver.Assert(expression.(z3.Bool))
			sat, _ := solver.Check()
			reachable = sat
			solver.Pop()

			fmt.Printf("Path conditions: %s\n", interpreter.PathCondition.String())
			fmt.Printf("Path SAT: %b\n", reachable)
			fmt.Println("---------------------------")
		}

		funBuilder.PrintBlock(util.IndexOf(stackFrame.Function.Blocks, stackFrame.Block), stackFrame.Block)

		fmt.Println("---------------------------")
		for _, instr := range instructions {
			nextStates := interpreter.interpretDynamically(instr)

			for _, s := range nextStates {
				an.StatesQueue.Push(&Item{
					value:    s,
					priority: an.PathSelector.CalculatePriority(s),
				})
			}

			// Для прервываемой стратегии выполнения
			/*if len(nextStates) != 0 {
				// Если количество стейтов отлично от нуля, должны обратиться к пас селектору
				// Например: вызов функции, ветвление или прыжoк
				break
			}*/

		}

		// Для прервываемой стратегии выполнения
		//if interpreter.completed() {
		an.Results = append(an.Results, interpreter)
		//}
	}

	return an.Results
}
