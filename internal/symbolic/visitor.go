package symbolic

// Visitor интерфейс для обхода символьных выражений (Visitor Pattern)
type Visitor interface {
	VisitVariable(expr *SymbolicVariable) interface{}
	VisitIntConstant(expr *IntConstant) interface{}
	VisitFloatConstant(expr *FloatConstant) interface{}
	VisitBoolConstant(expr *BoolConstant) interface{}
	VisitBinaryOperation(expr *BinaryOperation) interface{}
	VisitLogicalOperation(expr *LogicalOperation) interface{}
	VisitUnaryOperation(expr *UnaryOperation) interface{}
	VisitConditionalExpression(expr *ConditionalExpression) interface{}
	VisitFunction(expr *Function) interface{}
	VisitArraySelect(expr *ArraySelect) interface{}
	VisitArrayStore(expr *ArrayStore) interface{}
	VisitFunctionCall(expr *FunctionCall) interface{}
	VisitRef(ref *Ref) interface{}
	VisitFieldRead(f *FieldRead) interface{}
	VisitFieldWrite(f *FieldWrite) interface{}
	// TODO: Добавьте методы для других типов выражений по мере необходимости
}
