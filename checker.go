package sl

import (
	"fmt"

	"github.com/yywing/sl/ast"
)

// CheckError represents type checking error
type CheckError struct {
	Message string
	Node    ast.ASTNode
}

func (e *CheckError) Error() string {
	return fmt.Sprintf("type check error: %s", e.Message)
}

// Checker implements type checking
type Checker struct {
	env     *Env
	program *Program
}

// NewChecker creates a new type checker
func NewChecker(env *Env, program *Program) *Checker {
	return &Checker{env: env, program: program}
}

// Check checks the type of expression
func (tc *Checker) Check() (ast.ValueType, error) {
	return tc.check(tc.program.ASTNode)
}

func (tc *Checker) check(node ast.ASTNode) (ast.ValueType, error) {
	result, err := node.Accept(tc)
	if err != nil {
		return nil, err
	}
	if typ, ok := result.(ast.ValueType); ok {
		return typ, nil
	}
	return nil, fmt.Errorf("internal error: type checker returned non-type")
}

func (tc *Checker) VisitLiteral(node *ast.LiteralNode) (interface{}, error) {
	return node.Value.Type(), nil
}

func (tc *Checker) VisitIdent(node *ast.IdentNode) (interface{}, error) {
	if t, exists := tc.program.GetVariable(node.Name); exists {
		return t, nil
	}

	return nil, &CheckError{
		Message: fmt.Sprintf("undefined identifier: %s", node.Name),
		Node:    node,
	}
}

func (tc *Checker) VisitMemberAccess(node *ast.MemberAccessNode) (interface{}, error) {
	objectType, err := tc.check(node.Object)
	if err != nil {
		return nil, err
	}

	if !objectType.HasTrait(ast.SelectorType) {
		return nil, &CheckError{
			Message: fmt.Sprintf("cannot access member of type %s", objectType.String()),
			Node:    node,
		}
	}

	memberType := objectType.Member(node.Member)
	if memberType == nil {
		return nil, &CheckError{
			Message: fmt.Sprintf("member %s not found in type %s", node.Member, objectType.String()),
			Node:    node,
		}
	}

	return memberType, nil
}

func (tc *Checker) VisitFunctionCall(node *ast.FunctionCallNode) (interface{}, error) {
	var fnName string
	var args []ast.ASTNode
	switch fn := node.Function.(type) {
	case *ast.IdentNode:
		fnName = fn.Name
		args = node.Args
	case *ast.MemberAccessNode:
		fnName = fn.Member
		args = append([]ast.ASTNode{fn.Object}, node.Args...)
	default:
		return nil, &CheckError{
			Message: fmt.Sprintf("function call must be an identifier or member access, got %s", node.Function.String()),
			Node:    node,
		}
	}

	f, ok := tc.env.GetFunction(fnName)
	if !ok {
		return nil, &CheckError{
			Message: fmt.Sprintf("function %s not found", fnName),
			Node:    node,
		}
	}

	// Get argument types
	argTypes := make([]ast.ValueType, len(args))
	for i, arg := range args {
		argType, err := tc.check(arg)
		if err != nil {
			return nil, err
		}
		argTypes[i] = argType
	}

	// Check argument types
	fnTypes := f.Types()
	resultEnv := make(map[string]ast.ValueType)
	var foundFnType *ast.FunctionType
	for _, fnType := range fnTypes {
		var ok bool
		resultEnv, ok = ast.MatchFunctionTypes(fnType.ParamTypes(), argTypes)
		if !ok {
			continue
		}
		foundFnType = &fnType
		break
	}

	if foundFnType == nil {
		return nil, &CheckError{
			Message: fmt.Sprintf("function %s not found with args %v", f.Name(), argTypes),
			Node:    node,
		}
	}

	resultType := foundFnType.ReturnType()
	if resultType.IsDyn() {
		var err error
		resultType, err = ast.ResolveDynamicType(resultEnv, resultType, nil)
		if err != nil {
			return nil, err
		}
	}

	return resultType, nil
}

