// Package translator содержит реализацию транслятора в Z3
package translator

import (
	"fmt"
	"math/big"
	"symbolic-execution-course/internal/symbolic"
	"symbolic-execution-course/internal/util"

	"github.com/ebukreev/go-z3/z3"
)

// Z3Translator транслирует символьные выражения в Z3 формулы
type Z3Translator struct {
	ctx    *z3.Context
	config *z3.Config
	vars   map[string]z3.Value // Кэш переменных
	funcs  map[string]z3.FuncDecl
}

// NewZ3Translator создаёт новый экземпляр Z3 транслятора
func NewZ3Translator() *Z3Translator {
	config := &z3.Config{}
	ctx := z3.NewContext(config)

	return &Z3Translator{
		ctx:    ctx,
		config: config,
		vars:   make(map[string]z3.Value),
		funcs:  make(map[string]z3.FuncDecl),
	}
}

// GetContext возвращает Z3 контекст
func (zt *Z3Translator) GetContext() interface{} {
	return zt.ctx
}

// Reset сбрасывает состояние транслятора
func (zt *Z3Translator) Reset() {
	zt.vars = make(map[string]z3.Value)
}

// Close освобождает ресурсы
func (zt *Z3Translator) Close() {
	// Z3 контекст закрывается автоматически
}

