package ast

import (
	"fmt"
	"strings"
)

// ASTNode 表示抽象语法树的节点
type ASTNode interface {
	String() string
	Accept(visitor ASTVisitor) (interface{}, error)
}

// ASTVisitor 访问者模式接口
type ASTVisitor interface {
	VisitLiteral(node *LiteralNode) (interface{}, error)
	VisitIdent(node *IdentNode) (interface{}, error)
	VisitMemberAccess(node *MemberAccessNode) (interface{}, error)
	VisitFunctionCall(node *FunctionCallNode) (interface{}, error)
	VisitIndex(node *IndexNode) (interface{}, error)
	VisitConditional(node *ConditionalNode) (interface{}, error)
	VisitList(node *ListNode) (interface{}, error)
	VisitMap(node *MapNode) (interface{}, error)
	VisitStruct(node *StructNode) (interface{}, error)
}

// LiteralNode 字面量节点
type LiteralNode struct {
	Value Value
}

func (n *LiteralNode) String() string {
	return n.Value.String()
}

func (n *LiteralNode) Accept(visitor ASTVisitor) (interface{}, error) {
	return visitor.VisitLiteral(n)
}

// IdentNode 标识符节点
type IdentNode struct {
	Name       string
	LeadingDot bool // 是否有前导点（.identifier）
}

func (n *IdentNode) String() string {
	if n.LeadingDot {
		return "." + n.Name
	}
	return n.Name
}

func (n *IdentNode) Accept(visitor ASTVisitor) (interface{}, error) {
	return visitor.VisitIdent(n)
}

// MemberAccessNode 成员访问节点 (obj.member)
type MemberAccessNode struct {
	Object   ASTNode
	Member   string
	Optional bool // 是否是可选访问 (obj.?member)
}

func (n *MemberAccessNode) String() string {
	op := "."
	if n.Optional {
		op = ".?"
	}
	return fmt.Sprintf("(%s%s%s)", n.Object.String(), op, n.Member)
}

func (n *MemberAccessNode) Accept(visitor ASTVisitor) (interface{}, error) {
	return visitor.VisitMemberAccess(n)
}

// FunctionCallNode 函数调用节点
type FunctionCallNode struct {
	Function ASTNode   // 函数表达式，可以是标识符或成员访问
	Args     []ASTNode // 参数列表
}

func (n *FunctionCallNode) String() string {
	args := make([]string, len(n.Args))
	for i, arg := range n.Args {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", n.Function.String(), strings.Join(args, ", "))
}

func (n *FunctionCallNode) Accept(visitor ASTVisitor) (interface{}, error) {
	return visitor.VisitFunctionCall(n)
}

// IndexNode 索引访问节点 (obj[index])
type IndexNode struct {
	Object   ASTNode
	Index    ASTNode
	Optional bool // 是否是可选索引 (obj[?index])
}

func (n *IndexNode) String() string {
	op := "["
	if n.Optional {
		op = "[?"
	}
	return fmt.Sprintf("(%s%s%s])", n.Object.String(), op, n.Index.String())
}

func (n *IndexNode) Accept(visitor ASTVisitor) (interface{}, error) {
	return visitor.VisitIndex(n)
}

// ConditionalNode 条件表达式节点 (condition ? trueExpr : falseExpr)
type ConditionalNode struct {
	Condition ASTNode
	TrueExpr  ASTNode
	FalseExpr ASTNode
}

func (n *ConditionalNode) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", n.Condition.String(), n.TrueExpr.String(), n.FalseExpr.String())
}

func (n *ConditionalNode) Accept(visitor ASTVisitor) (interface{}, error) {
	return visitor.VisitConditional(n)
}

// ListNode 列表字面量节点
type ListNode struct {
	Elements []ASTNode
}

func (n *ListNode) String() string {
	elements := make([]string, len(n.Elements))
	for i, elem := range n.Elements {
		elements[i] = elem.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (n *ListNode) Accept(visitor ASTVisitor) (interface{}, error) {
	return visitor.VisitList(n)
}

// MapNode 映射字面量节点
type MapNode struct {
	Entries []MapEntry
}

func (n *MapNode) String() string {
	entries := make([]string, len(n.Entries))
	for i, entry := range n.Entries {
		entries[i] = fmt.Sprintf("%s: %s", entry.Key.String(), entry.Value.String())
	}
	return fmt.Sprintf("{%s}", strings.Join(entries, ", "))
}

func (n *MapNode) Accept(visitor ASTVisitor) (interface{}, error) {
	return visitor.VisitMap(n)
}

// MapEntry 映射条目
type MapEntry struct {
	Key      ASTNode
	Value    ASTNode
	Optional bool // 是否是可选键 (?key: value)
}

// StructNode 结构体字面量节点 (Message{field: value})
type StructNode struct {
	TypeName      string
	Fields        []StructField
	ReceiverStyle bool // 是否有前导点 (.Message{})
}

func (n *StructNode) String() string {
	fields := make([]string, len(n.Fields))
	for i, field := range n.Fields {
		opt := ""
		if field.Optional {
			opt = "?"
		}
		fields[i] = fmt.Sprintf("%s%s: %s", opt, field.Name, field.Value.String())
	}
	typeName := n.TypeName
	if n.ReceiverStyle {
		typeName = "." + typeName
	}
	return fmt.Sprintf("%s{%s}", typeName, strings.Join(fields, ", "))
}

func (n *StructNode) Accept(visitor ASTVisitor) (interface{}, error) {
	return visitor.VisitStruct(n)
}

// StructField 表示结构体字段
type StructField struct {
	Name     string
	Value    ASTNode
	Optional bool
}

// 便利函数用于创建AST节点
func NewLiteral(value Value) *LiteralNode {
	return &LiteralNode{Value: value}
}

func NewIdent(name string, leadingDot bool) *IdentNode {
	return &IdentNode{Name: name, LeadingDot: leadingDot}
}

func NewMemberAccess(object ASTNode, member string, optional bool) *MemberAccessNode {
	return &MemberAccessNode{Object: object, Member: member, Optional: optional}
}

func NewFunctionCall(function ASTNode, args []ASTNode) *FunctionCallNode {
	return &FunctionCallNode{Function: function, Args: args}
}

func NewIndex(object ASTNode, index ASTNode, optional bool) *IndexNode {
	return &IndexNode{Object: object, Index: index, Optional: optional}
}

func NewConditional(condition, trueExpr, falseExpr ASTNode) *ConditionalNode {
	return &ConditionalNode{Condition: condition, TrueExpr: trueExpr, FalseExpr: falseExpr}
}

func NewList(elements []ASTNode) *ListNode {
	return &ListNode{Elements: elements}
}

func NewMap(entries []MapEntry) *MapNode {
	return &MapNode{Entries: entries}
}

func NewMapEntry(key, value ASTNode, optional bool) MapEntry {
	return MapEntry{Key: key, Value: value, Optional: optional}
}

func NewStruct(typeName string, fields []StructField, receiverStyle bool) *StructNode {
	return &StructNode{
		TypeName:      typeName,
		Fields:        fields,
		ReceiverStyle: receiverStyle,
	}
}
