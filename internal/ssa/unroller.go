package ssa

import (
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"strconv"
	"strings"
)

const defaultUnrollCount = 1

var breakLabelCounter int
var continueLabelCounter int

func UnrollLoops(c *astutil.Cursor) bool {
	if forStmt, ok := c.Node().(*ast.ForStmt); ok {

		breakLabelCounter++
		breakLabel := fmt.Sprintf("unrollBreakLabel%d", breakLabelCounter)

		//continueLabelCounter++
		//continueLabel := fmt.Sprintf("unrollContinueLabel%d", continueLabelCounter)

		var stmts []ast.Stmt

		if forStmt.Init != nil {
			stmts = append(stmts, forStmt.Init)
		}

		if forStmt.Cond == nil {
			return true
		}

		iterations, canUnroll := getIterations(forStmt)
		if !canUnroll {
			iterations = defaultUnrollCount
		}
		needLabel := false
		for i := 0; i < iterations; i++ {
			body, repl := replaceBreakWithGoto(forStmt.Body.List, breakLabel, token.BREAK)
			stmts = append(stmts, body...)
			needLabel = repl || needLabel

			body, repl = replaceBreakWithGoto(body, breakLabel, token.CONTINUE)
			//if repl {
			//	stmts = append(stmts, &ast.LabeledStmt{
			//		Label: ast.NewIdent(continueLabel),
			//		Stmt:  &ast.EmptyStmt{},
			//	})
			//	continueLabelCounter++
			//	continueLabel = fmt.Sprintf("unrollContinueLabel%d", continueLabelCounter)
			//}

			if forStmt.Post != nil {
				stmts = append(stmts, forStmt.Post)
			}
		}

		if needLabel {
			stmts = append(stmts, &ast.LabeledStmt{
				Label: ast.NewIdent(breakLabel),
				Stmt:  &ast.EmptyStmt{},
			})
		}

		unrolled := &ast.BlockStmt{
			Lbrace: forStmt.Body.Lbrace,
			Rbrace: forStmt.Body.Rbrace,
			List:   stmts,
		}

		c.Replace(unrolled)
	}

	if rangeStmt, ok := c.Node().(*ast.RangeStmt); ok {
		breakLabelCounter++
		breakLabel := fmt.Sprintf("unrollBreakLabel%d", breakLabelCounter)

		var stmts []ast.Stmt
		var preStmts []ast.Stmt

		/*iName := fmt.Sprintf("unrollRangeIdx%d", breakLabelCounter)
		iIdent := ast.NewIdent(iName)
		iInit := &ast.AssignStmt{
			Lhs: []ast.Expr{iIdent},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}},
		}
		preStmts = append(preStmts, iInit)*/

		// Определяем key и value
		var keyIdent, valueIdent *ast.Ident
		if rangeStmt.Key != nil {
			if ident, ok := rangeStmt.Key.(*ast.Ident); ok && ident.Name != "_" {
				keyIdent = ident
				keyStmt := &ast.AssignStmt{
					Lhs: []ast.Expr{keyIdent},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}},
				}
				preStmts = append(preStmts, keyStmt)
			}
		}
		if rangeStmt.Value != nil {
			if ident, ok := rangeStmt.Value.(*ast.Ident); ok && ident.Name != "_" {
				valueIdent = ident
			}
		}

		// Если используется value, объявляем переменную нужного типа
		if valueIdent != nil {
			decl := &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{valueIdent},
							Type:  nil, // Тип не всегда можно вывести, оставим nil
						},
					},
				},
			}
			preStmts = append(preStmts, decl)
		}

		needLabel := false
		for i := 0; i < defaultUnrollCount; i++ {
			var iterStmts []ast.Stmt

			// key = i
			if keyIdent != nil {
				assignKey := &ast.AssignStmt{
					Lhs: []ast.Expr{keyIdent},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(i)}},
				}
				iterStmts = append(iterStmts, assignKey)
			}

			// value = x[i]
			if valueIdent != nil {
				indexExpr := &ast.IndexExpr{
					X:     rangeStmt.X,
					Index: &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(i)},
				}
				assignValue := &ast.AssignStmt{
					Lhs: []ast.Expr{valueIdent},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{indexExpr},
				}
				iterStmts = append(iterStmts, assignValue)
			}

			// Клонируем тело цикла
			bodyList := cloneStmtList(rangeStmt.Body.List)

			// Заменяем break/continue на goto
			bodyList, replBreak := replaceBreakWithGoto(bodyList, breakLabel, token.BREAK)
			bodyList, _ = replaceBreakWithGoto(bodyList, breakLabel, token.CONTINUE)
			needLabel = needLabel || replBreak

			iterStmts = append(iterStmts, bodyList...)

			/*// i++
			incI := &ast.IncDecStmt{
				X:   iIdent,
				Tok: token.INC,
			}
			iterStmts = append(iterStmts, incI)*/

			stmts = append(stmts, iterStmts...)
		}

		if needLabel {
			stmts = append(stmts, &ast.LabeledStmt{
				Label: ast.NewIdent(breakLabel),
				Stmt:  &ast.EmptyStmt{},
			})
		}

		// Объединяем preStmts и stmts
		allStmts := append(preStmts, stmts...)

		unrolled := &ast.BlockStmt{
			Lbrace: rangeStmt.Body.Lbrace,
			Rbrace: rangeStmt.Body.Rbrace,
			List:   allStmts,
		}

		c.Replace(unrolled)
	}

	return true
}

func cloneStmtList(list []ast.Stmt) []ast.Stmt {
	out := make([]ast.Stmt, len(list))
	copy(out, list)
	return out
}

// Рекурсивно заменяет break на goto <label>
func replaceBreakWithGoto(stmts []ast.Stmt, label string, tok token.Token) ([]ast.Stmt, bool) {
	var out []ast.Stmt
	hasBreak := false
	inRepl := false
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.BranchStmt:
			if s.Tok == tok && s.Label == nil {
				out = append(out, &ast.BranchStmt{
					Tok:   token.GOTO,
					Label: ast.NewIdent(label),
				})
				hasBreak = true
				continue
			}

			if s.Label != nil && strings.HasPrefix(s.Label.Name, "unrollContinueLabel") && tok == token.CONTINUE {
				out = append(out, &ast.BranchStmt{
					Tok:   token.GOTO,
					Label: ast.NewIdent(label),
				})
				hasBreak = true
				continue
			}
		case *ast.BlockStmt:
			s.List, inRepl = replaceBreakWithGoto(s.List, label, tok)
			hasBreak = hasBreak || inRepl
			out = append(out, s)
			continue
		case *ast.IfStmt:
			s.Body.List, inRepl = replaceBreakWithGoto(s.Body.List, label, tok)
			hasBreak = hasBreak || inRepl
			if s.Else != nil {
				if elseBlock, ok := s.Else.(*ast.BlockStmt); ok {
					elseBlock.List, inRepl = replaceBreakWithGoto(elseBlock.List, label, tok)
					hasBreak = hasBreak || inRepl
				}
			}
			out = append(out, s)
			continue
		}
		out = append(out, stmt)
	}
	return out, hasBreak
}

func getLit(expr ast.Expr) (int, bool) {
	init, ok := expr.(*ast.BasicLit)
	if !ok {
		return 0, false
	}

	counter, err := strconv.Atoi(init.Value)
	return counter, err == nil
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
