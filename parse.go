package sl

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/antlr4-go/antlr/v4"
	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/parser"
)

// ParseError represents a parsing error
type ParseError struct {
	Message string
	Line    int
	Column  int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// ErrorListener implements the antlr.ErrorListener interface
type ErrorListener struct {
	Errors []ParseError
}

func (e *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, ex antlr.RecognitionException) {
	e.Errors = append(e.Errors, ParseError{
		Message: msg,
		Line:    line,
		Column:  column,
	})
}

func (e *ErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	// Can choose to handle ambiguity
}

func (e *ErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	// Can choose to handle full context attempts
}

func (e *ErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, prediction int, configs *antlr.ATNConfigSet) {
	// Can choose to handle context sensitivity
}

// ASTBuilder implements parser.BaseSLVisitor to build AST
type ASTBuilder struct {
	parser.BaseSLVisitor
}

// Parse parses an expression string and returns an AST
func Parse(expression string) (ast.ASTNode, error) {
	// Create input stream
	input := antlr.NewInputStream(expression)

	// Create lexer
	lexer := parser.NewSLLexer(input)

	// Create token stream
	stream := antlr.NewCommonTokenStream(lexer, 0)

	// Create parser
	p := parser.NewSLParser(stream)

	// Add error listener
	errorListener := &ErrorListener{}
	p.RemoveErrorListeners()
	p.AddErrorListener(errorListener)

	// Parse
	tree := p.Start_()

	// Check for parsing errors
	if len(errorListener.Errors) > 0 {
		return nil, &errorListener.Errors[0]
	}

	// Build AST
	builder := &ASTBuilder{}
	result := builder.Visit(tree)

	switch result := result.(type) {
	case ast.ASTNode:
		return result, nil
	case error:
		return nil, result
	default:
		return nil, fmt.Errorf("failed to build AST")
	}
}

// Visit dispatches to specific visit methods based on node type
func (v *ASTBuilder) Visit(tree antlr.ParseTree) interface{} {
	switch t := tree.(type) {
	case parser.IStartContext:
		return v.VisitStart(t)
	case parser.IExprContext:
		return v.VisitExpr(t)
	case parser.IConditionalOrContext:
		return v.VisitConditionalOr(t)
	case parser.IConditionalAndContext:
		return v.VisitConditionalAnd(t)
	case parser.IRelationContext:
		return v.VisitRelation(t)
	case parser.ICalcContext:
		return v.VisitCalc(t)
	case parser.IUnaryContext:
		return v.VisitUnary(t)
	case parser.IMemberContext:
		return v.VisitMember(t)
	case parser.IPrimaryContext:
		return v.VisitPrimary(t)
	case parser.ILiteralContext:
		return v.VisitLiteral(t)
	default:
		return v.VisitChildren(tree.(antlr.RuleNode))
	}
}

func (v *ASTBuilder) VisitStart(ctx parser.IStartContext) interface{} {
	if ctx.GetE() != nil {
		return v.Visit(ctx.GetE())
	}
	return nil
}

func (v *ASTBuilder) VisitExpr(ctx parser.IExprContext) interface{} {
	condition := v.Visit(ctx.GetE()).(ast.ASTNode)

	if ctx.GetOp() != nil && ctx.GetE1() != nil && ctx.GetE2() != nil {
		// Conditional expression condition ? trueExpr : falseExpr
		trueExpr := v.Visit(ctx.GetE1()).(ast.ASTNode)
		falseExpr := v.Visit(ctx.GetE2()).(ast.ASTNode)
		return ast.NewConditional(condition, trueExpr, falseExpr)
	}

	return condition
}

func (v *ASTBuilder) VisitConditionalOr(ctx parser.IConditionalOrContext) interface{} {
	left := v.Visit(ctx.GetE()).(ast.ASTNode)

	if len(ctx.GetOps()) > 0 {
		for i, op := range ctx.GetOps() {
			right := v.Visit(ctx.GetE1()[i]).(ast.ASTNode)

			var functionName string
			if op.GetText() == "||" {
				functionName = ast.LogicalOr
			} else if op.GetText() == "&&" {
				functionName = ast.LogicalAnd
			}

			left = ast.NewFunctionCall(ast.NewIdent(functionName, false), []ast.ASTNode{left, right})
		}
	}

	return left
}

func (v *ASTBuilder) VisitConditionalAnd(ctx parser.IConditionalAndContext) interface{} {
	left := v.Visit(ctx.GetE()).(ast.ASTNode)

	if len(ctx.GetOps()) > 0 {
		for i, op := range ctx.GetOps() {
			right := v.Visit(ctx.GetE1()[i]).(ast.ASTNode)

			var functionName string
			if op.GetText() == "||" {
				functionName = ast.LogicalOr
			} else if op.GetText() == "&&" {
				functionName = ast.LogicalAnd
			}

			left = ast.NewFunctionCall(ast.NewIdent(functionName, false), []ast.ASTNode{left, right})
		}
	}

	return left
}

func (v *ASTBuilder) VisitRelation(ctx parser.IRelationContext) interface{} {
	if ctx.Calc() != nil {
		return v.Visit(ctx.Calc())
	}

	// Binary relation operations
	if ctx.GetChildCount() == 3 {
		left := v.Visit(ctx.GetChild(0).(antlr.ParseTree)).(ast.ASTNode)
		right := v.Visit(ctx.GetChild(2).(antlr.ParseTree)).(ast.ASTNode)

		op := ctx.GetOp().GetText()
		var functionName string
		if op == "==" {
			functionName = ast.Equals
		} else if op == "!=" {
			functionName = ast.NotEquals
		} else if op == "<" {
			functionName = ast.Less
		} else if op == "<=" {
			functionName = ast.LessEquals
		} else if op == ">" {
			functionName = ast.Greater
		} else if op == ">=" {
			functionName = ast.GreaterEquals
		} else if op == "in" {
			functionName = ast.In
		}

		return ast.NewFunctionCall(ast.NewIdent(functionName, false), []ast.ASTNode{left, right})
	}

	return nil
}

func (v *ASTBuilder) VisitCalc(ctx parser.ICalcContext) interface{} {
	if ctx.Unary() != nil {
		return v.Visit(ctx.Unary())
	}

	// Binary arithmetic operations
	if ctx.GetChildCount() == 3 {
		left := v.Visit(ctx.GetChild(0).(antlr.ParseTree)).(ast.ASTNode)
		right := v.Visit(ctx.GetChild(2).(antlr.ParseTree)).(ast.ASTNode)

		op := ctx.GetOp().GetText()
		var functionName string
		if op == "+" {
			functionName = ast.Add
		} else if op == "-" {
			functionName = ast.Subtract
		} else if op == "*" {
			functionName = ast.Multiply
		} else if op == "/" {
			functionName = ast.Divide
		} else if op == "%" {
			functionName = ast.Modulo
		}

		return ast.NewFunctionCall(ast.NewIdent(functionName, false), []ast.ASTNode{left, right})
	}

	return nil
}

func (v *ASTBuilder) VisitUnary(ctx parser.IUnaryContext) interface{} {
	// Check specific Context type
	switch unaryCtx := ctx.(type) {
	case *parser.MemberExprContext:
		if member := unaryCtx.Member(); member != nil {
			return v.Visit(member)
		}
	case *parser.LogicalNotContext:
		if member := unaryCtx.Member(); member != nil {
			memberNode := v.Visit(member).(ast.ASTNode)
			// Handle logical NOT operator
			ops := unaryCtx.GetOps()
			result := memberNode
			for range ops {
				result = ast.NewFunctionCall(ast.NewIdent(ast.LogicalNot, false), []ast.ASTNode{result})
			}
			return result
		}
	case *parser.NegateContext:
		if member := unaryCtx.Member(); member != nil {
			memberNode := v.Visit(member).(ast.ASTNode)
			// Handle negation operator
			ops := unaryCtx.GetOps()
			result := memberNode
			for range ops {
				result = ast.NewFunctionCall(ast.NewIdent(ast.Negate, false), []ast.ASTNode{result})
			}
			return result
		}
	}
	return nil
}

func (v *ASTBuilder) VisitMember(ctx parser.IMemberContext) interface{} {
	// Check specific Context type
	switch memberCtx := ctx.(type) {
	case *parser.PrimaryExprContext:
		if primary := memberCtx.Primary(); primary != nil {
			return v.Visit(primary)
		}
	case *parser.SelectContext:
		if member := memberCtx.Member(); member != nil {
			memberNode := v.Visit(member).(ast.ASTNode)
			fieldName := ""
			// Get field name
			if id := memberCtx.GetId(); id != nil {
				if escapedCtx, ok := id.(*parser.SimpleIdentifierContext); ok {
					fieldName = escapedCtx.GetId().GetText()
				} else if escapedCtx, ok := id.(*parser.EscapedIdentifierContext); ok {
					fieldName = escapedCtx.GetId().GetText()
				}
			}
			optional := memberCtx.GetOpt() != nil
			return ast.NewMemberAccess(memberNode, fieldName, optional)
		}
	case *parser.MemberCallContext:
		if member := memberCtx.Member(); member != nil {
			memberNode := v.Visit(member).(ast.ASTNode)
			methodName := memberCtx.GetId().GetText()
			var args []ast.ASTNode
			if argsCtx := memberCtx.GetArgs(); argsCtx != nil {
				args = v.VisitExprList(argsCtx)
			}
			return ast.NewFunctionCall(ast.NewMemberAccess(memberNode, methodName, false), args)
		}
	case *parser.IndexContext:
		if member := memberCtx.Member(); member != nil {
			memberNode := v.Visit(member).(ast.ASTNode)
			indexNode := v.Visit(memberCtx.GetIndex()).(ast.ASTNode)
			optional := memberCtx.GetOpt() != nil
			return ast.NewIndex(memberNode, indexNode, optional)
		}
	}
	return nil
}

func (v *ASTBuilder) VisitPrimary(ctx parser.IPrimaryContext) interface{} {
	// Check specific Context type
	switch primaryCtx := ctx.(type) {
	case *parser.ConstantLiteralContext:
		if literal := primaryCtx.Literal(); literal != nil {
			return v.Visit(literal)
		}
	case *parser.NestedContext:
		if expr := primaryCtx.Expr(); expr != nil {
			return v.Visit(expr)
		}
	case *parser.IdentContext:
		name := primaryCtx.GetId().GetText()
		receiverStyle := primaryCtx.GetLeadingDot() != nil
		return ast.NewIdent(name, receiverStyle)
	case *parser.GlobalCallContext:
		funcName := primaryCtx.GetId().GetText()
		var args []ast.ASTNode
		if argsCtx := primaryCtx.GetArgs(); argsCtx != nil {
			args = v.VisitExprList(argsCtx)
		}
		receiverStyle := primaryCtx.GetLeadingDot() != nil
		if receiverStyle {
			return ast.NewFunctionCall(ast.NewIdent("."+funcName, false), args)
		}
		return ast.NewFunctionCall(ast.NewIdent(funcName, false), args)
	case *parser.CreateListContext:
		var elements []ast.ASTNode
		if elemsCtx := primaryCtx.GetElems(); elemsCtx != nil {
			elements = v.VisitListInit(elemsCtx)
		}
		return ast.NewList(elements)
	case *parser.CreateStructContext:
		var entries []ast.MapEntry
		if entriesCtx := primaryCtx.GetEntries(); entriesCtx != nil {
			entries = v.VisitMapInitializerList(entriesCtx)
		}
		return ast.NewMap(entries)
	case *parser.CreateMessageContext:
		// Build type name
		var typeName string
		ids := primaryCtx.GetIds()
		if len(ids) > 0 {
			parts := make([]string, len(ids))
			for i, id := range ids {
				parts[i] = id.GetText()
			}
			typeName = strings.Join(parts, ".")
		}
		receiverStyle := primaryCtx.GetLeadingDot() != nil

		// Handle field initialization list
		var fields []ast.StructField
		if entriesCtx := primaryCtx.GetEntries(); entriesCtx != nil {
			fields = v.VisitFieldInitializerList(entriesCtx.(*parser.FieldInitializerListContext))
		}
		return ast.NewStruct(typeName, fields, receiverStyle)
	}
	return nil
}

func (v *ASTBuilder) VisitLiteral(ctx parser.ILiteralContext) interface{} {
	// Check specific child Context type
	switch literalCtx := ctx.(type) {
	case *parser.IntContext:
		text := literalCtx.GetTok().GetText()
		base := 10
		if strings.HasPrefix(text, "0x") {
			base = 16
			text = text[2:]
		}
		sign := literalCtx.GetSign()
		if sign != nil {
			text = sign.GetText() + text
		}

		val, err := strconv.ParseInt(text, base, 64)
		if err != nil {
			return err
		}
		return ast.NewLiteral(ast.NewIntValue(val))
	case *parser.UintContext:
		text := literalCtx.GetTok().GetText()
		if strings.HasSuffix(text, "u") || strings.HasSuffix(text, "U") {
			text = text[:len(text)-1]
		}
		base := 10
		if strings.HasPrefix(text, "0x") {
			base = 16
			text = text[2:]
		}
		val, err := strconv.ParseUint(text, base, 64)
		if err != nil {
			return err
		}
		return ast.NewLiteral(ast.NewUintValue(val))
	case *parser.DoubleContext:
		text := literalCtx.GetTok().GetText()
		sign := literalCtx.GetSign()
		if sign != nil {
			text = sign.GetText() + text
		}
		val, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return err
		}
		return ast.NewLiteral(ast.NewDoubleValue(val))
	case *parser.StringContext:
		text := literalCtx.GetTok().GetText()
		val, err := unescape(text, false)
		if err != nil {
			return err
		}
		return ast.NewLiteral(ast.NewStringValue(val))
	case *parser.BytesContext:
		text := literalCtx.GetTok().GetText()
		if strings.HasPrefix(text, "b") || strings.HasPrefix(text, "B") {
			text = text[1:]
		}
		val, err := unescape(text, true)
		if err != nil {
			return err
		}
		return ast.NewLiteral(ast.NewBytesValue([]byte(val)))
	case *parser.BoolTrueContext:
		return ast.NewLiteral(ast.NewBoolValue(true))
	case *parser.BoolFalseContext:
		return ast.NewLiteral(ast.NewBoolValue(false))
	case *parser.NullContext:
		return ast.NewLiteral(ast.NewNullValue())
	}
	return nil
}

