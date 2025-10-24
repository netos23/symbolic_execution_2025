// Package symbolic определяет базовые типы символьных выражений
package symbolic

import "fmt"

// ExpressionType представляет тип символьного выражения
type ExpressionType int

const (
	IntType ExpressionType = iota
	FloatType
	BoolType
	ArrayType
	// Добавьте другие типы по необходимости
)

type GenericType struct {
	ExprType ExpressionType
	Generic  *GenericType
}

func (g *GenericType) String() string {
	if g.Generic == nil {
		return g.ExprType.String()
	}

	return fmt.Sprintf("%s[%s]", g.ExprType.String(), g.Generic.String())
}

// String возвращает строковое представление типа
func (et ExpressionType) String() string {
	switch et {
	case IntType:
		return "int"
	case FloatType:
		return "float"
	case BoolType:
		return "bool"
	case ArrayType:
		return "array"
	default:
		return "unknown"
	}
}
