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
	CallStack     []CallStackFrame
	Analyser      *Analyser
	PathCondition symbolic.SymbolicExpression
	Heap          memory.Memory
}

type CallStackFrame struct {
	Function    *ssa.Function
	LocalMemory map[string]symbolic.SymbolicExpression
	ReturnValue []symbolic.SymbolicExpression
	Block       *ssa.BasicBlock
}

func (interpreter *Interpreter) copy() *Interpreter {
	return &Interpreter{
		CallStack:     interpreter.CallStack,
		Analyser:      interpreter.Analyser,
		PathCondition: interpreter.PathCondition,
		Heap:          interpreter.Heap,
	}
}

func (interpreter *Interpreter) copyWithBasicBlock(b *ssa.BasicBlock) *Interpreter {
	cp := interpreter.copy()
	i := len(cp.CallStack)
	f := util.Last(cp.CallStack)
	stack := make([]CallStackFrame, i+1)
	copy(stack, cp.CallStack)
	stack[i] = CallStackFrame{
		Function:    f.Function,
		LocalMemory: f.LocalMemory,
		ReturnValue: f.ReturnValue,
		Block:       b,
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
	default:
		return []Interpreter{}
	}
}

func (interpreter *Interpreter) resolveExpression(value ssa.Value) symbolic.SymbolicExpression {
	switch v := value.(type) {
	case *ssa.Const:
		return interpreter.resolveConst(v)
	default:
		frame := util.Last(interpreter.CallStack)
		if expr, ok := frame.LocalMemory[value.Name()]; ok {
			return expr
		}
	}

	panic("Unsupported value")
}

func (interpreter *Interpreter) interpretDynamicallyUnOp(e *ssa.UnOp) []Interpreter {
	v := interpreter.resolveExpression(e.X)

	var expr symbolic.SymbolicExpression
	switch e.Op {
	//case token.NOT:
	//	expr = symbolic.NewUnaryOperation(v,symbolic.MINUS)
	case token.SUB:
		expr = symbolic.NewUnaryOperation(v, symbolic.MINUS)
	/*case token.ARROW:
		expr = symbolic.NewUnaryOperation(v,symbolic.MINUS)
	case token.MUL:
		expr = symbolic.NewUnaryOperation(v,symbolic.MINUS)*/
	case token.XOR:
		expr = symbolic.NewUnaryOperation(v, symbolic.CARET)
	default:
	}

	localMemory := util.Last(interpreter.CallStack).LocalMemory
	v, hasLocal := localMemory[e.Name()]
	if hasLocal {
		localMemory[e.Name()] = interpreter.Heap.Assign(v, expr)
	} else {
		v := interpreter.Heap.Allocate(
			expr.Type(), symbolic.ObjectNameFor(expr), symbolic.GenericFor(expr),
		)
		localMemory[e.Name()] = interpreter.Heap.Assign(v, expr)
	}

	return []Interpreter{}
}

func (interpreter *Interpreter) interpretDynamicallyBinOp(e *ssa.BinOp) []Interpreter {
	expr := interpreter.resolveBinOp(e)

	localMemory := util.Last(interpreter.CallStack).LocalMemory
	v, hasLocal := localMemory[e.Name()]
	if hasLocal {
		localMemory[e.Name()] = interpreter.Heap.Assign(v, expr)
	} else {
		v := interpreter.Heap.Allocate(
			expr.Type(), symbolic.ObjectNameFor(expr), symbolic.GenericFor(expr),
		)
		localMemory[e.Name()] = interpreter.Heap.Assign(v, expr)
	}

	return []Interpreter{}
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
		return []Interpreter{*interpreter}
	}

	f.LocalMemory[e.Name()] = interpreter.Heap.Allocate(builder.ConvertToSymbolic(e.Type()))

	return []Interpreter{*interpreter}
}

func (interpreter *Interpreter) interpretDynamicallyJump(e *ssa.Jump) []Interpreter {
	return util.Convert(
		e.Block().Succs,
		func(b *ssa.BasicBlock) Interpreter {
			cp := interpreter.copyWithBasicBlock(b)
			return *cp
		},
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

	return []Interpreter{*trueInterp, *falseInterp}
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

	return []Interpreter{}
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
