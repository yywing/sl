package ast

import (
	"fmt"
	"strings"
)

// ASTNode represents an abstract syntax tree node
type ASTNode interface {
	String() string
	Accept(visitor ASTVisitor) (interface{}, error)
}

// ASTVisitor visitor pattern interface
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

// LiteralNode literal value node
type LiteralNode struct {
	Value Value
}

func (n *LiteralNode) String() string {
	return n.Value.String()
}

func (n *LiteralNode) Accept(visitor ASTVisitor) (interface{}, error) {
	return visitor.VisitLiteral(n)
}

// IdentNode identifier node
type IdentNode struct {
	Name       string
	LeadingDot bool // whether there is a leading dot (.identifier)
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

// MemberAccessNode member access node (obj.member)
type MemberAccessNode struct {
	Object   ASTNode
	Member   string
	Optional bool // whether it's optional access (obj.?member)
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

// FunctionCallNode function call node
type FunctionCallNode struct {
	Function ASTNode   // function expression, can be identifier or member access
	Args     []ASTNode // argument list
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

// IndexNode index access node (obj[index])
type IndexNode struct {
	Object   ASTNode
	Index    ASTNode
	Optional bool // whether it's optional index (obj[?index])
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

// ConditionalNode conditional expression node (condition ? trueExpr : falseExpr)
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

// ListNode list literal node
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

// MapNode map literal node
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

// MapEntry map entry
type MapEntry struct {
	Key      ASTNode
	Value    ASTNode
	Optional bool // whether it's an optional key (?key: value)
}

// StructNode struct literal node (Message{field: value})
type StructNode struct {
	TypeName      string
	Fields        []StructField
	ReceiverStyle bool // whether there is a leading dot (.Message{})
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

// StructField represents a struct field
type StructField struct {
	Name     string
	Value    ASTNode
	Optional bool
}

// Convenience functions for creating AST nodes
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