func (tc *Checker) VisitIndex(node *ast.IndexNode) (interface{}, error) {
	objectType, err := tc.check(node.Object)
	if err != nil {
		return nil, err
	}

	indexType, err := tc.check(node.Index)
	if err != nil {
		return nil, err
	}

	switch objType := objectType.(type) {
	case *ast.ListType:
		if indexType.Kind() != ast.TypeKindInt && indexType.Kind() != ast.TypeKindUint {
			return nil, &CheckError{
				Message: fmt.Sprintf("list index must be integer, got %s", indexType.String()),
				Node:    node,
			}
		}
		return objType.ElementType(), nil
	case *ast.MapType:
		if !tc.isCompatible(indexType, objType.KeyType()) {
			return nil, &CheckError{
				Message: fmt.Sprintf("map key type mismatch: expected %s, got %s", objType.KeyType().String(), indexType.String()),
				Node:    node,
			}
		}
		return objType.ValueType(), nil
	default:
		return nil, &CheckError{
			Message: fmt.Sprintf("cannot index type %s", objectType.String()),
			Node:    node,
		}
	}
}

func (tc *Checker) VisitConditional(node *ast.ConditionalNode) (interface{}, error) {
	conditionType, err := tc.check(node.Condition)
	if err != nil {
		return nil, err
	}

	if conditionType.Kind() != ast.TypeKindBool {
		return nil, &CheckError{
			Message: fmt.Sprintf("conditional expression requires bool condition, got %s", conditionType.String()),
			Node:    node,
		}
	}

	trueType, err := tc.check(node.TrueExpr)
	if err != nil {
		return nil, err
	}

	falseType, err := tc.check(node.FalseExpr)
	if err != nil {
		return nil, err
	}

	// Check if the types of both branches are compatible
	if tc.isCompatible(trueType, falseType) {
		return trueType, nil
	} else if tc.isCompatible(falseType, trueType) {
		return falseType, nil
	} else {
		return nil, &CheckError{
			Message: fmt.Sprintf("conditional branches have incompatible types: %s and %s", trueType.String(), falseType.String()),
			Node:    node,
		}
	}
}

func (tc *Checker) VisitList(node *ast.ListNode) (interface{}, error) {
	if len(node.Elements) == 0 {
		return ast.NewListType(ast.AnyType), nil
	}

	// Check the type of the first element as the list element type
	firstElemType, err := tc.check(node.Elements[0])
	if err != nil {
		return nil, err
	}

	// Check if all element types are consistent
	for _, elem := range node.Elements[1:] {
		elemType, err := tc.check(elem)
		if err != nil {
			return nil, err
		}

		if !tc.isCompatible(elemType, firstElemType) {
			firstElemType = ast.AnyType
		}
	}

	return ast.NewListType(firstElemType), nil
}

func (tc *Checker) VisitMap(node *ast.MapNode) (interface{}, error) {
	if len(node.Entries) == 0 {
		// Empty map
		return ast.NewMapType(ast.AnyType, ast.AnyType), nil
	}

	// Check the type of the first entry
	firstEntry := node.Entries[0]
	keyType, err := tc.check(firstEntry.Key)
	if err != nil {
		return nil, err
	}

	valueType, err := tc.check(firstEntry.Value)
	if err != nil {
		return nil, err
	}

	// Check type consistency of all entries
	for i, entry := range node.Entries[1:] {
		entryKeyType, err := tc.check(entry.Key)
		if err != nil {
			return nil, err
		}

		entryValueType, err := tc.check(entry.Value)
		if err != nil {
			return nil, err
		}

		if !tc.isCompatible(entryKeyType, keyType) {
			return nil, &CheckError{
				Message: fmt.Sprintf("map entry %d key has type %s, expected %s", i+1, entryKeyType.String(), keyType.String()),
				Node:    node,
			}
		}

		if !tc.isCompatible(entryValueType, valueType) {
			valueType = ast.AnyType
		}
	}

	return ast.NewMapType(keyType, valueType), nil
}

func (tc *Checker) isCompatible(t1, t2 ast.ValueType) bool {
	return t1.Equals(t2)
}

func (tc *Checker) VisitStruct(node *ast.StructNode) (interface{}, error) {
	// TODO: struct syntax is not supported yet
	return nil, &CheckError{
		Message: "struct is not supported",
		Node:    node,
	}
}