func (v *ASTBuilder) VisitExprList(ctx parser.IExprListContext) []ast.ASTNode {
	if ctx == nil {
		return nil
	}

	exprListCtx := ctx.(*parser.ExprListContext)
	var args []ast.ASTNode
	for _, e := range exprListCtx.GetE() {
		args = append(args, v.Visit(e).(ast.ASTNode))
	}
	return args
}

func (v *ASTBuilder) VisitListInit(ctx parser.IListInitContext) []ast.ASTNode {
	if ctx == nil {
		return nil
	}

	listInitCtx := ctx.(*parser.ListInitContext)
	var elements []ast.ASTNode
	for _, elem := range listInitCtx.GetElems() {
		// Handle optional expressions
		if optExpr, ok := elem.(*parser.OptExprContext); ok {
			elements = append(elements, v.Visit(optExpr.GetE()).(ast.ASTNode))
		}
	}
	return elements
}

func (v *ASTBuilder) VisitMapInitializerList(ctx parser.IMapInitializerListContext) []ast.MapEntry {
	if ctx == nil {
		return nil
	}

	mapInitCtx := ctx.(*parser.MapInitializerListContext)
	var entries []ast.MapEntry
	for i, key := range mapInitCtx.GetKeys() {
		if i < len(mapInitCtx.GetValues()) {
			keyNode := v.visitOptExpr(key)
			valueNode := v.Visit(mapInitCtx.GetValues()[i]).(ast.ASTNode)
			entries = append(entries, ast.NewMapEntry(keyNode, valueNode, false))
		}
	}
	return entries
}

