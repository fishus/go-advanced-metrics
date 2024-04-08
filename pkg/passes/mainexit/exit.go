// Package mainexit Анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main.
package mainexit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "mainexit",
	Doc:  "Checking the call os.Exit in the main function.",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
skipGenerated:
	for _, file := range pass.Files {
		for _, cg := range file.Comments {
			for _, c := range cg.List {
				// Skip generated files
				if strings.Contains(c.Text, "DO NOT EDIT") {
					continue skipGenerated
				}
			}
		}

		if file.Name.Name != "main" {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			if fMain, ok := node.(*ast.FuncDecl); ok && fMain.Name.Name == "main" {
				for _, stmt := range fMain.Body.List {
					if isExitCall(stmt) {
						pass.Reportf(stmt.Pos(), "shouldn't call os.Exit in main function")
					}
				}
			}

			return true
		})
	}
	return nil, nil
}

func isExitCall(node ast.Stmt) bool {
	exprStmt, ok := node.(*ast.ExprStmt)
	if !ok {
		return false
	}

	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return false
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	if ident.Name == "os" && selExpr.Sel.Name == "Exit" {
		return true
	}

	return false
}
