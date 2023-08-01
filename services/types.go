package services

import "github.com/lukasjarosch/go-docx"

type Data struct {
	Event string
	How   string
	Date  string
	Time  string
	Table Table
	Full  bool
}

type Table struct {
	ItemsNumber int
	Items       []Item
	CarsNumber  int
	Cars        []Car
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
	ReplaceMap docx.PlaceholderMap
}
