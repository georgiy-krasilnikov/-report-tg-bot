package doc

import (
	"baliance.com/gooxml/document"
	"github.com/lukasjarosch/go-docx"
)

type Data struct {
	Event string
	How   string
	Date  string
	Time  string
	Table Table
}

type Table struct {
	Items []Item
	Cars  []Car
}

type Item struct {
	Name  string
	Count string
}

type Car struct {
	Brand     string
	Number    string
	FullName  string
	Telephone string
}

type Doc struct {
	DocName    string
	DocPath    string
	DocX       *docx.Document
	Doc        *document.Document
	ReplaceMap docx.PlaceholderMap
}
