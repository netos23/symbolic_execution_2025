package internal

import (
	"go/constant"
	"go/token"
	"golang.org/x/tools/go/ssa"
	"symbolic-execution-course/internal/memory"
	builder "symbolic-execution-course/internal/ssa"
	"symbolic-execution-course/internal/symbolic"
	"symbolic-execution-course/internal/util"
)

type Interpreter struct {
	CallStack      []CallStackFrame
	Analyser       *Analyser
	PathCondition  symbolic.SymbolicExpression
	Heap           memory.Memory
	RecursionLimit int
}

type CallStackFrame struct {
	Function           *ssa.Function
	LocalMemory        map[string]symbolic.SymbolicExpression
	ReturnValue        []symbolic.SymbolicExpression
	Block              *ssa.BasicBlock
	InstructionPointer int
	ReturnAddress      symbolic.SymbolicExpression
}

func (interpreter *Interpreter) copy() *Interpreter {
	return &Interpreter{
		CallStack:      interpreter.CallStack,
		Analyser:       interpreter.Analyser,
		PathCondition:  interpreter.PathCondition,
		Heap:           interpreter.Heap,
		RecursionLimit: interpreter.RecursionLimit,
	}
}

func (interpreter *Interpreter) copyWithBasicBlock(b *ssa.BasicBlock) *Interpreter {
	cp := interpreter.copy()
	i := len(cp.CallStack)
	f := util.Last(cp.CallStack)
	stack := make([]CallStackFrame, i+1)
	copy(stack, cp.CallStack)
	stack[i] = CallStackFrame{
		Function:           f.Function,
		LocalMemory:        f.LocalMemory,
		ReturnValue:        f.ReturnValue,
		Block:              b,
		InstructionPointer: 0,
	}
	cp.CallStack = stack
	return cp
}

func (interpreter *Interpreter) interpretDynamically(element ssa.Instruction) []Interpreter {
	switch e := element.(type) {
	case *ssa.Return:
		return interpreter.interpretDynamicallyReturn(e)
	case *ssa.If:
		return interpreter.interpretDynamicallyIf(e)
	case *ssa.Jump:
		return interpreter.interpretDynamicallyJump(e)
	case *ssa.UnOp:
		return interpreter.interpretDynamicallyUnOp(e)
	case *ssa.BinOp:
		return interpreter.interpretDynamicallyBinOp(e)
	case *ssa.Alloc:
		return interpreter.interpretDynamicallyAlloc(e)
	case *ssa.Call:
		return interpreter.interpretDynamicallyCall(e)
	case *ssa.Store:
		return interpreter.interpretDynamicallyStore(e)
	case *ssa.Phi:
		return interpreter.interpretDynamicallyPhi(e)
	case *ssa.ChangeType:
		return interpreter.interpretDynamicallyChangeType(e)
	case *ssa.Convert:
		return interpreter.interpretDynamicallyConvert(e)
	case *ssa.MakeInterface:
		return interpreter.interpretDynamicallyMakeInterface(e)
	case *ssa.FieldAddr:
		return interpreter.interpretDynamicallyFieldAddr(e)
	case *ssa.Field:
		return interpreter.interpretDynamicallyField(e)
	case *ssa.IndexAddr:
		return interpreter.interpretDynamicallyIndexAddr(e)
	case *ssa.Index:
		return interpreter.interpretDynamicallyIndex(e)
	default:
		return []Interpreter{}
	}
}

func (interpreter *Interpreter) resolveExpression(value ssa.Value) symbolic.SymbolicExpression {
	switch v := value.(type) {
	case *ssa.Const:
		return interpreter.resolveConst(v)
	case *ssa.UnOp:
		return interpreter.resolveUnOp(v)
	case *ssa.BinOp:
		return interpreter.resolveBinOp(v)
	default:
		frame := util.Last(interpreter.CallStack)
		if expr, ok := frame.LocalMemory[value.Name()]; ok {
			return expr
		}
	}

	panic("Unsupported value")
}

func (interpreter *Interpreter) interpretDynamicallyUnOp(e *ssa.UnOp) []Interpreter {
	expr := interpreter.resolveUnOp(e)

	interpreter.assign(e.Name(), expr)

	return interpreter.nextInstruction()
}

func (interpreter *Interpreter) resolveUnOp(e *ssa.UnOp) symbolic.SymbolicExpression {
	v := interpreter.resolveExpression(e.X)

	var expr symbolic.SymbolicExpression
	switch e.Op {
	//case token.NOT:
	//	expr = symbolic.NewUnaryOperation(v,symbolic.MINUS)
	case token.SUB:
		expr = symbolic.NewUnaryOperation(v, symbolic.MINUS)
	/*case token.ARROW:
	expr = symbolic.NewUnaryOperation(v,symbolic.MINUS)

	*/
	// TODO implement load
	/*
		case token.MUL:
			expr = symbolic.NewUnaryOperation(v,symbolic.MINUS)*/
	case token.XOR:
		expr = symbolic.NewUnaryOperation(v, symbolic.CARET)
	default:
	}
	return expr
}

