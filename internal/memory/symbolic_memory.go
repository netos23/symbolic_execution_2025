package memory

import (
	"fmt"
	"symbolic-execution-course/internal/symbolic"
)

type Memory interface {
	Allocate(tpe symbolic.ExpressionType, typeName string, genericType *symbolic.GenericType) *symbolic.Ref

	MakeRef(tpe symbolic.ExpressionType, typeName string, genericType *symbolic.GenericType) *symbolic.Ref

	MakeRefRaw(expr symbolic.SymbolicExpression) *symbolic.Ref

	AssignPrimitive(ref *symbolic.Ref, value symbolic.SymbolicExpression) symbolic.SymbolicExpression

	ReadPrimitive(ref *symbolic.Ref) symbolic.SymbolicExpression

	AssignField(ref *symbolic.Ref, fieldIdx int, value symbolic.SymbolicExpression) symbolic.SymbolicExpression

	GetFieldValue(ref *symbolic.Ref, fieldIdx int) symbolic.SymbolicExpression

	AssignToArray(ref *symbolic.Ref, index symbolic.SymbolicExpression, value symbolic.SymbolicExpression) symbolic.SymbolicExpression

	GetFromArray(ref *symbolic.Ref, index symbolic.SymbolicExpression) symbolic.SymbolicExpression

	Assign(lhs symbolic.SymbolicExpression, rhs symbolic.SymbolicExpression) symbolic.SymbolicExpression
}

type PrimitiveHolder struct {
	RefSeq        int64
	PrimitiveType symbolic.ExpressionType
	Slots         symbolic.SymbolicExpression
}

func NewPrimitivesHolder(tpe symbolic.ExpressionType, slots symbolic.SymbolicExpression) *PrimitiveHolder {
	return &PrimitiveHolder{
		RefSeq:        1,
		PrimitiveType: tpe,
		Slots:         slots,
	}
}

type ObjectHolder struct {
	RefSeq       int64
	ObjectDef    *symbolic.Object
	FieldsHolder []symbolic.SymbolicExpression
}

func NewObjectHolder(obj *symbolic.Object) *ObjectHolder {
	return &ObjectHolder{
		RefSeq:       1,
		ObjectDef:    obj,
		FieldsHolder: make([]symbolic.SymbolicExpression, 0),
	}
}

type ArrayHolder struct {
	RefSeq  int64
	Generic *symbolic.GenericType
	Slots   symbolic.SymbolicExpression
}

func NewArrayHolder(generic *symbolic.GenericType) *ArrayHolder {
	return &ArrayHolder{
		RefSeq:  1,
		Generic: generic,
		Slots: symbolic.NewSymbolicVariable(
			fmt.Sprintf("$%s", generic.String()),
			symbolic.ArrayType,
			&symbolic.GenericType{
				ExprType: symbolic.ArrayType,
				Generic:  generic,
			},
			nil,
		),
	}
}

type SymbolicMemory struct {
	RefId         int64
	Refs          symbolic.SymbolicExpression
	PrimitivePool map[symbolic.ExpressionType]*PrimitiveHolder
	ObjectPool    map[string]*ObjectHolder
	ArrayPool     map[string]*ArrayHolder
}

func NewSymbolicMemory() *SymbolicMemory {
	return &SymbolicMemory{
		RefId: 1,
		Refs: symbolic.NewSymbolicVariable(
			"$addr", symbolic.ArrayType,
			&symbolic.GenericType{symbolic.IntType, nil, nil},
			nil,
		),
		PrimitivePool: map[symbolic.ExpressionType]*PrimitiveHolder{
			symbolic.IntType: NewPrimitivesHolder(
				symbolic.IntType, symbolic.NewSymbolicVariable(
					"$addrI", symbolic.ArrayType,
					&symbolic.GenericType{symbolic.IntType, nil, nil},
					nil,
				),
			),
			symbolic.FloatType: NewPrimitivesHolder(
				symbolic.FloatType, symbolic.NewSymbolicVariable(
					"$addrF", symbolic.ArrayType,
					&symbolic.GenericType{symbolic.FloatType, nil, nil},
					nil,
				),
			),
			symbolic.BoolType: NewPrimitivesHolder(
				symbolic.BoolType, symbolic.NewSymbolicVariable(
					"$addrB", symbolic.ArrayType,
					&symbolic.GenericType{symbolic.BoolType, nil, nil},
					nil,
				),
			),
		},
		ObjectPool: make(map[string]*ObjectHolder),
		ArrayPool:  make(map[string]*ArrayHolder),
	}
}

