// Code generated from /root/self/sl/SL.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // SL

import "github.com/antlr4-go/antlr/v4"

// BaseSLListener is a complete listener for a parse tree produced by SLParser.
type BaseSLListener struct{}

var _ SLListener = &BaseSLListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseSLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseSLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseSLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseSLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BaseSLListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseSLListener) ExitStart(ctx *StartContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseSLListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseSLListener) ExitExpr(ctx *ExprContext) {}

// EnterConditionalOr is called when production conditionalOr is entered.
func (s *BaseSLListener) EnterConditionalOr(ctx *ConditionalOrContext) {}

// ExitConditionalOr is called when production conditionalOr is exited.
func (s *BaseSLListener) ExitConditionalOr(ctx *ConditionalOrContext) {}

// EnterConditionalAnd is called when production conditionalAnd is entered.
func (s *BaseSLListener) EnterConditionalAnd(ctx *ConditionalAndContext) {}

// ExitConditionalAnd is called when production conditionalAnd is exited.
func (s *BaseSLListener) ExitConditionalAnd(ctx *ConditionalAndContext) {}

// EnterRelation is called when production relation is entered.
func (s *BaseSLListener) EnterRelation(ctx *RelationContext) {}

// ExitRelation is called when production relation is exited.
func (s *BaseSLListener) ExitRelation(ctx *RelationContext) {}

// EnterCalc is called when production calc is entered.
func (s *BaseSLListener) EnterCalc(ctx *CalcContext) {}

// ExitCalc is called when production calc is exited.
func (s *BaseSLListener) ExitCalc(ctx *CalcContext) {}

// EnterMemberExpr is called when production MemberExpr is entered.
func (s *BaseSLListener) EnterMemberExpr(ctx *MemberExprContext) {}

// ExitMemberExpr is called when production MemberExpr is exited.
func (s *BaseSLListener) ExitMemberExpr(ctx *MemberExprContext) {}

// EnterLogicalNot is called when production LogicalNot is entered.
func (s *BaseSLListener) EnterLogicalNot(ctx *LogicalNotContext) {}

// ExitLogicalNot is called when production LogicalNot is exited.
func (s *BaseSLListener) ExitLogicalNot(ctx *LogicalNotContext) {}

// EnterNegate is called when production Negate is entered.
func (s *BaseSLListener) EnterNegate(ctx *NegateContext) {}

// ExitNegate is called when production Negate is exited.
func (s *BaseSLListener) ExitNegate(ctx *NegateContext) {}

// EnterMemberCall is called when production MemberCall is entered.
func (s *BaseSLListener) EnterMemberCall(ctx *MemberCallContext) {}

// ExitMemberCall is called when production MemberCall is exited.
func (s *BaseSLListener) ExitMemberCall(ctx *MemberCallContext) {}

// EnterSelect is called when production Select is entered.
func (s *BaseSLListener) EnterSelect(ctx *SelectContext) {}

// ExitSelect is called when production Select is exited.
func (s *BaseSLListener) ExitSelect(ctx *SelectContext) {}

// EnterPrimaryExpr is called when production PrimaryExpr is entered.
func (s *BaseSLListener) EnterPrimaryExpr(ctx *PrimaryExprContext) {}

// ExitPrimaryExpr is called when production PrimaryExpr is exited.
func (s *BaseSLListener) ExitPrimaryExpr(ctx *PrimaryExprContext) {}

// EnterIndex is called when production Index is entered.
func (s *BaseSLListener) EnterIndex(ctx *IndexContext) {}

// ExitIndex is called when production Index is exited.
func (s *BaseSLListener) ExitIndex(ctx *IndexContext) {}

// EnterIdent is called when production Ident is entered.
func (s *BaseSLListener) EnterIdent(ctx *IdentContext) {}

// ExitIdent is called when production Ident is exited.
func (s *BaseSLListener) ExitIdent(ctx *IdentContext) {}

// EnterGlobalCall is called when production GlobalCall is entered.
func (s *BaseSLListener) EnterGlobalCall(ctx *GlobalCallContext) {}

// ExitGlobalCall is called when production GlobalCall is exited.
func (s *BaseSLListener) ExitGlobalCall(ctx *GlobalCallContext) {}

// EnterNested is called when production Nested is entered.
func (s *BaseSLListener) EnterNested(ctx *NestedContext) {}

// ExitNested is called when production Nested is exited.
func (s *BaseSLListener) ExitNested(ctx *NestedContext) {}

// EnterCreateList is called when production CreateList is entered.
func (s *BaseSLListener) EnterCreateList(ctx *CreateListContext) {}

// ExitCreateList is called when production CreateList is exited.
func (s *BaseSLListener) ExitCreateList(ctx *CreateListContext) {}

// EnterCreateStruct is called when production CreateStruct is entered.
func (s *BaseSLListener) EnterCreateStruct(ctx *CreateStructContext) {}

// ExitCreateStruct is called when production CreateStruct is exited.
func (s *BaseSLListener) ExitCreateStruct(ctx *CreateStructContext) {}

// EnterCreateMessage is called when production CreateMessage is entered.
func (s *BaseSLListener) EnterCreateMessage(ctx *CreateMessageContext) {}

