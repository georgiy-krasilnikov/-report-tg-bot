package services

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"baliance.com/gooxml/document"
	"github.com/lukasjarosch/go-docx"
)

func (h *Handler) NewDoc() error {
	doc, err := docx.Open("file.docx")
	if err != nil {
		return fmt.Errorf("failed to open doc: %s", err.Error())
	}

	h.doc = &Doc{
		DocName: "Рапорт." + h.data.Date + ".docx",
		DocPath: "docs/Рапорт." + h.data.Date + ".docx",
		DocX:    doc,
		ReplaceMap: docx.PlaceholderMap{
			"дд.мм.гггг": h.data.Date,
			"xxx":        h.data.Event,
			"yyy":        h.data.How,
			"zzz":        h.data.Time,
		},
	}

	return nil
}

func (h *Handler) CreateDocument() error {
	if err := h.NewDoc(); err != nil || h.doc == nil {
		return fmt.Errorf("failed to create new doc: %s", err.Error())
	}

	if err := h.doc.DocX.ReplaceAll(h.doc.ReplaceMap); err != nil {
		return fmt.Errorf("failed to replace letters: %s", err.Error())
	}

	if err := h.doc.DocX.WriteToFile(h.doc.DocPath); err != nil {
		return fmt.Errorf("failed to write to file: %s", err.Error())
	}

	doc, err := document.Open(h.doc.DocPath)
	if err != nil {
		return fmt.Errorf("failed to open document: %s", err.Error())
	}

	for i := 0; i < h.data.Table.ItemsNumber; i++ {
		row := doc.Tables()[1].InsertRowAfter(doc.Tables()[1].Rows()[i])
		for i := 0; i < 5; i++ {
			row.AddCell().AddParagraph()
		}

		row.Cells()[0].Paragraphs()[0].AddRun().AddText(strconv.Itoa(i + 1))
		row.Cells()[1].Paragraphs()[0].AddRun().AddText(h.data.Table.Items[i].Name)
		row.Cells()[2].Paragraphs()[0].AddRun().AddText(h.data.Table.Items[i].Count)
	}

	for i := h.data.Table.ItemsNumber; i-h.data.Table.ItemsNumber < h.data.Table.CarsNumber; i++ {
		row := doc.Tables()[1].InsertRowAfter(doc.Tables()[1].Rows()[i+2])

		row.AddCell().AddParagraph().AddRun().AddText(strconv.Itoa(i - h.data.Table.ItemsNumber + 1))
		row.AddCell().AddParagraph().AddRun().AddText(h.data.Table.Cars[i-h.data.Table.ItemsNumber].Brand)
		row.AddCell().AddParagraph().AddRun().AddText(h.data.Table.Cars[i-h.data.Table.ItemsNumber].Number)
		row.AddCell().AddParagraph().AddRun().AddText(h.data.Table.Cars[i-h.data.Table.ItemsNumber].FullName)
		row.AddCell().AddParagraph().AddRun().AddText(h.data.Table.Cars[i-h.data.Table.ItemsNumber].Telephone)
	}

	if err := doc.SaveToFile(h.doc.DocPath); err != nil {
		return fmt.Errorf("failed to save new file: %s", err.Error())
	}

	return nil
}

func GetListOfDocuments() ([]string, error) {
	m, err := filepath.Glob("docs/*.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to get list of files: %s", err.Error())
	}

	var lst []string
	for _, v := range m {
		lst = append(lst, strings.TrimPrefix(v, "docs/"))
	}

	return lst, nil
}

// func (h *Handler) DeleteDocument() error {
// 	lst, err := GetListOfDocuments()
// 	if err != nil {
// 		return fmt.Errorf("failed to get list of docs: %s", err.Error())
// 	}

// 	for _, name := range lst {
// 		if time.Now().Format("01.02.2006") == strings.TrimPrefix(name, "Рапорт.") && strconv.Itoa(time.Now().Hour()) == "23" && strconv.Itoa(time.Now().Minute()) == "59" {
// 			if err := os.Remove("docs/" + h.doc.DocName); err != nil {
// 				return fmt.Errorf("failed to delete document: %s", err.Error())
// 			}
// 		}
// 	}

// 	return nil
// }

func (h *Handler) EditDate() error {
	_, err := time.Parse(h.data.Date, "02.01.2006")
	if err != nil {
		return fmt.Errorf("invalid date format: %s", err.Error())
	}

	doc, err := document.Open(h.doc.DocPath)
	if err != nil {
		return fmt.Errorf("failed to open document: %s", err.Error())
	}

	doc.Paragraphs()[4].AddRun().AddText(h.data.Date)
	doc.Paragraphs()[4].SetStyle("Text Body")
	doc.Paragraphs()[4].RemoveRun(doc.Paragraphs()[4].Runs()[0])

	if err := doc.SaveToFile("docs/Рапорт." + h.data.Date + ".docx"); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}

func (h *Handler) GetListOfItems() ([][]string, error) {
	doc, err := document.Open(h.doc.DocPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open document: %s", err.Error())
	}

	var lst [][]string
	for i := 1; strconv.Itoa(i) == doc.Tables()[1].Rows()[i].Cells()[0].Paragraphs()[0].Runs()[0].Text(); i++ {
		var row []string
		for j := 0; j < 3; j++ {
			row = append(row, doc.Tables()[1].Rows()[i].Cells()[j].Paragraphs()[0].Runs()[0].Text())
		}
		lst = append(lst, row)
	}

	return lst, nil
}

func (h *Handler) EditRow(id string) error {
	doc, err := document.Open(h.doc.DocPath)
	if err != nil {
		return fmt.Errorf("failed to open document: %s", err.Error())
	}

	for i := 1; strconv.Itoa(i) == doc.Tables()[1].Rows()[i].Cells()[0].Paragraphs()[0].Runs()[0].Text(); i++ {
		if strconv.Itoa(i) == id {
			doc.Tables()[1].Rows()[i].Cells()[1].Paragraphs()[0].Runs()[0].ClearContent()
			doc.Tables()[1].Rows()[i].Cells()[2].Paragraphs()[0].Runs()[0].ClearContent()

		}
	}

	if err := doc.SaveToFile(h.doc.DocPath); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}
