package xml

import (
	"github.com/beevik/etree"
)

type XMLNode struct {
	Tag       string
	TagFull   string
	Namespace string
	Text      string
	Attribute map[string]string
	Items     []XMLNode
}

// "github.com/beevik/etree"
// 利用etree 解析xmlString
func XmlParse(xmlString string) (retNode XMLNode, err error) {
	// map[string][]map[string]string,
	xmldoc := etree.NewDocument()
	if err = xmldoc.ReadFromString(xmlString); err != nil {
		return
	}
	retNode, err = XmlParseNode(xmldoc.Root())
	return
}

func XmlParseNode(element *etree.Element) (XMLNode, error) {
	var xmlnodes XMLNode
	attr := make(map[string]string)
	xmlnodes.Tag = element.Tag
	xmlnodes.TagFull = element.FullTag()
	xmlnodes.Namespace = element.NamespaceURI()
	xmlnodes.Text = element.Text()
	for _, att := range element.Attr {
		attr[att.Key+"_namespace"] = att.NamespaceURI()
		attr[att.Key] = att.Value
		attr[att.FullKey()] = att.Value
	}
	xmlnodes.Attribute = attr

	// var items []map[string]interface{}
	for _, elemt := range element.ChildElements() {
		node, errs := XmlParseNode(elemt)
		if errs == nil {
			xmlnodes.Items = append(xmlnodes.Items, node)
		}
	}
	return xmlnodes, nil
}
func (xmlNode *XMLNode) GetItem(idx int) (item XMLNode) {
	if len(xmlNode.Items) <= idx {
		// err = fmt.Errorf("超出范围")
		return
	}
	item = xmlNode.Items[idx]
	return
}
