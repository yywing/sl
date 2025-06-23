package types

import "github.com/yywing/sl/ast"

const (
	TypeKindXML = "xml"
)

var (
	XMLType = ast.NewPrimitiveType(TypeKindXML, 0)
)

type XMLValue struct {
	XML string
}

func NewXMLValue(xml string) *XMLValue {
	return &XMLValue{
		XML: xml,
	}
}

func (v *XMLValue) Type() ast.ValueType {
	return XMLType
}

func (v *XMLValue) String() string {
	return v.XML
}

func (v *XMLValue) Equal(other ast.Value) bool {
	otherValue, ok := other.(*XMLValue)
	if !ok {
		return false
	}
	return v.XML == otherValue.XML
}
