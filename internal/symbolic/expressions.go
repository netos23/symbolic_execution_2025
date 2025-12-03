// Package symbolic содержит конкретные реализации символьных выражений
package symbolic

import (
	"fmt"
	"slices"
	"strings"
	"symbolic-execution-course/internal/util"
)

// SymbolicExpression - базовый интерфейс для всех символьных выражений
type SymbolicExpression interface {
	// Type возвращает тип выражения
	Type() ExpressionType

	// String возвращает строковое представление выражения
	String() string

	// Accept принимает visitor для обхода дерева выражений
	Accept(visitor Visitor) interface{}
}

// SymbolicVariable представляет символьную переменную
type SymbolicVariable struct {
	Name        string
	VarType     ExpressionType
	TypeGeneric *GenericType
	ObjType     *Object
}

// NewSymbolicVariable создаёт новую символьную переменную
func NewSymbolicVariable(name string, exprType ExpressionType, generic *GenericType, objType *Object) *SymbolicVariable {
	return &SymbolicVariable{
		Name:        name,
		VarType:     exprType,
		TypeGeneric: generic,
		ObjType:     objType,
	}
}

// Type возвращает тип переменной
func (sv *SymbolicVariable) Type() ExpressionType {
	return sv.VarType
}

// String возвращает строковое представление переменной
func (sv *SymbolicVariable) String() string {
	return sv.Name
}

// Accept реализует Visitor pattern
func (sv *SymbolicVariable) Accept(visitor Visitor) interface{} {
	return visitor.VisitVariable(sv)
}

// IntConstant представляет целочисленную константу
type IntConstant struct {
	Value int64
}

// NewIntConstant создаёт новую целочисленную константу
func NewIntConstant(value int64) *IntConstant {
	return &IntConstant{Value: value}
}

// Type возвращает тип константы
func (ic *IntConstant) Type() ExpressionType {
	return IntType
}

// String возвращает строковое представление константы
func (ic *IntConstant) String() string {
	return fmt.Sprintf("%d", ic.Value)
}

// Accept реализует Visitor pattern
func (ic *IntConstant) Accept(visitor Visitor) interface{} {
	return visitor.VisitIntConstant(ic)
}

// FloatConstant представляет целочисленную константу
type FloatConstant struct {
	Value float64
}

// NewFloatConstant создаёт новую целочисленную константу
func NewFloatConstant(value float64) *FloatConstant {
	return &FloatConstant{Value: value}
}

// Type возвращает тип константы
func (ic *FloatConstant) Type() ExpressionType {
	return FloatType
}

// String возвращает строковое представление константы
func (ic *FloatConstant) String() string {
	return fmt.Sprintf("%f", ic.Value)
}

// Accept реализует Visitor pattern
func (ic *FloatConstant) Accept(visitor Visitor) interface{} {
	return visitor.VisitFloatConstant(ic)
}

// BoolConstant представляет булеву константу
type BoolConstant struct {
	Value bool
}

// NewBoolConstant создаёт новую булеву константу
func NewBoolConstant(value bool) *BoolConstant {
	return &BoolConstant{Value: value}
}

// Type возвращает тип константы
func (bc *BoolConstant) Type() ExpressionType {
	return BoolType
}

// String возвращает строковое представление константы
func (bc *BoolConstant) String() string {
	return fmt.Sprintf("%t", bc.Value)
}

// Accept реализует Visitor pattern
func (bc *BoolConstant) Accept(visitor Visitor) interface{} {
	return visitor.VisitBoolConstant(bc)
}

// BinaryOperation представляет бинарную операцию
type BinaryOperation struct {
	Left     SymbolicExpression
	Right    SymbolicExpression
	Operator BinaryOperator
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// NewBinaryOperation создаёт новую бинарную операцию
func NewBinaryOperation(left, right SymbolicExpression, op BinaryOperator) *BinaryOperation {
	if left.Type() != right.Type() {
		return nil
	}

	switch op {
	case ADD:
		if left.Type() == BoolType {
			return nil
		}

		if left.Type() == ArrayType {
			return nil
		}
	case SUB:
		if left.Type() == BoolType {
			return nil
		}

		if left.Type() == ArrayType {
			return nil
		}
	case MUL:
		if left.Type() == BoolType {
			return nil
		}

		if left.Type() == ArrayType {
			return nil
		}
	case DIV:
		if left.Type() == BoolType {
			return nil
		}

		if left.Type() == ArrayType {
			return nil
		}
	case MOD:
		if left.Type() == FloatType {
			return nil
		}

		if left.Type() == BoolType {
			return nil
		}

		if left.Type() == ArrayType {
			return nil
		}
	}

	return &BinaryOperation{
		Operator: op,
		Left:     left,
		Right:    right,
	}
}

// Type возвращает результирующий тип операции
func (bo *BinaryOperation) Type() ExpressionType {
	switch bo.Operator {
	case ADD, SUB, MUL, DIV, MOD:
		return bo.Left.Type()
	case EQ, NE, LT, LE, GT, GE:
		return BoolType
	default:
		panic("unknown type of binary operation")
	}
}

// String возвращает строковое представление операции
func (bo *BinaryOperation) String() string {
	return fmt.Sprintf("(%s %s %s)", bo.Left.String(), bo.Operator.String(), bo.Right.String())
}

// Accept реализует Visitor pattern
func (bo *BinaryOperation) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryOperation(bo)
}

