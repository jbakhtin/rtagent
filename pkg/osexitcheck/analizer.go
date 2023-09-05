// Package osexitcheck пакет производит проверку вызова функции os.Exit в функции main
package osexitcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for call os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	isMainFunction := func(x *ast.FuncDecl) bool {
		return x.Name.Name == "main"
	}

	findOsExit := func(x *ast.CallExpr) {
		if call, ok := x.Fun.(*ast.SelectorExpr); ok {
			pkgIdent, ok := call.X.(*ast.Ident)
			if ok && pkgIdent.Name == "os" {
				if ok && call.Sel.Name == "Exit" {
					pass.Reportf(x.Pos(), "os.Exit call detected")
				}
			}
		}
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if isMainFunction(x) {
					ast.Inspect(node, func(node ast.Node) bool {
						switch x := node.(type) {
						case *ast.CallExpr:
							findOsExit(x)
						}
						return true
					})
				}
			}
			return true
		})
	}
	return nil, nil
}
