package main

import "encoding/xml"

type SAMLAttribute struct {
	Name  string   `xml:"Name,attr"`
	Value []string `xml:"AttributeValue"`
}

type SAMLXml struct {
	XMLName    xml.Name        `xml:"Response"`
	Attributes []SAMLAttribute `xml:"Assertion>AttributeStatement>Attribute"`
}