// LogicalOperation представляет логическую операцию
type LogicalOperation struct {
	Operands []SymbolicExpression
	Operator LogicalOperator
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// NewLogicalOperation создаёт новую логическую операцию
func NewLogicalOperation(operands []SymbolicExpression, op LogicalOperator) *LogicalOperation {
	if slices.ContainsFunc(operands, func(e SymbolicExpression) bool {
		return e.Type() != BoolType
	}) {
		return nil
	}

	switch op {
	case NOT:
		if len(operands) != 1 {
			return nil
		}

		return &LogicalOperation{
			Operands: operands,
			Operator: op,
		}
	case OR, AND:
		return &LogicalOperation{
			Operands: operands,
			Operator: op,
		}
	case IMPLIES:
		if len(operands) != 2 {
			return nil
		}

		return &LogicalOperation{
			Operands: operands,
			Operator: op,
		}
	}

	return nil
}

// Type возвращает тип логической операции (всегда bool)
func (lo *LogicalOperation) Type() ExpressionType {
	return BoolType
}

// String возвращает строковое представление логической операции
func (lo *LogicalOperation) String() string {
	// TODO: Реализовать
	// Для NOT: "!operand"
	// Для AND/OR: "(operand1 && operand2 && ...)"
	// Для IMPLIES: "(operand1 => operand2)"
	switch lo.Operator {
	case NOT:
		return fmt.Sprintf("!%s", lo.Operands[len(lo.Operands)-1].String())
	case AND, OR:
		return fmt.Sprintf(
			"(%s)",
			strings.Join(
				util.Convert(
					lo.Operands, func(e SymbolicExpression) string {
						return e.String()
					},
				),
				lo.Operator.String(),
			),
		)
	case IMPLIES:
		return fmt.Sprintf("(%s => %s)", lo.Operands[0].String(), lo.Operands[len(lo.Operands)-1].String())
	default:
		return ""
	}
}

// Accept реализует Visitor pattern
func (lo *LogicalOperation) Accept(visitor Visitor) interface{} {
	return visitor.VisitLogicalOperation(lo)
}

// Операторы для бинарных выражений
type BinaryOperator int

const (
	// Арифметические операторы
	ADD BinaryOperator = iota
	SUB
	MUL
	DIV
	MOD

	// Операторы сравнения
	EQ // равно
	NE // не равно
	LT // меньше
	LE // меньше или равно
	GT // больше
	GE // больше или равно
)

// String возвращает строковое представление оператора
func (op BinaryOperator) String() string {
	switch op {
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"
	case MOD:
		return "%"
	case EQ:
		return "=="
	case NE:
		return "!="
	case LT:
		return "<"
	case LE:
		return "<="
	case GT:
		return ">"
	case GE:
		return ">="
	default:
		return "unknown"
	}
}

// Логические операторы
type LogicalOperator int

const (
	AND LogicalOperator = iota
	OR
	NOT
	IMPLIES
)

// String возвращает строковое представление логического оператора
func (op LogicalOperator) String() string {
	switch op {
	case AND:
		return "&&"
	case OR:
		return "||"
	case NOT:
		return "!"
	case IMPLIES:
		return "=>"
	default:
		return "unknown"
	}
}

type UnaryOperator int

const (
	PLUS UnaryOperator = iota
	MINUS
	CARET
	INCREMENT
	DECREMENT
)

func (op UnaryOperator) String() string {
	switch op {
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case CARET:
		return "^"
	case INCREMENT:
		return "++"
	case DECREMENT:
		return "--"
	default:
		return "unknown"
	}
}

type Ref struct {
	Address     int64
	VarType     ExpressionType
	TypeGeneric *GenericType
	ObjType     *Object
	Deref       SymbolicExpression
}

func NewRef(
	address int64,
	tpe ExpressionType,
	generic *GenericType,
	obj *Object,
	deref SymbolicExpression,
) *Ref {
	return &Ref{
		Address:     address,
		VarType:     tpe,
		TypeGeneric: generic,
		ObjType:     obj,
		Deref:       deref,
	}
}

func (ref *Ref) Type() ExpressionType {
	return RefType
}

func (ref *Ref) String() string {
	return fmt.Sprintf("0x%d", ref.Address)
}

func (ref *Ref) Accept(visitor Visitor) interface{} {
	return visitor.VisitRef(ref)
}

// TODO: Добавьте дополнительные типы выражений по необходимости:
// - UnaryOperation (унарные операции: -x, !x)

type UnaryOperation struct {
	Operand  SymbolicExpression
	Operator UnaryOperator
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// NewUnaryOperation создаёт новую бинарную операцию
func NewUnaryOperation(operand SymbolicExpression, op UnaryOperator) *UnaryOperation {
	switch op {
	case PLUS:
		if operand.Type() != IntType && operand.Type() != FloatType {
			return nil
		}
	case MINUS:
		if operand.Type() != IntType && operand.Type() != FloatType {
			return nil
		}
	case CARET:
		if operand.Type() != IntType {
			return nil
		}
	case INCREMENT:
		if operand.Type() != IntType {
			return nil
		}
	case DECREMENT:
		if operand.Type() != IntType {
			return nil
		}
	}

	return &UnaryOperation{
		Operand:  operand,
		Operator: op,
	}
}

// Type возвращает результирующий тип операции
func (uo *UnaryOperation) Type() ExpressionType {
	return uo.Operand.Type()
}

// String возвращает строковое представление операции
func (uo *UnaryOperation) String() string {
	switch uo.Operator {
	case PLUS, MINUS, CARET:
		return uo.Operator.String() + uo.Operand.String()
	default:
		return uo.Operand.String() + uo.Operator.String()
	}
}

// Accept реализует Visitor pattern
func (uo *UnaryOperation) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryOperation(uo)
}

// - ArrayAccess (доступ к элементам массива: arr[index])

type ArraySelect struct {
	Array SymbolicExpression // Сам массив
	Index SymbolicExpression // Тип элемента массива
}

// NewArraySelect создаёт выражение arr[idx]
func NewArraySelect(arr SymbolicExpression, idx SymbolicExpression) *ArraySelect {
	return &ArraySelect{Array: arr, Index: idx}
}

// Type возвращает тип элемента массива
func (as *ArraySelect) Type() ExpressionType {
	arrayType := as.Array
	return GenericFor(arrayType).ExprType
}

func ObjectFor(obj SymbolicExpression) *Object {
	switch o := obj.(type) {
	case *SymbolicVariable:
		return o.ObjType
	case *Ref:
		return o.ObjType
	}

	panic("Wrong object")
}

func GenericFor(arrayType SymbolicExpression) *GenericType {
	if arrayType.Type() != ArrayType {
		return nil
	}

	for arrayType.Type() == ArrayType {
		if inner, ok := arrayType.(*ArrayStore); ok {
			arrayType = inner.Array
		}

		if val, ok := arrayType.(*SymbolicVariable); ok {
			return val.TypeGeneric
		}

		if val, ok := arrayType.(*FieldRead); ok {
			switch f := val.Obj.(type) {
			case *Ref:
				return f.ObjType.Fields[val.Index].Generic
			case *SymbolicVariable:
				return f.ObjType.Fields[val.Index].Generic
			}
		}
	}

	panic("Unknown reciver")
}

// String:  arr[idx]
func (as *ArraySelect) String() string {
	return fmt.Sprintf("%s[%s]", as.Array.String(), as.Index.String())
}

// Accept реализует Visitor pattern
func (as *ArraySelect) Accept(visitor Visitor) interface{} { return visitor.VisitArraySelect(as) }

// - ArrayAccess (доступ к элементам массива: arr[index])

type ArrayStore struct {
	Array SymbolicExpression
	Index SymbolicExpression
	Value SymbolicExpression
}

// NewArrayStore создаёт выражение arr[idx]
func NewArrayStore(arr SymbolicExpression, idx SymbolicExpression, v SymbolicExpression) *ArrayStore {

	return &ArrayStore{Array: arr, Index: idx, Value: v}
}

// Type возвращает тип элемента массива
func (as *ArrayStore) Type() ExpressionType {
	return ArrayType
}

// String:  arr[idx]
func (as *ArrayStore) String() string {
	return fmt.Sprintf("(%s[%s] = %s)", as.Array.String(), as.Index.String(), as.Value.String())
}

// Accept реализует Visitor pattern
func (as *ArrayStore) Accept(visitor Visitor) interface{} { return visitor.VisitArrayStore(as) }

// - FunctionCall (вызовы функций: f(x, y))

type Function struct {
	Name       string
	Args       []SymbolicVariable
	ReturnType GenericType
}

// NewFunction создаёт выражение arr[idx]
func NewFunction(name string, args []SymbolicVariable, returnType GenericType) *Function {
	return &Function{
		Name:       name,
		Args:       args,
		ReturnType: returnType,
	}
}

// Type возвращает результирующий тип операции
func (f *Function) Type() ExpressionType {
	return f.ReturnType.ExprType
}

// String возвращает строковое представление операции
func (f *Function) String() string {
	return fmt.Sprintf(
		"%s %s(%s)",
		f.Type(),
		f.Name,
		strings.Join(
			util.Convert(
				f.Args, func(e SymbolicVariable) string {
					return e.Name
				},
			),
			",",
		),
	)
}

// Accept реализует Visitor pattern
func (f *Function) Accept(visitor Visitor) interface{} {
	return visitor.VisitFunction(f)
}

type FunctionCall struct {
	Func Function
	Args []SymbolicExpression
}

// NewArraySelect создаёт выражение arr[idx]
func NewFunctionCall(fun Function, args []SymbolicExpression) *FunctionCall {
	for i, e := range args {
		if e.Type() != fun.Args[i].Type() {
			return nil
		}
	}

	return &FunctionCall{
		Func: fun,
		Args: args,
	}
}

// Type возвращает результирующий тип операции
func (fc *FunctionCall) Type() ExpressionType {
	return fc.Func.Type()
}

// String возвращает строковое представление операции
func (fc *FunctionCall) String() string {
	return fmt.Sprintf(
		"%s(%s)",
		fc.Func.Name,
		strings.Join(
			util.Convert(
				fc.Args, func(e SymbolicExpression) string {
					return e.String()
				},
			),
			",",
		),
	)
}

// Accept реализует Visitor pattern
func (fc *FunctionCall) Accept(visitor Visitor) interface{} {
	return visitor.VisitFunctionCall(fc)
}

// - ConditionalExpression (тернарный оператор: condition ? true_expr : false_expr)

type ConditionalExpression struct {
	Condition LogicalOperation
	Then      SymbolicExpression
	Else      SymbolicExpression
}

func NewConditionalExpression(
	condition LogicalOperation,
	then SymbolicExpression,
	elze SymbolicExpression,
) *ConditionalExpression {
	if then.Type() != elze.Type() {
		return nil
	}

	return &ConditionalExpression{
		Condition: condition,
		Then:      then,
		Else:      elze,
	}
}

// Type возвращает результирующий тип операции
func (ce *ConditionalExpression) Type() ExpressionType {
	return ce.Then.Type()
}

// String возвращает строковое представление операции
func (ce *ConditionalExpression) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", ce.Condition.String(), ce.Then.String(), ce.Else.String())
}

