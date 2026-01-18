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
	ObjectType
	RefType
	// Добавьте другие типы по необходимости
)

type ObjectField struct {
	ExprType   ExpressionType
	ObjectType *Object
	Generic    *GenericType
}

func NewObjectField(
	exprType ExpressionType,
	objectType *Object,
	generic *GenericType,
) *ObjectField {
	return &ObjectField{
		exprType,
		objectType,
		generic,
	}
}

type Object struct {
	Name   string
	Fields []*ObjectField
}

func NewObject(name string) *Object {
	return &Object{
		Name:   name,
		Fields: make([]*ObjectField, 0),
	}
}

type GenericType struct {
	ExprType   ExpressionType
	ObjectType *Object
	Generic    *GenericType
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
