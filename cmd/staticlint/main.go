package staticlint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

// OsExitCheckAnalyzer - Анализатор поиска os.Exit
var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "os.Exit checker",
	Doc:  "check for os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				selexpr, ok := x.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				ident, ok := selexpr.X.(*ast.Ident)
				if !ok || ident.Name != "os" {
					return true
				}
				if selexpr.Sel.Name == "Exit" {
					return false
				}
			}
			return true
		})
	}
	return nil, nil
}

func main() {
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
		OsExitCheckAnalyzer,
	}
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	multichecker.Main(
		mychecks...,
	)
}
