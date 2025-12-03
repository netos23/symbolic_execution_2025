package main

import (
	"log"
	"symbolic-execution-course/internal/memory"
	"symbolic-execution-course/internal/symbolic"
	"symbolic-execution-course/internal/translator"
	"symbolic-execution-course/pkg/z3wrapper"

	"github.com/ebukreev/go-z3/z3"
)

func main() {
	var mem = memory.NewSymbolicMemory()
	var array = mem.Allocate(symbolic.ArrayType, "", &symbolic.GenericType{ExprType: symbolic.IntType})

	mem.AssignToArray(array, 5, symbolic.NewIntConstant(10))

	var fromArray = mem.GetFromArray(array, 5)
	println(fromArray.String())

	var anotherFromArray = mem.GetFromArray(array, 10)
	println(anotherFromArray.String())

	var obj = mem.Allocate(symbolic.ObjectType, "Foo", nil)

	var ref1 = mem.MakeRef(symbolic.ObjectType, "Foo", nil)
	var ref2 = mem.MakeRef(symbolic.ObjectType, "Foo", nil)

	mem.AssignField(ref1, 1, symbolic.NewIntConstant(10))
	mem.AssignField(ref2, 1, symbolic.NewIntConstant(5))

	tr := translator.NewZ3Translator()
	defer tr.Close()
	solver := z3wrapper.NewSolver()
	defer solver.Close()

	f1Expr, err := tr.TranslateExpression(mem.GetFieldValue(ref1, 1))
	if err != nil {
		log.Fatalf("Ошибка трансляции: %v", err)
	}

	f2Expr, err := tr.TranslateExpression(mem.GetFieldValue(ref2, 1))
	if err != nil {
		log.Fatalf("Ошибка трансляции: %v", err)
	}

	objExpr, err := tr.TranslateExpression(obj.Deref)
	if err != nil {
		log.Fatalf("Ошибка трансляции: %v", err)
	}
	ref1Expr, err := tr.TranslateExpression(ref1.Deref)
	if err != nil {
		log.Fatalf("Ошибка трансляции: %v", err)
	}
	ref2Expr, err := tr.TranslateExpression(ref2.Deref)
	if err != nil {
		log.Fatalf("Ошибка трансляции: %v", err)
	}

	solver.Assert((objExpr.(z3.Array)).Eq(ref1Expr.(z3.Array)))
	solver.Assert((objExpr.(z3.Array)).Eq(ref2Expr.(z3.Array)))
	solver.Assert((f1Expr.(z3.Int)).Eq(solver.Context().FromInt(5, solver.Context().IntSort()).(z3.Int)))
	solver.Assert((f2Expr.(z3.Int)).Eq(solver.Context().FromInt(5, solver.Context().IntSort()).(z3.Int)))

	sat, err := solver.Check()

	println(sat)

}