func (v *ASTBuilder) visitOptExpr(ctx parser.IOptExprContext) ast.ASTNode {
	if optExpr, ok := ctx.(*parser.OptExprContext); ok {
		return v.Visit(optExpr.GetE()).(ast.ASTNode)
	}
	return nil
}

// VisitFieldInitializerList handles field initialization list
func (v *ASTBuilder) VisitFieldInitializerList(ctx *parser.FieldInitializerListContext) []ast.StructField {
	if ctx == nil {
		return nil
	}

	var fields []ast.StructField
	for i, fieldCtx := range ctx.GetFields() {
		if i < len(ctx.GetValues()) {
			fieldName := ""
			optional := false

			// Get field name and whether it's optional
			if optFieldCtx, ok := fieldCtx.(*parser.OptFieldContext); ok {
				optional = optFieldCtx.GetOpt() != nil
				if escapedCtx := optFieldCtx.EscapeIdent(); escapedCtx != nil {
					if simpleCtx, ok := escapedCtx.(*parser.SimpleIdentifierContext); ok {
						fieldName = simpleCtx.GetId().GetText()
					} else if escapedCtx, ok := escapedCtx.(*parser.EscapedIdentifierContext); ok {
						fieldName = escapedCtx.GetId().GetText()
					}
				}
			}

			valueNode := v.Visit(ctx.GetValues()[i]).(ast.ASTNode)
			fields = append(fields, ast.StructField{
				Name:     fieldName,
				Value:    valueNode,
				Optional: optional,
			})
		}
	}
	return fields
}

