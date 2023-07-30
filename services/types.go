package services

import "github.com/lukasjarosch/go-docx"

type Data struct {
	Event string
	How   string
	Date  string
	Time  string
	Table Table
}

type Table struct {
	ItemsNumber int
	Items       []string
	CountItems  []string
	CarsNumber  int
	Cars        []string
	CountCars   []string
}

type Doc struct {
	DocName    string
	DocPath    string
	DocX       *docx.Document
	ReplaceMap docx.PlaceholderMap
}
