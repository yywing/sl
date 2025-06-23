package functions

import (
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/lib/types"
	"github.com/yywing/sl/native"
)

func init() {
	LibFunctions[FunctionXMLPath] = ast.NewBaseFunction(
		FunctionXMLPath,
		native.MustNewNativeFunction(FunctionXMLPath, XMLPath).Definitions(),
	)
	LibFunctions[FunctionXMLAttr] = ast.NewBaseFunction(
		FunctionXMLAttr,
		native.MustNewNativeFunction(FunctionXMLAttr, XMLAttr).Definitions(),
	)
	LibFunctions[FunctionXMLElement] = ast.NewBaseFunction(
		FunctionXMLElement,
		native.MustNewNativeFunction(FunctionXMLElement, XMLElement).Definitions(),
	)
	LibFunctions[FunctionXMLText] = ast.NewBaseFunction(
		FunctionXMLText,
		native.MustNewNativeFunction(FunctionXMLText, XMLText).Definitions(),
	)
}

const (
	FunctionXMLPath    = "xmlPath"
	FunctionXMLAttr    = "xmlAttr"
	FunctionXMLElement = "xmlElement"
	FunctionXMLText    = "xmlText"
)

func XMLPath(json, path string) ([]*types.XMLValue, error) {
	doc, err := xmlquery.Parse(strings.NewReader(json))
	if err != nil {
		return nil, err
	}

	nodes := xmlquery.Find(doc, path)

	values := make([]*types.XMLValue, 0, len(nodes))
	for _, node := range nodes {
		values = append(values, types.NewXMLValue(node.OutputXML(true)))
	}

	return values, nil
}

func node(v *types.XMLValue) (*xmlquery.Node, error) {
	return xmlquery.Parse(strings.NewReader(v.XML))
}

func XMLAttr(v *types.XMLValue, attr string) (string, error) {
	n, err := node(v)
	if err != nil {
		return "", err
	}
	// get element node
	return n.LastChild.SelectAttr(attr), nil
}

func XMLElement(v *types.XMLValue, element string) ([]*types.XMLValue, error) {
	n, err := node(v)
	if err != nil {
		return nil, err
	}
	// get element node
	nodes := n.LastChild.SelectElements(element)
	values := make([]*types.XMLValue, 0, len(nodes))
	for _, node := range nodes {
		values = append(values, types.NewXMLValue(node.OutputXML(true)))
	}
	return values, nil
}

func XMLText(v *types.XMLValue) (string, error) {
	n, err := node(v)
	if err != nil {
		return "", err
	}
	// get element node
	return n.LastChild.InnerText(), nil
}