// copy from cel-go/parser/unescape.go

// Unescape takes a quoted string, unquotes, and unescapes it.
//
// This function performs escaping compatible with GoogleSQL.
func unescape(value string, isBytes bool) (string, error) {
	// All strings normalize newlines to the \n representation.
	value = newlineNormalizer.Replace(value)
	n := len(value)

	// Nothing to unescape / decode.
	if n < 2 {
		return value, errors.New("unable to unescape string")
	}

	// Raw string preceded by the 'r|R' prefix.
	isRawLiteral := false
	if value[0] == 'r' || value[0] == 'R' {
		value = value[1:]
		n = len(value)
		isRawLiteral = true
	}

	// Quoted string of some form, must have same first and last char.
	if value[0] != value[n-1] || (value[0] != '"' && value[0] != '\'') {
		return value, errors.New("unable to unescape string")
	}

	// Normalize the multi-line CEL string representation to a standard
	// Go quoted string.
	if n >= 6 {
		if strings.HasPrefix(value, "'''") {
			if !strings.HasSuffix(value, "'''") {
				return value, errors.New("unable to unescape string")
			}
			value = "\"" + value[3:n-3] + "\""
		} else if strings.HasPrefix(value, `"""`) {
			if !strings.HasSuffix(value, `"""`) {
				return value, errors.New("unable to unescape string")
			}
			value = "\"" + value[3:n-3] + "\""
		}
		n = len(value)
	}
	value = value[1 : n-1]
	// If there is nothing to escape, then return.
	if isRawLiteral || !strings.ContainsRune(value, '\\') {
		return value, nil
	}

	// Otherwise the string contains escape characters.
	// The following logic is adapted from `strconv/quote.go`
	var runeTmp [utf8.UTFMax]byte
	buf := make([]byte, 0, 3*n/2)
	for len(value) > 0 {
		c, encode, rest, err := unescapeChar(value, isBytes)
		if err != nil {
			return "", err
		}
		value = rest
		if c < utf8.RuneSelf || !encode {
			buf = append(buf, byte(c))
		} else {
			n := utf8.EncodeRune(runeTmp[:], c)
			buf = append(buf, runeTmp[:n]...)
		}
	}
	return string(buf), nil
}