// ExitCreateMessage is called when production CreateMessage is exited.
func (s *BaseSLListener) ExitCreateMessage(ctx *CreateMessageContext) {}

// EnterConstantLiteral is called when production ConstantLiteral is entered.
func (s *BaseSLListener) EnterConstantLiteral(ctx *ConstantLiteralContext) {}

// ExitConstantLiteral is called when production ConstantLiteral is exited.
func (s *BaseSLListener) ExitConstantLiteral(ctx *ConstantLiteralContext) {}

// EnterExprList is called when production exprList is entered.
func (s *BaseSLListener) EnterExprList(ctx *ExprListContext) {}

// ExitExprList is called when production exprList is exited.
func (s *BaseSLListener) ExitExprList(ctx *ExprListContext) {}

// EnterListInit is called when production listInit is entered.
func (s *BaseSLListener) EnterListInit(ctx *ListInitContext) {}

// ExitListInit is called when production listInit is exited.
func (s *BaseSLListener) ExitListInit(ctx *ListInitContext) {}

// EnterFieldInitializerList is called when production fieldInitializerList is entered.
func (s *BaseSLListener) EnterFieldInitializerList(ctx *FieldInitializerListContext) {}

// ExitFieldInitializerList is called when production fieldInitializerList is exited.
func (s *BaseSLListener) ExitFieldInitializerList(ctx *FieldInitializerListContext) {}

// EnterOptField is called when production optField is entered.
func (s *BaseSLListener) EnterOptField(ctx *OptFieldContext) {}

// ExitOptField is called when production optField is exited.
func (s *BaseSLListener) ExitOptField(ctx *OptFieldContext) {}

// EnterMapInitializerList is called when production mapInitializerList is entered.
func (s *BaseSLListener) EnterMapInitializerList(ctx *MapInitializerListContext) {}

// ExitMapInitializerList is called when production mapInitializerList is exited.
func (s *BaseSLListener) ExitMapInitializerList(ctx *MapInitializerListContext) {}

// EnterSimpleIdentifier is called when production SimpleIdentifier is entered.
func (s *BaseSLListener) EnterSimpleIdentifier(ctx *SimpleIdentifierContext) {}

// ExitSimpleIdentifier is called when production SimpleIdentifier is exited.
func (s *BaseSLListener) ExitSimpleIdentifier(ctx *SimpleIdentifierContext) {}

// EnterEscapedIdentifier is called when production EscapedIdentifier is entered.
func (s *BaseSLListener) EnterEscapedIdentifier(ctx *EscapedIdentifierContext) {}

// ExitEscapedIdentifier is called when production EscapedIdentifier is exited.
func (s *BaseSLListener) ExitEscapedIdentifier(ctx *EscapedIdentifierContext) {}

// EnterOptExpr is called when production optExpr is entered.
func (s *BaseSLListener) EnterOptExpr(ctx *OptExprContext) {}

// ExitOptExpr is called when production optExpr is exited.
func (s *BaseSLListener) ExitOptExpr(ctx *OptExprContext) {}

// EnterInt is called when production Int is entered.
func (s *BaseSLListener) EnterInt(ctx *IntContext) {}

// ExitInt is called when production Int is exited.
func (s *BaseSLListener) ExitInt(ctx *IntContext) {}

// EnterUint is called when production Uint is entered.
func (s *BaseSLListener) EnterUint(ctx *UintContext) {}

// ExitUint is called when production Uint is exited.
func (s *BaseSLListener) ExitUint(ctx *UintContext) {}

// EnterDouble is called when production Double is entered.
func (s *BaseSLListener) EnterDouble(ctx *DoubleContext) {}

// ExitDouble is called when production Double is exited.
func (s *BaseSLListener) ExitDouble(ctx *DoubleContext) {}

// EnterString is called when production String is entered.
func (s *BaseSLListener) EnterString(ctx *StringContext) {}

// ExitString is called when production String is exited.
func (s *BaseSLListener) ExitString(ctx *StringContext) {}

// EnterBytes is called when production Bytes is entered.
func (s *BaseSLListener) EnterBytes(ctx *BytesContext) {}

// ExitBytes is called when production Bytes is exited.
func (s *BaseSLListener) ExitBytes(ctx *BytesContext) {}

// EnterBoolTrue is called when production BoolTrue is entered.
func (s *BaseSLListener) EnterBoolTrue(ctx *BoolTrueContext) {}

// ExitBoolTrue is called when production BoolTrue is exited.
func (s *BaseSLListener) ExitBoolTrue(ctx *BoolTrueContext) {}

// EnterBoolFalse is called when production BoolFalse is entered.
func (s *BaseSLListener) EnterBoolFalse(ctx *BoolFalseContext) {}

// ExitBoolFalse is called when production BoolFalse is exited.
func (s *BaseSLListener) ExitBoolFalse(ctx *BoolFalseContext) {}

// EnterNull is called when production Null is entered.
func (s *BaseSLListener) EnterNull(ctx *NullContext) {}

// ExitNull is called when production Null is exited.
func (s *BaseSLListener) ExitNull(ctx *NullContext) {}
