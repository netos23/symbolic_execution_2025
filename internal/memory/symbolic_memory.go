package memory

import (
	"fmt"
	"symbolic-execution-course/internal/symbolic"
)

type Memory interface {
	Allocate(tpe symbolic.ExpressionType, typeName string, genericType *symbolic.GenericType) *symbolic.Ref

	MakeRef(tpe symbolic.ExpressionType, typeName string, genericType *symbolic.GenericType) *symbolic.Ref

	AssignPrimitive(ref *symbolic.Ref, value symbolic.SymbolicExpression) symbolic.SymbolicExpression

	ReadPrimitive(ref *symbolic.Ref) symbolic.SymbolicExpression

	AssignField(ref *symbolic.Ref, fieldIdx int, value symbolic.SymbolicExpression) symbolic.SymbolicExpression

	GetFieldValue(ref *symbolic.Ref, fieldIdx int) symbolic.SymbolicExpression

	AssignToArray(ref *symbolic.Ref, index int64, value symbolic.SymbolicExpression) symbolic.SymbolicExpression

	GetFromArray(ref *symbolic.Ref, index int64) symbolic.SymbolicExpression
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
				symbolic.ArrayType,
				nil,
				generic,
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
					"$addr", symbolic.ArrayType,
					&symbolic.GenericType{symbolic.IntType, nil, nil},
					nil,
				),
			),
			symbolic.FloatType: NewPrimitivesHolder(
				symbolic.FloatType, symbolic.NewSymbolicVariable(
					"$addr", symbolic.ArrayType,
					&symbolic.GenericType{symbolic.FloatType, nil, nil},
					nil,
				),
			),
			symbolic.BoolType: NewPrimitivesHolder(
				symbolic.BoolType, symbolic.NewSymbolicVariable(
					"$addr", symbolic.ArrayType,
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
	case symbolic.IntType:
	case symbolic.FloatType:
	case symbolic.BoolType:
		holder, _ := mem.PrimitivePool[tpe]
		refId := holder.RefSeq
		holder.RefSeq += 1
		address := mem.RefId
		mem.RefId += 1
		deref := symbolic.NewArrayStore(mem.Refs, symbolic.NewIntConstant(address), symbolic.NewIntConstant(refId))
		mem.Refs = deref

		return symbolic.NewRef(
			address, tpe, nil, nil,
			symbolic.NewArraySelect(deref, symbolic.NewIntConstant(address)),
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
			symbolic.NewArraySelect(deref, symbolic.NewIntConstant(address)),
		)
	case symbolic.ObjectType:
		holder, hasHolder := mem.ObjectPool[typeName]
		if !hasHolder {
			holder = NewObjectHolder(symbolic.NewObject(typeName))
			mem.ObjectPool[typeName] = holder
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
	case symbolic.IntType:
	case symbolic.FloatType:
	case symbolic.BoolType:
		address := mem.RefId
		mem.RefId += 1

		return symbolic.NewRef(
			address, tpe, nil, nil,
			symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(address)),
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
			symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(address)),
		)
	case symbolic.ObjectType:
		holder, hasHolder := mem.ObjectPool[typeName]
		if !hasHolder {
			holder = NewObjectHolder(symbolic.NewObject(typeName))
			mem.ObjectPool[typeName] = holder
		}

		address := mem.RefId
		mem.RefId += 1

		return symbolic.NewRef(
			address, tpe, nil, holder.ObjectDef,
			symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(address)),
		)
	case symbolic.RefType:

	}

	panic("Make ref unsupported")
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
		case symbolic.IntType:
			fallthrough
		case symbolic.FloatType:
			fallthrough
		case symbolic.BoolType:
			typ.Fields[fieldIdx] = symbolic.NewObjectField(value.Type(), nil, nil)
			holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
				fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, value.Type()),
				symbolic.ArrayType, &symbolic.GenericType{value.Type(), nil, nil},
				nil,
			)
			break
		case symbolic.ArrayType:
			typ.Fields[fieldIdx] = symbolic.NewObjectField(value.Type(), nil, symbolic.GenericFor(value))
			holder.FieldsHolder[fieldIdx] = symbolic.NewSymbolicVariable(
				fmt.Sprintf("$%s%d$%s", typ.Name, fieldIdx, value.Type()),
				symbolic.ArrayType,
				&symbolic.GenericType{value.Type(), nil, symbolic.GenericFor(value)},
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

	f = symbolic.NewArrayStore(f, symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(ref.Address)), value)
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

	f := holder.FieldsHolder[fieldIdx]

	return symbolic.NewFieldRead(
		ref,
		fieldIdx,
		symbolic.NewArraySelect(
			f,
			symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(ref.Address)),
		),
	)
}

func (mem *SymbolicMemory) AssignToArray(ref *symbolic.Ref, index int64, value symbolic.SymbolicExpression) symbolic.SymbolicExpression {
	holder, _ := mem.ArrayPool[ref.TypeGeneric.String()]

	holder.Slots = symbolic.NewArrayStore(
		symbolic.NewArraySelect(holder.Slots, symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(ref.Address))),
		symbolic.NewIntConstant(index),
		value,
	)

	return holder.Slots
}

func (mem *SymbolicMemory) GetFromArray(ref *symbolic.Ref, index int64) symbolic.SymbolicExpression {
	holder, _ := mem.ArrayPool[ref.TypeGeneric.String()]

	return symbolic.NewArraySelect(
		symbolic.NewArraySelect(holder.Slots, symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(ref.Address))),
		symbolic.NewIntConstant(index),
	)
}
func (mem *SymbolicMemory) AssignPrimitive(ref *symbolic.Ref, value symbolic.SymbolicExpression) symbolic.SymbolicExpression {
	holder, _ := mem.PrimitivePool[ref.VarType]

	holder.Slots = symbolic.NewArrayStore(
		holder.Slots,
		symbolic.NewArraySelect(mem.Refs, symbolic.NewIntConstant(ref.Address)),
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