func (mem *SymbolicMemory) Allocate(tpe symbolic.ExpressionType, typeName string, genericType *symbolic.GenericType) *symbolic.Ref {
	switch tpe {
	case symbolic.BoolType, symbolic.IntType, symbolic.FloatType:
		holder, _ := mem.PrimitivePool[tpe]
		refId := holder.RefSeq
		holder.RefSeq += 1
		address := mem.RefId
		mem.RefId += 1
		deref := symbolic.NewArrayStore(mem.Refs, symbolic.NewIntConstant(address), symbolic.NewIntConstant(refId))
		mem.Refs = deref

		return symbolic.NewRef(
			address, tpe, nil, nil,
			symbolic.NewArraySelect(
				mem.PrimitivePool[tpe].Slots,
				symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(address)),
			),
		)
	case symbolic.ArrayType:
		holder, hasHolder := mem.ArrayPool[genericType.String()]
		if !hasHolder {
			holder = NewArrayHolder(genericType)
			mem.ArrayPool[genericType.String()] = holder
		}

		refId := holder.RefSeq
		holder.RefSeq += 1
		address := mem.RefId
		mem.RefId += 1
		deref := symbolic.NewArrayStore(mem.Refs, symbolic.NewIntConstant(address), symbolic.NewIntConstant(refId))
		mem.Refs = deref

		return symbolic.NewRef(
			address, tpe, genericType, nil,
			symbolic.NewArraySelect(holder.Slots, symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(address))),
		)
	case symbolic.ObjectType:
		holder, hasHolder := mem.ObjectPool[typeName]
		if !hasHolder {
			holder = NewObjectHolder(symbolic.NewObject(typeName))
			mem.ObjectPool[typeName] = holder
		}

		if genericType != nil && genericType.ObjectType != nil {
			typ := holder.ObjectDef
			target := genericType.ObjectType
			if cap(typ.Fields) <= len(target.Fields) {
				tmp := make([]*symbolic.ObjectField, len(target.Fields)+1)
				copy(tmp, typ.Fields)
				typ.Fields = tmp

				aux := make([]symbolic.SymbolicExpression, len(target.Fields)+1)
				copy(aux, holder.FieldsHolder)
				holder.FieldsHolder = aux
			}

			for fieldIdx, f := range target.Fields {
				if typ.Fields[fieldIdx] == nil {
					switch f.ExprType {
					case symbolic.IntType, symbolic.FloatType, symbolic.BoolType:
						typ.Fields[fieldIdx] = symbolic.NewObjectField(f.ExprType, nil, nil)
						holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
							fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, f.ExprType),
							symbolic.ArrayType, &symbolic.GenericType{ExprType: f.ExprType},
							nil,
						)
						break
					case symbolic.ArrayType:
						typ.Fields[fieldIdx] = symbolic.NewObjectField(f.ExprType, nil, f.Generic)
						holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
							fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, f.ExprType),
							symbolic.ArrayType,
							&symbolic.GenericType{ExprType: f.ExprType, Generic: f.Generic},
							nil,
						)
						break
					case symbolic.RefType:
						fallthrough
					case symbolic.ObjectType:
						typ.Fields[fieldIdx] = symbolic.NewObjectField(symbolic.RefType, f.ObjectType, nil)
						holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
							fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, f.ExprType),
							symbolic.ArrayType,
							&symbolic.GenericType{symbolic.RefType, f.ObjectType, nil},
							nil,
						)
						break
					}
				}
			}
		}

		refId := holder.RefSeq
		holder.RefSeq += 1
		address := mem.RefId
		mem.RefId += 1
		deref := symbolic.NewArrayStore(mem.Refs, symbolic.NewIntConstant(address), symbolic.NewIntConstant(refId))
		mem.Refs = deref

		return symbolic.NewRef(
			address, tpe, nil, holder.ObjectDef,
			symbolic.NewArraySelect(deref, symbolic.NewIntConstant(address)),
		)
	case symbolic.RefType:

	}

	panic("Allocate ref unsupported")
}