func (interpreter *Interpreter) interpretDynamicallyBinOp(e *ssa.BinOp) []Interpreter {
	expr := interpreter.resolveBinOp(e)

	interpreter.assign(e.Name(), expr)

	return interpreter.nextInstruction()
}

func (interpreter *Interpreter) assign(name string, expr symbolic.SymbolicExpression) {
	localMemory := util.Last(interpreter.CallStack).LocalMemory
	v, hasLocal := localMemory[name]
	if hasLocal {
		localMemory[name] = interpreter.Heap.Assign(v, expr)
	} else {
		v := interpreter.Heap.Allocate(
			expr.Type(), symbolic.ObjectNameFor(expr), symbolic.GenericFor(expr),
		)
		localMemory[name] = interpreter.Heap.Assign(v, expr)
	}
}

func (interpreter *Interpreter) resolveBinOp(e *ssa.BinOp) symbolic.SymbolicExpression {
	lhs := interpreter.resolveExpression(e.X)
	rhs := interpreter.resolveExpression(e.Y)

	var expr symbolic.SymbolicExpression
	switch e.Op {
	case token.ADD:
		expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.ADD)
	case token.SUB:
		expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.SUB)
	case token.MUL:
		expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.MUL)
	case token.QUO:
		expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.DIV)
	case token.REM:
		expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.MOD)
	// TODO: Add bit operator support
	/*	case token.AND:
			expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.AND)
		case token.OR:
			expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.OR)
		case token.XOR:
			expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.XOR)
		case token.SHL:
			expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.SHL)
		case token.SHR:
			expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.SHR)
		case token.AND_NOT:
			expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.AND_NOT)*/
	case token.EQL:
		expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.EQ)
	case token.LSS:
		expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.LT)
	case token.GTR:
		expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.GT)
	case token.NEQ:
		expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.NE)
	case token.LEQ:
		expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.LE)
	case token.GEQ:
		expr = symbolic.NewBinaryOperation(lhs, rhs, symbolic.GE)
	default:
		panic("Unsupported binop")
	}
	return expr
}

func (interpreter *Interpreter) interpretDynamicallyAlloc(e *ssa.Alloc) []Interpreter {
	f := util.Last(interpreter.CallStack)

	if _, ok := f.LocalMemory[e.Name()]; ok {
		return interpreter.nextInstruction()
	}

	f.LocalMemory[e.Name()] = interpreter.Heap.Allocate(builder.ConvertToSymbolic(e.Type()))

	return interpreter.nextInstruction()
}

func (interpreter *Interpreter) interpretDynamicallyJump(e *ssa.Jump) []Interpreter {
	return interpreter.nextStates(
		util.Convert(
			e.Block().Succs,
			func(b *ssa.BasicBlock) Interpreter {
				cp := interpreter.copyWithBasicBlock(b)
				return *cp
			},
		),
	)
}

// internal/dynamic_interpreter.go
func (interpreter *Interpreter) interpretDynamicallyIf(e *ssa.If) []Interpreter {
	condExpr := interpreter.resolveExpression(e.Cond)

	succs := e.Block().Succs
	if len(succs) != 2 {
		panic("If instruction must have exactly two successors")
	}

	// True branch
	trueBlock := succs[0]
	trueInterp := interpreter.copyWithBasicBlock(trueBlock)
	trueInterp.addPathCondition(condExpr)

	// False branch
	falseBlock := succs[1]
	falseInterp := interpreter.copyWithBasicBlock(falseBlock)
	negCond := symbolic.NewLogicalOperation([]symbolic.SymbolicExpression{condExpr}, symbolic.NOT)
	falseInterp.addPathCondition(negCond)

	return interpreter.nextStates([]Interpreter{*trueInterp, *falseInterp})
}

func (interpreter *Interpreter) addPathCondition(condExpr symbolic.SymbolicExpression) {
	if interpreter.PathCondition != nil {
		interpreter.PathCondition = symbolic.NewLogicalOperation(
			[]symbolic.SymbolicExpression{interpreter.PathCondition, condExpr}, symbolic.AND,
		)
	} else {
		interpreter.PathCondition = condExpr
	}
}

func (interpreter *Interpreter) interpretDynamicallyReturn(e *ssa.Return) []Interpreter {
	f := util.Last(interpreter.CallStack)

	f.ReturnValue = util.Convert(e.Results, func(e ssa.Value) symbolic.SymbolicExpression {
		return interpreter.resolveExpression(e)
	})

	return interpreter.nextInstruction()
}

func (interpreter *Interpreter) resolveConst(v *ssa.Const) symbolic.SymbolicExpression {
	switch v.Value.Kind() {
	case constant.Bool:
		return symbolic.NewBoolConstant(constant.BoolVal(v.Value))
	case constant.Int:
		return symbolic.NewIntConstant(v.Int64())
	case constant.Float:
		return symbolic.NewFloatConstant(v.Float64())
	default:
		panic("unsupported const")
	}
}

