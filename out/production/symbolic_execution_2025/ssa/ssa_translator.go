package ssa

import (
	"go/types"
	"symbolic-execution-course/internal/symbolic"
)

func ConvertToSymbolic(param types.Type) (symbolic.ExpressionType, string, *symbolic.GenericType) {
	switch t := param.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int8:
			fallthrough
		case types.Int16:
			fallthrough
		case types.Int32:
			fallthrough
		case types.Int64:
			fallthrough
		case types.Uint8:
			fallthrough
		case types.Uint16:
			fallthrough
		case types.Uint32:
			fallthrough
		case types.Uint64:
			fallthrough
		case types.Int:
			fallthrough
		case types.Uint:
			return symbolic.IntType, "", nil
		case types.Bool:
			return symbolic.BoolType, "", nil
		case types.Float32:
			return symbolic.FloatType, "", nil
		default:
		}
	case *types.Slice:
		// TODO: Improve types for struct arrays
		elimType, _, elimGeneric := ConvertToSymbolic(t.Elem())
		return symbolic.ArrayType, "", &symbolic.GenericType{elimType, nil, elimGeneric}
	case *types.Named:
		return symbolic.ObjectType, t.Obj().Name(), nil
	case *types.Pointer:
		elimType, _, elimGeneric := ConvertToSymbolic(t.Elem())
		return symbolic.RefType, "", &symbolic.GenericType{elimType, nil, elimGeneric}
	}

	panic("Unsupported type")
}