func (mem *SymbolicMemory) MakeRef(tpe symbolic.ExpressionType, typeName string, genericType *symbolic.GenericType) *symbolic.Ref {
	switch tpe {
	case symbolic.IntType, symbolic.FloatType, symbolic.BoolType:
		address := mem.RefId
		mem.RefId += 1

		return symbolic.NewRef(
			address, tpe, nil, nil,
			symbolic.NewArraySelect(
				mem.PrimitivePool[tpe].Slots,
				symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(address)),
			),
		)
	case symbolic.ArrayType:
		holder, hasHolder := mem.ArrayPool[genericType.String()]
		if !hasHolder {
			holder = NewArrayHolder(genericType)
			mem.ArrayPool[genericType.String()] = holder
		}
		address := mem.RefId
		mem.RefId += 1

		return symbolic.NewRef(
			address, tpe, genericType, nil,
			symbolic.NewArraySelect(holder.Slots, symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(address))),
		)
	case symbolic.ObjectType:
		holder, hasHolder := mem.ObjectPool[typeName]
		if !hasHolder {
			holder = NewObjectHolder(symbolic.NewObject(typeName))
			mem.ObjectPool[typeName] = holder
		}

		if genericType != nil && genericType.ObjectType != nil {
			typ := holder.ObjectDef
			target := genericType.ObjectType
			if cap(typ.Fields) <= len(target.Fields) {
				tmp := make([]*symbolic.ObjectField, len(target.Fields)+1)
				copy(tmp, typ.Fields)
				typ.Fields = tmp

				aux := make([]symbolic.SymbolicExpression, len(target.Fields)+1)
				copy(aux, holder.FieldsHolder)
				holder.FieldsHolder = aux
			}

			for fieldIdx, f := range target.Fields {
				if typ.Fields[fieldIdx] == nil {
					switch f.ExprType {
					case symbolic.IntType, symbolic.FloatType, symbolic.BoolType:
						typ.Fields[fieldIdx] = symbolic.NewObjectField(f.ExprType, nil, nil)
						holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
							fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, f.ExprType),
							symbolic.ArrayType, &symbolic.GenericType{ExprType: f.ExprType},
							nil,
						)
						break
					case symbolic.ArrayType:
						typ.Fields[fieldIdx] = symbolic.NewObjectField(f.ExprType, nil, f.Generic)
						holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
							fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, f.ExprType),
							symbolic.ArrayType,
							&symbolic.GenericType{ExprType: f.ExprType, Generic: f.Generic},
							nil,
						)
						break
					case symbolic.RefType:
						fallthrough
					case symbolic.ObjectType:
						typ.Fields[fieldIdx] = symbolic.NewObjectField(symbolic.RefType, f.ObjectType, nil)
						holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
							fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, f.ExprType),
							symbolic.ArrayType,
							&symbolic.GenericType{symbolic.RefType, f.ObjectType, nil},
							nil,
						)
						break
					}
				}
			}
		}

		address := mem.RefId
		mem.RefId += 1

		return symbolic.NewRef(
			address, tpe, nil, holder.ObjectDef,
			symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(address)),
		)
	case symbolic.RefType:
		// TODO: work with pointers
	}

	panic("Make ref unsupported")
}

