// Code generated from sl/SL.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // SL

import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by SLParser.
type SLVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by SLParser#start.
	VisitStart(ctx *StartContext) interface{}

	// Visit a parse tree produced by SLParser#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by SLParser#conditionalOr.
	VisitConditionalOr(ctx *ConditionalOrContext) interface{}

	// Visit a parse tree produced by SLParser#conditionalAnd.
	VisitConditionalAnd(ctx *ConditionalAndContext) interface{}

	// Visit a parse tree produced by SLParser#relation.
	VisitRelation(ctx *RelationContext) interface{}

	// Visit a parse tree produced by SLParser#calc.
	VisitCalc(ctx *CalcContext) interface{}

	// Visit a parse tree produced by SLParser#MemberExpr.
	VisitMemberExpr(ctx *MemberExprContext) interface{}

	// Visit a parse tree produced by SLParser#LogicalNot.
	VisitLogicalNot(ctx *LogicalNotContext) interface{}

	// Visit a parse tree produced by SLParser#Negate.
	VisitNegate(ctx *NegateContext) interface{}

	// Visit a parse tree produced by SLParser#MemberCall.
	VisitMemberCall(ctx *MemberCallContext) interface{}

	// Visit a parse tree produced by SLParser#Select.
	VisitSelect(ctx *SelectContext) interface{}

	// Visit a parse tree produced by SLParser#PrimaryExpr.
	VisitPrimaryExpr(ctx *PrimaryExprContext) interface{}

	// Visit a parse tree produced by SLParser#Index.
	VisitIndex(ctx *IndexContext) interface{}

	// Visit a parse tree produced by SLParser#Ident.
	VisitIdent(ctx *IdentContext) interface{}

	// Visit a parse tree produced by SLParser#GlobalCall.
	VisitGlobalCall(ctx *GlobalCallContext) interface{}

	// Visit a parse tree produced by SLParser#Nested.
	VisitNested(ctx *NestedContext) interface{}

	// Visit a parse tree produced by SLParser#CreateList.
	VisitCreateList(ctx *CreateListContext) interface{}

	// Visit a parse tree produced by SLParser#CreateStruct.
	VisitCreateStruct(ctx *CreateStructContext) interface{}

	// Visit a parse tree produced by SLParser#CreateMessage.
	VisitCreateMessage(ctx *CreateMessageContext) interface{}

	// Visit a parse tree produced by SLParser#ConstantLiteral.
	VisitConstantLiteral(ctx *ConstantLiteralContext) interface{}

	// Visit a parse tree produced by SLParser#exprList.
	VisitExprList(ctx *ExprListContext) interface{}

	// Visit a parse tree produced by SLParser#listInit.
	VisitListInit(ctx *ListInitContext) interface{}

	// Visit a parse tree produced by SLParser#fieldInitializerList.
	VisitFieldInitializerList(ctx *FieldInitializerListContext) interface{}

	// Visit a parse tree produced by SLParser#optField.
	VisitOptField(ctx *OptFieldContext) interface{}

	// Visit a parse tree produced by SLParser#mapInitializerList.
	VisitMapInitializerList(ctx *MapInitializerListContext) interface{}

	// Visit a parse tree produced by SLParser#SimpleIdentifier.
	VisitSimpleIdentifier(ctx *SimpleIdentifierContext) interface{}

	// Visit a parse tree produced by SLParser#EscapedIdentifier.
	VisitEscapedIdentifier(ctx *EscapedIdentifierContext) interface{}

	// Visit a parse tree produced by SLParser#optExpr.
	VisitOptExpr(ctx *OptExprContext) interface{}

	// Visit a parse tree produced by SLParser#Int.
	VisitInt(ctx *IntContext) interface{}

	// Visit a parse tree produced by SLParser#Uint.
	VisitUint(ctx *UintContext) interface{}

	// Visit a parse tree produced by SLParser#Double.
	VisitDouble(ctx *DoubleContext) interface{}

	// Visit a parse tree produced by SLParser#String.
	VisitString(ctx *StringContext) interface{}

	// Visit a parse tree produced by SLParser#Bytes.
	VisitBytes(ctx *BytesContext) interface{}

	// Visit a parse tree produced by SLParser#BoolTrue.
	VisitBoolTrue(ctx *BoolTrueContext) interface{}

	// Visit a parse tree produced by SLParser#BoolFalse.
	VisitBoolFalse(ctx *BoolFalseContext) interface{}

	// Visit a parse tree produced by SLParser#Null.
	VisitNull(ctx *NullContext) interface{}
}
