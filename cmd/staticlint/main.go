package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/staticcheck"
)

var exitInMainAnalyzer analysis.Analyzer = analysis.Analyzer{
	Name: "exitinmain",
	Run:  checkExitInMain,
}

func main() {
	multichecker.Main(
		with(
			translate(
				staticcheck.Analyzers,
				func(a *lint.Analyzer) *analysis.Analyzer {
					return a.Analyzer
				},
			),
			loopclosure.Analyzer,
			printf.Analyzer,
			unusedresult.Analyzer,
			&exitInMainAnalyzer,
		)...,
	)
}

func translate[T any, R any](a []T, f func(T) R) []R {
	result := make([]R, len(a))
	for _, t := range a {
		result = append(result, f(t))
	}
	return result
}

func with[T any](source []T, others ...T) []T {
	result := make([]T, len(source))
	result = append(result, source...)
	result = append(result, others...)
	return result
}

func checkExitInMain(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					for _, statement := range x.Body.List {
						if callExpr, ok := statement.(*ast.ExprStmt); ok {
							if call, ok := callExpr.X.(*ast.CallExpr); ok {
								if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
									if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "os" && sel.Sel.Name == "Exit" {
										pass.Reportf(statement.Pos(), "os.Exit in main function")
									}
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
