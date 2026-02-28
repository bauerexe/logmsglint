package analysis

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"github.com/bauerexe/logmsglint/internal/domain"
	"golang.org/x/tools/go/analysis"
)

type ExtractedLogCall struct {
	Call       domain.LogCall
	MessagePos token.Pos
	MessageEnd token.Pos
}

func extractLogCall(pass *analysis.Pass, call *ast.CallExpr) (*ExtractedLogCall, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	method := sel.Sel.Name
	if !isSupportedMethod(method) || len(call.Args) == 0 {
		return nil, false
	}

	msgLit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || msgLit.Kind != token.STRING {
		return nil, false
	}

	message, err := strconv.Unquote(msgLit.Value)
	if err != nil {
		return nil, false
	}

	kind, ok := resolveLoggerKind(pass, sel)
	if !ok {
		return nil, false
	}

	return &ExtractedLogCall{
		Call:       domain.LogCall{Kind: kind, Method: method, Message: message},
		MessagePos: msgLit.Pos(),
		MessageEnd: msgLit.End(),
	}, true
}

func isSupportedMethod(method string) bool {
	switch method {
	case "Info", "Warn", "Error", "Debug":
		return true
	default:
		return false
	}
}

func resolveLoggerKind(pass *analysis.Pass, sel *ast.SelectorExpr) (domain.LoggerKind, bool) {
	if ident, ok := sel.X.(*ast.Ident); ok {
		if obj, ok := pass.TypesInfo.Uses[ident].(*types.PkgName); ok {
			if obj.Imported().Path() == "log/slog" {
				return domain.LoggerSlog, true
			}
		}
	}

	t := pass.TypesInfo.TypeOf(sel.X)
	if t == nil {
		return "", false
	}

	named := unwrapNamed(t)
	if named == nil || named.Obj() == nil || named.Obj().Pkg() == nil {
		return "", false
	}

	pkgPath := named.Obj().Pkg().Path()
	name := named.Obj().Name()

	if pkgPath == "log/slog" && name == "Logger" {
		return domain.LoggerSlog, true
	}

	if pkgPath == "go.uber.org/zap" && (name == "Logger" || name == "SugaredLogger") {
		return domain.LoggerZap, true
	}

	return "", false
}

func unwrapNamed(t types.Type) *types.Named {
	switch v := t.(type) {
	case *types.Named:
		return v
	case *types.Pointer:
		if named, ok := v.Elem().(*types.Named); ok {
			return named
		}
	}
	return nil
}
