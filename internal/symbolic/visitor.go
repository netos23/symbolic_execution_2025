package symbolic

// Visitor интерфейс для обхода символьных выражений (Visitor Pattern)
type Visitor interface {
	VisitVariable(expr SymbolicVariable) interface{}
	VisitIntConstant(expr *IntConstant) interface{}
	VisitBoolConstant(expr *BoolConstant) interface{}
	VisitBinaryOperation(expr *BinaryOperation) interface{}
	VisitLogicalOperation(expr *LogicalOperation) interface{}
	VisitUnaryOperation(expr *UnaryOperation) interface{}
	VisitConditionalExpression(expr *ConditionalExpression) interface{}
	VisitFunction(expr *Function) interface{}
	VisitArraySelect(expr *ArraySelect) interface{}
	VisitFunctionCall(expr *FunctionCall) interface{}
	// TODO: Добавьте методы для других типов выражений по мере необходимости
}
