// Code generated from /root/self/sl/SL.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // SL

import "github.com/antlr4-go/antlr/v4"

type BaseSLVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseSLVisitor) VisitStart(ctx *StartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitExpr(ctx *ExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitConditionalOr(ctx *ConditionalOrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitConditionalAnd(ctx *ConditionalAndContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitRelation(ctx *RelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitCalc(ctx *CalcContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitMemberExpr(ctx *MemberExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitLogicalNot(ctx *LogicalNotContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitNegate(ctx *NegateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitMemberCall(ctx *MemberCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitSelect(ctx *SelectContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitPrimaryExpr(ctx *PrimaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitIndex(ctx *IndexContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitIdent(ctx *IdentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitGlobalCall(ctx *GlobalCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitNested(ctx *NestedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitCreateList(ctx *CreateListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitCreateStruct(ctx *CreateStructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitCreateMessage(ctx *CreateMessageContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitConstantLiteral(ctx *ConstantLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitExprList(ctx *ExprListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitListInit(ctx *ListInitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitFieldInitializerList(ctx *FieldInitializerListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitOptField(ctx *OptFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitMapInitializerList(ctx *MapInitializerListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitSimpleIdentifier(ctx *SimpleIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitEscapedIdentifier(ctx *EscapedIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitOptExpr(ctx *OptExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitInt(ctx *IntContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitUint(ctx *UintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitDouble(ctx *DoubleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitString(ctx *StringContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitBytes(ctx *BytesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitBoolTrue(ctx *BoolTrueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitBoolFalse(ctx *BoolFalseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSLVisitor) VisitNull(ctx *NullContext) interface{} {
	return v.VisitChildren(ctx)
}
