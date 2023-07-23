package services

import "github.com/lukasjarosch/go-docx"

type Data struct {
	Event      string
	How        string
	Date       string
	Time       string
	Count      int
	Items      []string
	CountItems []string
}

type Doc struct {
	DocName    string
	DocX       *docx.Document
	ReplaceMap docx.PlaceholderMap
}
