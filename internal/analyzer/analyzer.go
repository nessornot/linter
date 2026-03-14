package analyzer

import (
	"go/ast"
	"regexp"
	"go/types"
	"strings"
	"go/token"
	"strconv"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// regexp for checks
var (
	// allow only eng letters, digits, spaсes, basic punctuation
	allowedCharsRx = regexp.MustCompile(`^[a-zA-Z0-9\s\-_:.,='"()\[\]{}/\\%!?]+$`)
	specialEndRx = regexp.MustCompile(`[!?.]$`)

	// key words
	sensitiveDataRx = regexp.MustCompile(`(password|token|api_key|secret)\s*[:=]`)
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "Checks logs for formatting and sensitive data",
	Run:      run,
	Requires:[]*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// taking only call expressions
	nodeFilter :=[]ast.Node{
		(*ast.CallExpr)(nil),
	}

	// iterating through nodes
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		// check if function call
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		// get info about function
		obj := pass.TypesInfo.Uses[sel.Sel]
		if obj == nil {
			return
		}
		
		funcDecl, ok := obj.(*types.Func)
		if !ok {
			return
		}

		// check if package == slog / zap
		pkg := funcDecl.Pkg()
		if pkg == nil {
			return
		}
		pkgPath := pkg.Path()
		isSlog := pkgPath == "log/slog"
		isZap := strings.HasPrefix(pkgPath, "go.uber.org/zap")

		// return if smth else
		if !isSlog && !isZap {
			return
		}

		// check if func name is {"Debug", "Info", "Warn" etc}
		methodName := funcDecl.Name()
		if !isLogMethod(methodName) {
			return
		}

		if len(call.Args) == 0 {
			return
		}
		
		// trying to find string literal
		msg, pos, ok := extractStrings(call.Args[0])
		if !ok || len(msg) == 0 {
			return
		}

		// --- RULES ---

		// 1. must start with lowercase letter
		firstRune := []rune(msg)[0]
		if unicode.IsLetter(firstRune) && unicode.IsUpper(firstRune) {
			pass.Reportf(pos, "log message must start with lowercase letter")
		}

		// 2. must include only eng letters, nums and basic punctuation
		if !allowedCharsRx.MatchString(msg) {
			pass.Reportf(pos, "log message must contain only english letters, numbers and basic punctuation (no emojis or special chars)")
		}

		// 3. should not end with punctuation marks
		if specialEndRx.MatchString(msg) {
			pass.Reportf(pos, "log message should not end with punctuation marks")
		}

		// 4. log message must not contain sensitive data
		lowerMsg := strings.ToLower(msg)
		if match := sensitiveDataRx.FindStringSubmatch(lowerMsg); match != nil {
			pass.Reportf(pos, "log message contains potentially sensitive data: %s", match[1])
		}
	})

	return nil, nil
}

// checks if func is used for logging
func isLogMethod(name string) bool {
	switch name {
	case "Debug", "Info", "Warn", "Error", "Fatal", "Panic", "Log":
		return true
	}
	return false
}

// get values from string literals
func extractStrings(expr ast.Expr) (string, token.Pos, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			val, err := strconv.Unquote(e.Value)
			if err == nil {
				return val, e.Pos(), true
			}
		}
	case *ast.BinaryExpr:
		// if its concatenation
		if e.Op == token.ADD {
			lStr, lPos, lOk := extractStrings(e.X)
			rStr, rPos, rOk := extractStrings(e.Y)

			pos := lPos
			if !lOk {
				pos = rPos
			}

			return lStr + rStr, pos, lOk || rOk
		}
	}
	return "", token.NoPos, false
}