func (mem *SymbolicMemory) MakeRefRaw(expr symbolic.SymbolicExpression) *symbolic.Ref {
	var typ symbolic.ExpressionType
	var gen *symbolic.GenericType
	var obj *symbolic.Object
	objName := ""

	switch rawRef := expr.(type) {
	case *symbolic.ArraySelect:
		typ = expr.Type()
		gen = symbolic.GenericFor(rawRef)
		pGen := symbolic.GenericFor(rawRef.Array)

		if (typ == symbolic.ObjectType || typ == symbolic.RefType) && pGen.ObjectType != nil {
			objName = pGen.ObjectType.Name
		}
	default:
		typ = expr.Type()
		gen = symbolic.GenericFor(expr)
	}

	if holder, ok := mem.ObjectPool[objName]; ok {
		obj = holder.ObjectDef
	}

	return symbolic.NewRef(int64(-1), typ, gen, obj, expr)
}

func (mem *SymbolicMemory) AssignField(ref *symbolic.Ref, fieldIdx int, value symbolic.SymbolicExpression) symbolic.SymbolicExpression {
	holder, _ := mem.ObjectPool[ref.ObjType.Name]
	typ := holder.ObjectDef

	if cap(typ.Fields) <= fieldIdx {
		tmp := make([]*symbolic.ObjectField, fieldIdx+1)
		copy(tmp, typ.Fields)
		typ.Fields = tmp

		aux := make([]symbolic.SymbolicExpression, fieldIdx+1)
		copy(aux, holder.FieldsHolder)
		holder.FieldsHolder = aux
	}

	if typ.Fields[fieldIdx] == nil {
		switch value.Type() {
		case symbolic.IntType, symbolic.FloatType, symbolic.BoolType:
			typ.Fields[fieldIdx] = symbolic.NewObjectField(value.Type(), nil, nil)
			holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
				fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, value.Type()),
				symbolic.ArrayType, &symbolic.GenericType{ExprType: value.Type()},
				nil,
			)
			break
		case symbolic.ArrayType:
			typ.Fields[fieldIdx] = symbolic.NewObjectField(value.Type(), nil, symbolic.GenericFor(value))
			holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
				fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, value.Type()),
				symbolic.ArrayType,
				&symbolic.GenericType{ExprType: value.Type(), Generic: symbolic.GenericFor(value)},
				nil,
			)
			break
		case symbolic.RefType:
			fallthrough
		case symbolic.ObjectType:
			typ.Fields[fieldIdx] = symbolic.NewObjectField(symbolic.RefType, symbolic.ObjectFor(value), nil)
			holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
				fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, value.Type()),
				symbolic.ArrayType,
				&symbolic.GenericType{symbolic.RefType, symbolic.ObjectFor(value), nil},
				nil,
			)
			break
		}
	}

	f := holder.FieldsHolder[fieldIdx]

	var addr symbolic.SymbolicExpression
	if ref.Address == int64(-1) {
		addr = ref.Deref
	} else {
		addr = symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(ref.Address))
	}

	f = symbolic.NewArrayStore(f, addr, value)
	holder.FieldsHolder[fieldIdx] = f

	return symbolic.NewFieldWrite(
		ref,
		fieldIdx,
		value,
		f,
	)
}

