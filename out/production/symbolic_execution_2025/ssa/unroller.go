package ssa

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"strconv"
)

func UnrollLoops(c *astutil.Cursor) bool {
	if forStmt, ok := c.Node().(*ast.ForStmt); ok {

		var stmts []ast.Stmt

		if forStmt.Init != nil {
			stmts = append(stmts, forStmt.Init)
		}

		if forStmt.Cond == nil {
			return true
		}
		iterations, canUnroll := getIterations(forStmt)

		if !canUnroll {
			return true
		}

		for i := 0; i < iterations; i++ {
			stmts = append(stmts, forStmt.Body.List...)

			if forStmt.Post != nil {
				stmts = append(stmts, forStmt.Post)
			}
		}

		unrolled := &ast.BlockStmt{
			Lbrace: forStmt.Body.Lbrace,
			Rbrace: forStmt.Body.Rbrace,
			List:   stmts,
		}

		c.Replace(unrolled)
	}
	return true
}

func getLit(expr ast.Expr) (int, bool) {
	init, ok := expr.(*ast.BasicLit)
	if !ok {
		return 0, false
	}

	counter, err := strconv.Atoi(init.Value)
	return counter, err != nil
}

func getIterations(stmt *ast.ForStmt) (int, bool) {
	ass, ok := stmt.Init.(*ast.AssignStmt)
	if !ok {
		return 0, false
	}

	counter, ok := getLit(ass.Rhs[0])
	if !ok {
		return 0, false
	}

	binaryExpr, ok := stmt.Cond.(*ast.BinaryExpr)
	if !ok {
		return 0, false
	}

	limit, ok := getLit(binaryExpr.Y)
	if !ok {
		return 0, false
	}

	switch stmt.Post.(type) {
	case *ast.IncDecStmt:
		switch binaryExpr.Op {
		case token.EQL:
			if limit == counter {
				return 1, true
			} else {
				return 0, true
			}
		case token.LSS:
			return limit - counter - 1, true
		case token.LEQ:
			return limit - counter, true
		default:
			return 0, false
		}
	}

	return 0, false
}
