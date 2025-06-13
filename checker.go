package sl

import (
	"fmt"

	"github.com/yywing/sl/ast"
)

// CheckError 表示类型检查错误
type CheckError struct {
	Message string
	Node    ast.ASTNode
}

func (e *CheckError) Error() string {
	return fmt.Sprintf("type check error: %s", e.Message)
}

// Checker 实现类型检查
type Checker struct {
	env *Env
}

// NewChecker 创建新的类型检查器
func NewChecker(env *Env) *Checker {
	return &Checker{env: env}
}

// Check 检查表达式的类型
func (tc *Checker) Check(node ast.ASTNode) (ast.ValueType, error) {
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
	if value, exists := tc.env.GetVariable(node.Name); exists {
		return value.Type(), nil
	}

	return nil, &CheckError{
		Message: fmt.Sprintf("undefined identifier: %s", node.Name),
		Node:    node,
	}
}

func (tc *Checker) VisitMemberAccess(node *ast.MemberAccessNode) (interface{}, error) {
	objectType, err := tc.Check(node.Object)
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

func (tc *Checker) resolveDynamicType(env map[string]ast.ValueType, dyn ast.ValueType, rt ast.ValueType) (ast.ValueType, error) {
	if !dyn.IsDyn() {
		return dyn, nil
	}

	switch dyn := dyn.(type) {
	case *ast.ListType:
		var innerType ast.ValueType
		if rt != nil {
			e, ok := rt.(*ast.ListType)
			if !ok {
				return nil, fmt.Errorf("cannot resolve dynamic type %s to %s", dyn.String(), rt.String())
			}
			innerType = e.ElementType()
		}
		result, err := tc.resolveDynamicType(env, dyn.ElementType(), innerType)
		if err != nil {
			return nil, err
		}
		return ast.NewListType(result), nil
	case *ast.MapType:
		var keyType ast.ValueType
		var valueType ast.ValueType
		if rt != nil {
			m, ok := rt.(*ast.MapType)
			if !ok {
				return nil, fmt.Errorf("cannot resolve dynamic type %s to %s", dyn.String(), rt.String())
			}
			keyType = m.KeyType()
			valueType = m.ValueType()
		}
		key, err := tc.resolveDynamicType(env, dyn.KeyType(), keyType)
		if err != nil {
			return nil, err
		}
		value, err := tc.resolveDynamicType(env, dyn.ValueType(), valueType)
		if err != nil {
			return nil, err
		}
		return ast.NewMapType(key, value), nil
	default:
		t, ok := env[dyn.Kind()]
		if rt == nil {
			if !ok {
				return nil, fmt.Errorf("dynamic type %s not found", dyn.String())
			}
			return t, nil
		}

		if ok {
			if !t.Equals(rt) {
				return nil, fmt.Errorf("dynamic type should %s but got %s", t.String(), rt.String())
			}
			env[dyn.Kind()] = tc.getDeterministicType(t, rt)
		} else {
			env[dyn.Kind()] = rt
		}

		return env[dyn.Kind()], nil
	}
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

	// 获取参数类型
	argTypes := make([]ast.ValueType, len(args))
	for i, arg := range args {
		argType, err := tc.Check(arg)
		if err != nil {
			return nil, err
		}
		argTypes[i] = argType
	}

	var err error
	// 检查参数类型
	fnTypes := f.Types()
	resultEnv := make(map[string]ast.ValueType)
	var foundFnType *ast.FunctionType
Outer:
	for _, fnType := range fnTypes {
		if len(fnType.ParamTypes()) != len(args) {
			continue
		}

		env := make(map[string]ast.ValueType)
		for i, argType := range argTypes {
			paramType := fnType.ParamTypes()[i]
			if paramType.IsDyn() {
				paramType, err = tc.resolveDynamicType(env, paramType, argType)
				if err != nil {
					continue Outer
				}
			}
			if !tc.isCompatible(argType, paramType) {
				continue Outer
			}
		}
		foundFnType = &fnType
		resultEnv = env
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
		resultType, err = tc.resolveDynamicType(resultEnv, resultType, nil)
		if err != nil {
			return nil, err
		}
	}

	return resultType, nil
}

func (tc *Checker) VisitIndex(node *ast.IndexNode) (interface{}, error) {
	objectType, err := tc.Check(node.Object)
	if err != nil {
		return nil, err
	}

	indexType, err := tc.Check(node.Index)
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
	conditionType, err := tc.Check(node.Condition)
	if err != nil {
		return nil, err
	}

	if conditionType.Kind() != ast.TypeKindBool {
		return nil, &CheckError{
			Message: fmt.Sprintf("conditional expression requires bool condition, got %s", conditionType.String()),
			Node:    node,
		}
	}

	trueType, err := tc.Check(node.TrueExpr)
	if err != nil {
		return nil, err
	}

	falseType, err := tc.Check(node.FalseExpr)
	if err != nil {
		return nil, err
	}

	// 检查两个分支的类型是否兼容
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

	// 检查第一个元素的类型作为列表元素类型
	firstElemType, err := tc.Check(node.Elements[0])
	if err != nil {
		return nil, err
	}

	// 检查所有元素类型是否一致
	for _, elem := range node.Elements[1:] {
		elemType, err := tc.Check(elem)
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
		// 空映射
		return ast.NewMapType(ast.AnyType, ast.AnyType), nil
	}

	// 检查第一个条目的类型
	firstEntry := node.Entries[0]
	keyType, err := tc.Check(firstEntry.Key)
	if err != nil {
		return nil, err
	}

	valueType, err := tc.Check(firstEntry.Value)
	if err != nil {
		return nil, err
	}

	// 检查所有条目的类型一致性
	for i, entry := range node.Entries[1:] {
		entryKeyType, err := tc.Check(entry.Key)
		if err != nil {
			return nil, err
		}

		entryValueType, err := tc.Check(entry.Value)
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

func (tc *Checker) getDeterministicType(ts ...ast.ValueType) ast.ValueType {
	for _, t := range ts {
		if t.Kind() != ast.TypeKindAny {
			return t
		}
	}
	return ast.AnyType
}

func (tc *Checker) VisitStruct(node *ast.StructNode) (interface{}, error) {
	// TODO: 暂时不支持 struct 语法
	return nil, &CheckError{
		Message: "struct is not supported",
		Node:    node,
	}
}
