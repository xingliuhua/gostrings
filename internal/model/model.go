package model

import "encoding/xml"

type StringRes struct {
	XMLName      xml.Name      `xml:"resources"`
	Strings      []StringItem  `xml:"string"`
	StringArrays []StringArray `xml:"string_array"`
}
type StringArray struct {
	Name  string   `xml:"name,attr"`
	Items []string `xml:"item"`
}

type StringItem struct {
	Name string `xml:"name,attr"`
	Text string `xml:",chardata"`
}