// TranslateExpression транслирует символьное выражение в Z3
func (zt *Z3Translator) TranslateExpression(expr symbolic.SymbolicExpression) (interface{}, error) {
	return expr.Accept(zt), nil
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// VisitVariable транслирует символьную переменную в Z3
func (zt *Z3Translator) VisitVariable(expr *symbolic.SymbolicVariable) interface{} {

	// Подсказки:
	// - Используйте zt.ctx.IntConst(name) для int переменных
	// - Используйте zt.ctx.BoolConst(name) для bool переменных
	// - Храните переменные в zt.vars для повторного использования

	if v, hasCache := zt.vars[expr.Name]; hasCache {
		return v
	}

	return zt.createZ3Variable(expr.Name, expr.Type(), expr.TypeGeneric)
}

// VisitIntConstant транслирует целочисленную константу в Z3
func (zt *Z3Translator) VisitIntConstant(expr *symbolic.IntConstant) interface{} {
	// Создать Z3 константу с помощью zt.ctx.FromBigInt или аналогичного метода

	return zt.ctx.FromBigInt(big.NewInt(expr.Value), zt.ctx.IntSort())
}

// VisitFloatConstant транслирует константу с плавающей точкой в Z3
func (zt *Z3Translator) VisitFloatConstant(expr *symbolic.FloatConstant) interface{} {
	// Создать Z3 константу с помощью zt.ctx.FromBigInt или аналогичного метода

	return zt.ctx.FromFloat64(expr.Value, zt.ctx.FloatSort(11, 53))
}

// VisitBoolConstant транслирует булеву константу в Z3
func (zt *Z3Translator) VisitBoolConstant(expr *symbolic.BoolConstant) interface{} {
	// Использовать zt.ctx.FromBool для создания Z3 булевой константы

	return zt.ctx.FromBool(expr.Value)
}

// VisitBinaryOperation транслирует бинарную операцию в Z3
func (zt *Z3Translator) VisitBinaryOperation(expr *symbolic.BinaryOperation) interface{} {
	// TODO: Реализовать
	// 1. Транслировать левый и правый операнды
	// 2. В зависимости от оператора создать соответствующую Z3 операцию

	// Подсказки по операциям в Z3:
	// - Арифметические: left.Add(right), left.Sub(right), left.Mul(right), left.Div(right)
	// - Сравнения: left.Eq(right), left.LT(right), left.LE(right), etc.
	// - Приводите типы: left.(z3.Int), right.(z3.Int) для int операций

	l := expr.Left.Accept(zt)
	r := expr.Right.Accept(zt)

	if expr.Left.Type() == symbolic.IntType {
		switch expr.Operator {
		case symbolic.ADD:
			return l.(z3.Int).Add(r.(z3.Int))
		case symbolic.SUB:
			return l.(z3.Int).Sub(r.(z3.Int))
		case symbolic.MUL:
			return l.(z3.Int).Mul(r.(z3.Int))
		case symbolic.DIV:
			return l.(z3.Int).Div(r.(z3.Int))
		case symbolic.MOD:
			return l.(z3.Int).Mod(r.(z3.Int))
		case symbolic.LT:
			return l.(z3.Int).LT(r.(z3.Int))
		case symbolic.LE:
			return l.(z3.Int).LE(r.(z3.Int))
		case symbolic.GT:
			return l.(z3.Int).GT(r.(z3.Int))
		case symbolic.GE:
			return l.(z3.Int).GE(r.(z3.Int))
		case symbolic.EQ:
			return l.(z3.Int).Eq(r.(z3.Int))
		case symbolic.NE:
			return l.(z3.Int).NE(r.(z3.Int))
		}
	}

	if expr.Left.Type() == symbolic.FloatType {
		switch expr.Operator {
		case symbolic.ADD:
			return l.(z3.Float).Add(r.(z3.Float))
		case symbolic.SUB:
			return l.(z3.Float).Sub(r.(z3.Float))
		case symbolic.MUL:
			return l.(z3.Float).Mul(r.(z3.Float))
		case symbolic.DIV:
			return l.(z3.Float).Div(r.(z3.Float))
		case symbolic.LT:
			return l.(z3.Float).LT(r.(z3.Float))
		case symbolic.LE:
			return l.(z3.Float).LE(r.(z3.Float))
		case symbolic.GT:
			return l.(z3.Float).GT(r.(z3.Float))
		case symbolic.GE:
			return l.(z3.Float).GE(r.(z3.Float))
		case symbolic.EQ:
			return l.(z3.Float).Eq(r.(z3.Float))
		case symbolic.NE:
			return l.(z3.Float).NE(r.(z3.Float))
		}
	}

	if expr.Operator == symbolic.EQ && expr.Left.Type() == symbolic.BoolType {
		return l.(z3.Bool).Eq(r.(z3.Bool))
	}

	if expr.Operator == symbolic.EQ && expr.Left.Type() == symbolic.ArrayType {
		return l.(z3.Array).Eq(r.(z3.Array))
	}

	if expr.Operator == symbolic.NE && expr.Left.Type() == symbolic.BoolType {
		return l.(z3.Bool).NE(r.(z3.Bool))
	}

	if expr.Operator == symbolic.NE && expr.Left.Type() == symbolic.ArrayType {
		return l.(z3.Array).NE(r.(z3.Array))
	}

	panic("unsupported binary operator")

}

// VisitLogicalOperation транслирует логическую операцию в Z3
func (zt *Z3Translator) VisitLogicalOperation(expr *symbolic.LogicalOperation) interface{} {

	// Подсказки:
	// - AND: zt.ctx.And(operands...)
	// - OR: zt.ctx.Or(operands...)
	// - NOT: operand.Not() (для единственного операнда)
	// - IMPLIES: antecedent.Implies(consequent)

	operands := util.Convert(expr.Operands, func(exp symbolic.SymbolicExpression) z3.Bool {
		return exp.Accept(zt).(z3.Bool)
	})

	switch expr.Operator {
	case symbolic.AND:
		return zt.ctx.FromBool(true).And(operands...)
	case symbolic.OR:
		return zt.ctx.FromBool(false).Or(operands...)
	case symbolic.NOT:
		return operands[0].Not()
	case symbolic.IMPLIES:
		return operands[0].Implies(operands[len(operands)-1])
	}

	panic("unsupported logical operator")
}

func (zt *Z3Translator) VisitUnaryOperation(expr *symbolic.UnaryOperation) interface{} {
	operand := expr.Operand.Accept(zt)
	switch expr.Operator {
	case symbolic.PLUS:
		return operand
	case symbolic.MINUS:
		if expr.Operand.Type() == symbolic.IntType {
			return operand.(z3.Int).Neg()
		} else {
			return operand.(z3.Float).Neg()
		}
	case symbolic.CARET:
		return operand.(z3.Int).ToBV(64).Not().SToInt()
	case symbolic.INCREMENT:
		return operand.(z3.Int).Add(zt.ctx.FromBigInt(big.NewInt(1), zt.ctx.IntSort()).(z3.Int))
	case symbolic.DECREMENT:
		return operand.(z3.Int).Sub(zt.ctx.FromBigInt(big.NewInt(1), zt.ctx.IntSort()).(z3.Int))
	}

	panic("unsupported unary operator")
}
func (zt *Z3Translator) VisitConditionalExpression(expr *symbolic.ConditionalExpression) interface{} {
	cond := expr.Condition.Accept(zt).(z3.Bool)
	then := expr.Then.Accept(zt).(z3.Value)
	elze := expr.Else.Accept(zt).(z3.Value)

	return cond.IfThenElse(then, elze)

}
func (zt *Z3Translator) VisitFunction(expr *symbolic.Function) interface{} {
	if v, hasCache := zt.funcs[expr.Name]; hasCache {
		return v
	}

	args := util.Convert(expr.Args, func(expr symbolic.SymbolicVariable) z3.Sort {
		return zt.sortForType(expr.Type(), expr.TypeGeneric)
	})

	ret := zt.sortForType(expr.ReturnType.ExprType, expr.ReturnType.Generic)

	zt.funcs[expr.Name] = zt.ctx.FuncDecl(expr.Name, args, ret)

	return zt.funcs[expr.Name]
}

func (zt *Z3Translator) VisitFunctionCall(expr *symbolic.FunctionCall) interface{} {
	fun := expr.Func.Accept(zt).(z3.FuncDecl)

	args := util.Convert(expr.Args, func(e symbolic.SymbolicExpression) z3.Value {
		return e.Accept(zt).(z3.Value)
	})

	return fun.Apply(args...)
}

func (zt *Z3Translator) VisitArraySelect(expr *symbolic.ArraySelect) interface{} {
	arr := expr.Array.Accept(zt).(z3.Array)
	i := expr.Index.Accept(zt).(z3.Int)

	return arr.Select(i)
}

func (zt *Z3Translator) VisitArrayStore(expr *symbolic.ArrayStore) interface{} {
	arr := expr.Array.Accept(zt).(z3.Array)
	i := expr.Index.Accept(zt).(z3.Int)
	v := expr.Value.Accept(zt).(z3.Value)

	return arr.Store(i, v)
}

func (zt *Z3Translator) VisitRef(ref *symbolic.Ref) interface{}{
	return ref.Deref.Accept(zt)
}
func (zt *Z3Translator) VisitFieldRead(f *symbolic.FieldRead) interface{} {
	return f.RawValue.Accept(zt)
}
func (zt *Z3Translator) VisitFieldWrite(f *symbolic.FieldWrite) interface{} {
	return f.RawValue.Accept(zt)
}

// Вспомогательные методы

// createZ3Variable создаёт Z3 переменную соответствующего типа
func (zt *Z3Translator) createZ3Variable(name string, exprType symbolic.ExpressionType, generic *symbolic.GenericType) z3.Value {
	// Создать Z3 переменную на основе типа
	zt.vars[name] = zt.ctx.Const(name, zt.sortForType(exprType, generic))
	return zt.vars[name]
}

func (zt *Z3Translator) sortForType(t symbolic.ExpressionType, generic *symbolic.GenericType) z3.Sort {
	switch t {
	case symbolic.IntType:
		return zt.ctx.IntSort()
	case symbolic.FloatType:
		return zt.ctx.FloatSort(11, 53)
	case symbolic.BoolType:
		return zt.ctx.BoolSort()
	case symbolic.ArrayType:
		indexSort := zt.ctx.IntSort()
		elementSort := zt.sortForType(generic.ExprType, generic.Generic)
		return zt.ctx.ArraySort(indexSort, elementSort)
	}

	panic("Unknown type")
}

// castToZ3Type приводит значение к нужному Z3 типу
func (zt *Z3Translator) castToZ3Type(value interface{}, targetType symbolic.ExpressionType) (z3.Value, error) {

	// Безопасно привести interface{} к конкретному Z3 типу
	switch targetType {
	case symbolic.IntType:
		return value.(z3.Int), nil
	case symbolic.FloatType:
		return value.(z3.Float), nil
	case symbolic.BoolType:
		return value.(z3.Bool), nil
	case symbolic.ArrayType:
		return value.(z3.Array), nil
	}

	return nil, fmt.Errorf("Unsupported cast")
}