func (mem *SymbolicMemory) GetFieldValue(ref *symbolic.Ref, fieldIdx int) symbolic.SymbolicExpression {
	holder, _ := mem.ObjectPool[ref.ObjType.Name]
	typ := holder.ObjectDef

	if cap(typ.Fields) <= fieldIdx {
		tmp := make([]*symbolic.ObjectField, fieldIdx+1)
		copy(tmp, typ.Fields)
		typ.Fields = tmp

		aux := make([]symbolic.SymbolicExpression, fieldIdx+1)
		copy(aux, holder.FieldsHolder)
		holder.FieldsHolder = aux
	}

	/*if typ.Fields[fieldIdx] == nil {
		switch value.Type() {
		case symbolic.IntType, symbolic.FloatType, symbolic.BoolType:
			typ.Fields[fieldIdx] = symbolic.NewObjectField(value.Type(), nil, nil)
			holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
				fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, value.Type()),
				symbolic.ArrayType, &symbolic.GenericType{ExprType: value.Type()},
				nil,
			)
			break
		case symbolic.ArrayType:
			typ.Fields[fieldIdx] = symbolic.NewObjectField(value.Type(), nil, symbolic.GenericFor(value))
			holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
				fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, value.Type()),
				symbolic.ArrayType,
				&symbolic.GenericType{ExprType: value.Type(), Generic: symbolic.GenericFor(value)},
				nil,
			)
			break
		case symbolic.RefType:
			fallthrough
		case symbolic.ObjectType:
			typ.Fields[fieldIdx] = symbolic.NewObjectField(symbolic.RefType, symbolic.ObjectFor(value), nil)
			holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
				fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, value.Type()),
				symbolic.ArrayType,
				&symbolic.GenericType{symbolic.RefType, symbolic.ObjectFor(value), nil},
				nil,
			)
			break
		}
	}*/

	f := holder.FieldsHolder[fieldIdx]

	var addr symbolic.SymbolicExpression
	if ref.Address == int64(-1) {
		addr = ref.Deref
	} else {
		addr = symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(ref.Address))
	}

	return symbolic.NewFieldRead(
		ref,
		fieldIdx,
		symbolic.NewArraySelect(
			f,
			addr,
		),
	)
}

func (mem *SymbolicMemory) AssignToArray(ref *symbolic.Ref, index symbolic.SymbolicExpression, value symbolic.SymbolicExpression) symbolic.SymbolicExpression {
	holder, _ := mem.ArrayPool[ref.TypeGeneric.String()]

	var addr symbolic.SymbolicExpression
	if ref.Address == int64(-1) {
		addr = ref.Deref
	} else {
		addr = symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(ref.Address))
	}

	holder.Slots = symbolic.NewArrayStore(
		holder.Slots,
		addr,
		symbolic.NewArrayStore(
			symbolic.NewArraySelect(holder.Slots, addr),
			index,
			value,
		),
	)

	return holder.Slots
}

func (mem *SymbolicMemory) GetFromArray(ref *symbolic.Ref, index symbolic.SymbolicExpression) symbolic.SymbolicExpression {
	holder, _ := mem.ArrayPool[ref.TypeGeneric.String()]

	var addr symbolic.SymbolicExpression
	if ref.Address == int64(-1) {
		addr = ref.Deref
	} else {
		addr = symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(ref.Address))
	}

	return symbolic.NewArraySelect(
		symbolic.NewArraySelect(holder.Slots, addr),
		index,
	)
}
func (mem *SymbolicMemory) AssignPrimitive(ref *symbolic.Ref, value symbolic.SymbolicExpression) symbolic.SymbolicExpression {
	holder, _ := mem.PrimitivePool[ref.VarType]

	var addr symbolic.SymbolicExpression
	if ref.Address == int64(-1) {
		addr = ref.Deref
	} else {
		addr = symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(ref.Address))
	}

	holder.Slots = symbolic.NewArrayStore(
		holder.Slots,
		addr,
		value,
	)

	return holder.Slots
}

func (mem *SymbolicMemory) ReadPrimitive(ref *symbolic.Ref) symbolic.SymbolicExpression {
	holder, _ := mem.PrimitivePool[ref.VarType]

	return symbolic.NewArraySelect(
		holder.Slots,
		symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(ref.Address)),
	)
}
func (mem *SymbolicMemory) Assign(lhs symbolic.SymbolicExpression, rhs symbolic.SymbolicExpression) symbolic.SymbolicExpression {
	if _, ok := lhs.(*symbolic.Ref); !ok {
		return rhs
	}

	ref := lhs.(*symbolic.Ref)

	switch ref.VarType {
	case symbolic.BoolType, symbolic.IntType, symbolic.FloatType:
		mem.AssignPrimitive(ref, rhs)

	default:
		// TODO: pointers support
	}

	return lhs
}