func (interpreter *Interpreter) interpretDynamicallyCall(e *ssa.Call) []Interpreter {

	var callee *ssa.Function
	if e.Call.StaticCallee() != nil {
		callee = e.Call.StaticCallee()
	} else {
		// Indirect calls, are not supported here
		return interpreter.nextInstruction()
	}

	// Count occurrences of callee in the call stack
	count := 0
	for _, frame := range interpreter.CallStack {
		if frame.Function == callee {
			count++
		}
	}
	if count >= interpreter.RecursionLimit {
		// Recursion limit reached, skip the call
		return interpreter.nextInstruction()
	}

	ep := util.FirstOrNil(callee.Blocks)
	functionCall := interpreter.copyWithBasicBlock(ep)
	functionFrame := util.Last(functionCall.CallStack)
	functionFrame.Function = callee

	callArgs := e.Call.Args
	pathCondition := make([]symbolic.SymbolicExpression, len(callee.Params))
	locals := make(map[string]symbolic.SymbolicExpression)
	for i, p := range callee.Params {
		locals[p.Name()] = functionCall.Heap.MakeRef(builder.ConvertToSymbolic(p.Type()))
		pathCondition[i] = symbolic.NewBinaryOperation(
			locals[p.Name()], interpreter.resolveExpression(callArgs[i]), symbolic.EQ,
		)
	}
	functionFrame.LocalMemory = locals
	functionCall.addPathCondition(symbolic.NewLogicalOperation(pathCondition, symbolic.AND))

	results := callee.Signature.Results()
	if results != nil {
		returnAddress := interpreter.Heap.Allocate(builder.ConvertToSymbolic(results.At(0).Type()))
		interpreter.assign(e.Name(), returnAddress)
		functionFrame.ReturnAddress = returnAddress
	}

	return interpreter.nextStates([]Interpreter{*functionCall})
}

func (interpreter *Interpreter) completed() bool {
	f := util.Last(interpreter.CallStack)

	return f.Block != nil && f.InstructionPointer >= len(f.Block.Instrs)
}

func (interpreter *Interpreter) nextInstruction() []Interpreter {
	return interpreter.nextStates([]Interpreter{})
}

func (interpreter *Interpreter) nextStates(interpreters []Interpreter) []Interpreter {
	f := util.Last(interpreter.CallStack)
	f.InstructionPointer++
	return interpreters
}

func (interpreter *Interpreter) interpretDynamicallyStore(e *ssa.Store) []Interpreter {
	// todo implement
	panic("implement me")
}

func (interpreter *Interpreter) interpretDynamicallyPhi(e *ssa.Phi) []Interpreter {
	frame := util.Last(interpreter.CallStack)
	currBlock := frame.Block

	for i := len(interpreter.CallStack) - 1; i >= 0; i-- {
		stackFrame := interpreter.CallStack[i]
		for predIdx, pred := range currBlock.Preds {
			if stackFrame.Block == pred {
				if predIdx < len(e.Edges) {
					val := e.Edges[predIdx]
					if v, ok := stackFrame.LocalMemory[val.Name()]; ok {
						frame.LocalMemory[e.Name()] = v
						return interpreter.nextInstruction()
					}
				}
			}
		}
	}

	panic("No matching predecessor found in call stack for Phi node")
}

func (interpreter *Interpreter) interpretDynamicallyChangeType(e *ssa.ChangeType) []Interpreter {
	return interpreter.nextInstruction()
}

func (interpreter *Interpreter) interpretDynamicallyConvert(e *ssa.Convert) []Interpreter {
	exp := interpreter.resolveExpression(e.X)

	frame := util.Last(interpreter.CallStack)
	frame.LocalMemory[e.Name()] = exp

	return interpreter.nextInstruction()
}

func (interpreter *Interpreter) interpretDynamicallyMakeInterface(e *ssa.MakeInterface) []Interpreter {
	frame := util.Last(interpreter.CallStack)

	/*if _, ok := frame.LocalMemory[e.Name()]; ok {
		return interpreter.nextInstruction()
	}*/

	frame.LocalMemory[e.Name()] = interpreter.Heap.Allocate(
		builder.ConvertToSymbolic(e.Type()),
	)

	return interpreter.nextInstruction()
}

func (interpreter *Interpreter) interpretDynamicallyFieldAddr(e *ssa.FieldAddr) []Interpreter {
	// todo implement
	panic("implement me")
}

func (interpreter *Interpreter) interpretDynamicallyField(e *ssa.Field) []Interpreter {
	frame := util.Last(interpreter.CallStack)

	base := interpreter.resolveExpression(e.X)
	fieldIndex := e.Field
	fieldValue := interpreter.Heap.GetFieldValue(base.(*symbolic.Ref), fieldIndex)
	frame.LocalMemory[e.Name()] = fieldValue

	return interpreter.nextInstruction()
}

func (interpreter *Interpreter) interpretDynamicallyIndexAddr(e *ssa.IndexAddr) []Interpreter {
	// todo implement
	panic("implement me")
}

func (interpreter *Interpreter) interpretDynamicallyIndex(e *ssa.Index) []Interpreter {
	// todo implement
	panic("implement me")
}
