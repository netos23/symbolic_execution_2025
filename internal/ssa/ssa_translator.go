package ssa

import (
	"go/types"
	"symbolic-execution-course/internal/symbolic"
)

func ConvertToSymbolic(param types.Type, args ...int) (symbolic.ExpressionType, string, *symbolic.GenericType) {
	depth := 0
	if len(args) > 0 {
		depth = args[0]
	}

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
		case types.Float32, types.Float64, types.UntypedFloat:
			return symbolic.FloatType, "", nil
		case types.String:
			return symbolic.ArrayType, "", &symbolic.GenericType{ExprType: symbolic.IntType}
		default:
		}
	case *types.Slice:
		elimType, name, elimGeneric := ConvertToSymbolic(t.Elem(), depth)
		return symbolic.ArrayType, "", &symbolic.GenericType{ExprType: elimType, ObjectType: symbolic.NewObject(name), Generic: elimGeneric}
	case *types.Array:
		elimType, name, elimGeneric := ConvertToSymbolic(t.Elem(), depth)

		return symbolic.ArrayType, "", &symbolic.GenericType{ExprType: elimType, ObjectType: symbolic.NewObject(name), Generic: elimGeneric}
	case *types.Named:

		var gen *symbolic.GenericType
		name := t.Obj().Name()

		if st, ok := t.Underlying().(*types.Struct); ok && depth < 3 {
			obj := symbolic.NewObject(name)
			obj.Fields = make([]*symbolic.ObjectField, st.NumFields())

			for i := 0; i < st.NumFields(); i++ {

				fType, fName, fGen := ConvertToSymbolic(st.Field(i).Type(), depth+1)
				if (fType != symbolic.ObjectType) && fType != symbolic.RefType {
					obj.Fields[i] = symbolic.NewObjectField(fType, nil, fGen)
				} else {
					var fObj *symbolic.Object
					if fGen == nil {
						fObj = symbolic.NewObject(fName)
					} else {
						fObj = fGen.ObjectType
					}

					obj.Fields[i] = symbolic.NewObjectField(fType, fObj, nil)
				}
			}

			gen = &symbolic.GenericType{
				ObjectType: obj,
			}
		}

		return symbolic.ObjectType, name, gen
	case *types.Pointer:
		return ConvertToSymbolic(t.Elem(), depth)

		//elimType, _, elimGeneric := ConvertToSymbolic(t.Elem())
		//return symbolic.RefType, "", &symbolic.GenericType{elimType, nil, elimGeneric}
	}

	panic("Unsupported type")
}

func DefaultSymbolic(param types.Type) symbolic.SymbolicExpression {
	typ, name, generic := ConvertToSymbolic(param)

	switch typ {
	case symbolic.IntType:
		return symbolic.NewIntConstant(0)
	case symbolic.BoolType:
		return symbolic.NewBoolConstant(false)
	case symbolic.FloatType:
		return symbolic.NewFloatConstant(0)
	default:
		return symbolic.NewSymbolicVariable("$nil"+param.String(), typ, generic, symbolic.NewObject(name))
	}
}
