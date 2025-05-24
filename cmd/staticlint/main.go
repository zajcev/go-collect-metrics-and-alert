// Package main implements a static analysis tool that combines multiple checkers.
//
// This tool aggregates various static analyzers to perform comprehensive code analysis.
// It combines analyzers from standard packages, staticcheck, stylecheck, and custom analyzers.
//
// # Usage
//
// Build the tool:
//
//	go build -o staticlint cmd/staticlint/main.go
//
// Run analysis:
//
//	./staticlint ./... # Analyze all packages recursively
//	./staticlint package/path/... # Analyze specific package
//	./staticlint file.go # Analyze single file
//
// # Analyzers Included
//
// The tool combines the following analyzers:
//
// 1. Standard analyzers:
//   - printf: checks consistency of Printf format strings and arguments
//   - shadow: checks for possible unintended variable shadowing
//   - structtag: checks struct field tags are well formed
//
// 2. Custom analyzers:
//   - ExitCheckAnalyzer: detects direct calls to os.Exit in main package
//
// 3. Staticcheck analyzers (SA class):
//   - Includes all SA analyzers for bugs detection, performance issues, etc.
//
// 4. Selected Stylecheck analyzers (ST class):
//   - ST1001: checks for dot imports (import . "package")
//   - ST1011: checks for poorly named variables in range statements
//
// 5. Additional third-party analyzers:
//   - errcheck: checks for unchecked errors
//   - bodyclose: checks whether HTTP response body is closed
//
// # Multichecker Mechanism
//
// The multichecker:
// 1. Initializes all specified analyzers
// 2. Parses command-line flags for each analyzer
// 3. Loads package information and type data
// 4. Run each analyzer on the specified packages
// 5. Collects and reports diagnostics from all analyzers
// 6. Prints diagnostics to stdout
//
// # Configuration
//
// Analyzers can be configured via flags:
//
//	Run "./staticlint -help" to show all available flags for all analyzers.
//
// Shows all available flags for all analyzers.
package main

import (
	"github.com/kisielk/errcheck/errcheck"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/staticlint"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	analyzers := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		staticlint.ExitCheckAnalyzer,
	}
	// Add all analayzers from staticcheck class SA
	for _, analyzer := range staticcheck.Analyzers {
		analyzers = append(analyzers, analyzer.Analyzer)
	}
	//analayzers from stylecheck class ST
	for _, analyzer := range stylecheck.Analyzers {
		if analyzer.Analyzer.Name == "ST1001" || analyzer.Analyzer.Name == "ST1011" {
			analyzers = append(analyzers, analyzer.Analyzer)
		}
	}

	analyzers = append(analyzers,
		errcheck.Analyzer,
		bodyclose.Analyzer,
	)
	multichecker.Main(analyzers...)
}