// Accept реализует Visitor pattern
func (ce *ConditionalExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitConditionalExpression(ce)
}

type FieldRead struct {
	Obj   SymbolicExpression
	Index int
	RawValue SymbolicExpression
}

func NewFieldRead(obj SymbolicExpression, idx int, value SymbolicExpression) *FieldRead {
	return &FieldRead{
		Obj:   obj,
		Index: idx,
		RawValue: value,
	}
}

func (f *FieldRead) Type() ExpressionType {
	switch v := f.Obj.(type) {
	case *SymbolicVariable:
		return v.ObjType.Fields[f.Index].ExprType
	case *Ref:
		return v.ObjType.Fields[f.Index].ExprType
	}

	panic("Wrong reciver")
}

func (f *FieldRead) String() string {
	return fmt.Sprintf("%s.%d", f.Obj, f.Index)
}

func (f *FieldRead) Accept(v Visitor) interface{} {
	return v.VisitFieldRead(f)
}

type FieldWrite struct {
	Obj   SymbolicExpression
	Index int
	Value SymbolicExpression
	RawValue SymbolicExpression
}

func NewFieldWrite(obj SymbolicExpression, index int, value SymbolicExpression, raw SymbolicExpression) *FieldWrite {
	return &FieldWrite{
		Obj:   obj,
		Index: index,
		Value: value,
		RawValue: value,
	}
}

func (f *FieldWrite) Type() ExpressionType {
	return f.Value.Type()
}

func (f *FieldWrite) String() string {
	return fmt.Sprintf("%s.%d = %s", f.Obj, f.Index, f.Value)
}

func (f *FieldWrite) Accept(visitor Visitor) interface{} {
	return visitor.VisitFieldWrite(f)
}
