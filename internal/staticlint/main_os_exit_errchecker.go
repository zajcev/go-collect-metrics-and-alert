package staticlint

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitCheck",
	Doc:  "check for os.Exit calls in main.go",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		if !strings.Contains(file.Name.Name, "main") {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			switch expr := call.Fun.(type) {
			case *ast.SelectorExpr:
				if pkgIdent, ok := expr.X.(*ast.Ident); ok && pkgIdent.Name == "os" && expr.Sel.Name == "Exit" {
					pass.Reportf(call.Pos(), "direct os.Exit() call in main package")
				}
			}
			return true
		})
	}
	return nil, nil
}