// unescapeChar takes a string input and returns the following info:
//
//	value - the escaped unicode rune at the front of the string.
//	encode - the value should be unicode-encoded
//	tail - the remainder of the input string.
//	err - error value, if the character could not be unescaped.
//
// When encode is true the return value may still fit within a single byte,
// but unicode encoding is attempted which is more expensive than when the
// value is known to self-represent as a single byte.
//
// If isBytes is set, unescape as a bytes literal so octal and hex escapes
// represent byte values, not unicode code points.
func unescapeChar(s string, isBytes bool) (value rune, encode bool, tail string, err error) {
	// 1. Character is not an escape sequence.
	switch c := s[0]; {
	case c >= utf8.RuneSelf:
		r, size := utf8.DecodeRuneInString(s)
		return r, true, s[size:], nil
	case c != '\\':
		return rune(s[0]), false, s[1:], nil
	}

	// 2. Last character is the start of an escape sequence.
	if len(s) <= 1 {
		err = errors.New("unable to unescape string, found '\\' as last character")
		return
	}

	c := s[1]
	s = s[2:]
	// 3. Common escape sequences shared with Google SQL
	switch c {
	case 'a':
		value = '\a'
	case 'b':
		value = '\b'
	case 'f':
		value = '\f'
	case 'n':
		value = '\n'
	case 'r':
		value = '\r'
	case 't':
		value = '\t'
	case 'v':
		value = '\v'
	case '\\':
		value = '\\'
	case '\'':
		value = '\''
	case '"':
		value = '"'
	case '`':
		value = '`'
	case '?':
		value = '?'

	// 4. Unicode escape sequences, reproduced from `strconv/quote.go`
	case 'x', 'X', 'u', 'U':
		n := 0
		encode = true
		switch c {
		case 'x', 'X':
			n = 2
			encode = !isBytes
		case 'u':
			n = 4
			if isBytes {
				err = errors.New("unable to unescape string")
				return
			}
		case 'U':
			n = 8
			if isBytes {
				err = errors.New("unable to unescape string")
				return
			}
		}
		var v rune
		if len(s) < n {
			err = errors.New("unable to unescape string")
			return
		}
		for j := 0; j < n; j++ {
			x, ok := unhex(s[j])
			if !ok {
				err = errors.New("unable to unescape string")
				return
			}
			v = v<<4 | x
		}
		s = s[n:]
		if !isBytes && !utf8.ValidRune(v) {
			err = errors.New("invalid unicode code point")
			return
		}
		value = v

	// 5. Octal escape sequences, must be three digits \[0-3][0-7][0-7]
	case '0', '1', '2', '3':
		if len(s) < 2 {
			err = errors.New("unable to unescape octal sequence in string")
			return
		}
		v := rune(c - '0')
		for j := 0; j < 2; j++ {
			x := s[j]
			if x < '0' || x > '7' {
				err = errors.New("unable to unescape octal sequence in string")
				return
			}
			v = v*8 + rune(x-'0')
		}
		if !isBytes && !utf8.ValidRune(v) {
			err = errors.New("invalid unicode code point")
			return
		}
		value = v
		s = s[2:]
		encode = !isBytes

		// Unknown escape sequence.
	default:
		err = errors.New("unable to unescape string")
	}

	tail = s
	return
}

func unhex(b byte) (rune, bool) {
	c := rune(b)
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

var (
	newlineNormalizer = strings.NewReplacer("\r\n", "\n", "\r", "\n")
)
