package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"baliance.com/gooxml/document"
	// "baliance.com/gooxml/measurement"
	// "baliance.com/gooxml/schema/soo/wml"
	"github.com/lukasjarosch/go-docx"
)

func (h *Handler) NewDocX() error {
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

func (h *Handler) NewDoc(name, path string) error {
	doc, err := document.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open document: %s", err.Error())
	}

	h.doc.Doc = doc
	h.doc.DocName = name
	h.doc.DocPath = path

	return nil
}

func (h *Handler) CreateDocument() error {
	if err := h.NewDocX(); err != nil || h.doc == nil {
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

	for i := 0; i < len(h.data.Table.Items); i++ {
		row := doc.Tables()[1].InsertRowAfter(doc.Tables()[1].Rows()[i])
		for i := 0; i < 5; i++ {
			row.AddCell().AddParagraph()
		}

		row.Cells()[0].Paragraphs()[0].AddRun().AddText(strconv.Itoa(i + 1))
		row.Cells()[1].Paragraphs()[0].AddRun().AddText(h.data.Table.Items[i].Name)
		row.Cells()[2].Paragraphs()[0].AddRun().AddText(h.data.Table.Items[i].Count)
	}

	for i := len(h.data.Table.Items); i-len(h.data.Table.Items) < len(h.data.Table.Cars); i++ {
		row := doc.Tables()[1].InsertRowAfter(doc.Tables()[1].Rows()[i+2])

		row.AddCell().AddParagraph().AddRun().AddText(strconv.Itoa(i - len(h.data.Table.Items) + 1))
		row.AddCell().AddParagraph().AddRun().AddText(h.data.Table.Cars[i-len(h.data.Table.Items)].Brand)
		row.AddCell().AddParagraph().AddRun().AddText(h.data.Table.Cars[i-len(h.data.Table.Items)].Number)
		row.AddCell().AddParagraph().AddRun().AddText(h.data.Table.Cars[i-len(h.data.Table.Items)].FullName)
		row.AddCell().AddParagraph().AddRun().AddText(h.data.Table.Cars[i-len(h.data.Table.Items)].Telephone)
	}

	if err := doc.SaveToFile(h.doc.DocPath); err != nil {
		return fmt.Errorf("failed to save new file: %s", err.Error())
	}

	return nil
}

func GetListOfDocuments() ([]string, error) {
	matches, err := filepath.Glob("docs/*.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to get list of files: %s", err.Error())
	}

	var docs []string
	for _, m := range matches {
		docs = append(docs, strings.TrimPrefix(m, "docs/"))
	}

	return docs, nil
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
	h.doc.Doc.Paragraphs()[4].Runs()[0].Clear()
	h.doc.Doc.Paragraphs()[4].Runs()[0].AddText(h.data.Date)

	if err := os.Remove(h.doc.DocPath); err != nil {
		return fmt.Errorf("failed to remove document: %s", err.Error())
	}

	h.doc.DocName = "Рапорт." + h.data.Date + ".docx"
	h.doc.DocPath = "docs/" + h.doc.DocName

	if err := h.doc.Doc.SaveToFile(h.doc.DocPath); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}

func (h *Handler) GetListOfItems() ([][]string, error) {
	var lst [][]string
	for i := 1; strconv.Itoa(i) == h.doc.Doc.Tables()[1].Rows()[i].Cells()[0].Paragraphs()[0].Runs()[0].Text(); i++ {
		var row []string
		for j := 0; j < 3; j++ {
			row = append(row, h.doc.Doc.Tables()[1].Rows()[i].Cells()[j].Paragraphs()[0].Runs()[0].Text())
		}
		lst = append(lst, row)
	}

	return lst, nil
}

func (h *Handler) GetListOfCars() ([][]string, error) {
	items, err := h.GetListOfItems()
	if err != nil {
		return nil, fmt.Errorf("failed to get list of items; %s", err.Error())
	}

	var id int
	var lst [][]string
	if len(h.doc.Doc.Tables()[1].Rows())-len(items) > 3 {
		id = len(items) + 3
	} else {
		return nil, nil
	}

	for i := id; i < len(h.doc.Doc.Tables()[1].Rows()); i++ {
		var row []string
		for j := 0; j < 5; j++ {
			row = append(row, h.doc.Doc.Tables()[1].Rows()[i].Cells()[j].Paragraphs()[0].Runs()[0].Text())
		}
		lst = append(lst, row)
	}

	return lst, nil
}

func (h *Handler) EditItemRow(id *string) error {
	for i := 1; strconv.Itoa(i) == h.doc.Doc.Tables()[1].Rows()[i].Cells()[0].Paragraphs()[0].Runs()[0].Text(); i++ {
		if strconv.Itoa(i) == *id {
			h.doc.Doc.Tables()[1].Rows()[i].Cells()[1].Paragraphs()[0].Runs()[0].ClearContent()
			h.doc.Doc.Tables()[1].Rows()[i].Cells()[2].Paragraphs()[0].Runs()[0].ClearContent()
			h.doc.Doc.Tables()[1].Rows()[i].Cells()[1].Paragraphs()[0].Runs()[0].AddText(h.data.Table.Items[0].Name)
			h.doc.Doc.Tables()[1].Rows()[i].Cells()[2].Paragraphs()[0].Runs()[0].AddText(h.data.Table.Items[0].Count)
		}
	}
	*id = ""

	if err := h.doc.Doc.SaveToFile(h.doc.DocPath); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}

func (h *Handler) AddItemRow() error {
	items, err := h.GetListOfItems()
	if err != nil {
		return fmt.Errorf("failed to get list of items: %s", err.Error())
	}

	id := len(items)

	for i := 0; i < len(h.data.Table.Items); i++ {
		row := h.doc.Doc.Tables()[1].InsertRowAfter(h.doc.Doc.Tables()[1].Rows()[id])
		for j := 0; j < 5; j++ {
			row.AddCell().AddParagraph()
		}

		row.Cells()[0].Paragraphs()[0].AddRun().AddText(strconv.Itoa(id + 1))
		row.Cells()[1].Paragraphs()[0].AddRun().AddText(h.data.Table.Items[i].Name)
		row.Cells()[2].Paragraphs()[0].AddRun().AddText(h.data.Table.Items[i].Count)
		id++
	}

	if err := h.doc.Doc.SaveToFile(h.doc.DocPath); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}

func (h *Handler) EditCarRow(id *string) error {
	items, err := h.GetListOfItems()
	if err != nil {
		return fmt.Errorf("failed to get list of items: %s", err.Error())
	}

	for i := len(items) + 2; i < len(h.doc.Doc.Tables()[1].Rows()); i++ {
		if h.doc.Doc.Tables()[1].Rows()[i].Cells()[0].Paragraphs()[0].Runs()[0].Text() == *id {
			for j := 1; j < 5; j++ {
				h.doc.Doc.Tables()[1].Rows()[i].Cells()[j].Paragraphs()[0].Runs()[0].ClearContent()
			}

			h.doc.Doc.Tables()[1].Rows()[i].Cells()[1].Paragraphs()[0].Runs()[0].AddText(h.data.Table.Cars[0].Brand)
			h.doc.Doc.Tables()[1].Rows()[i].Cells()[2].Paragraphs()[0].Runs()[0].AddText(h.data.Table.Cars[0].Number)
			h.doc.Doc.Tables()[1].Rows()[i].Cells()[3].Paragraphs()[0].Runs()[0].AddText(h.data.Table.Cars[0].FullName)
			h.doc.Doc.Tables()[1].Rows()[i].Cells()[4].Paragraphs()[0].Runs()[0].AddText(h.data.Table.Cars[0].Telephone)
		}
	}
	*id = ""

	if err := h.doc.Doc.SaveToFile(h.doc.DocPath); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}

func (h *Handler) AddCarRow() error {
	cars, err := h.GetListOfCars()
	if err != nil {
		return fmt.Errorf("failed to get list of cars: %s", err.Error())
	}

	id := len(h.doc.Doc.Tables()[1].Rows()) - len(cars)
	if id == len(h.doc.Doc.Tables()[1].Rows()) {
		s := h.doc.Doc.Paragraphs()[6].Runs()[0].Text()
		if strings.Contains(s, "КПП №1") {
			h.doc.Doc.Paragraphs()[6].RemoveRun(h.doc.Doc.Paragraphs()[6].Runs()[0])
			s = strings.ReplaceAll(s, "КПП №1", "гаражный въезд")
			h.doc.Doc.Paragraphs()[6].AddRun().AddText(s)
		}

		id -= 1
	}

	for i := 0; i < len(h.data.Table.Cars); i++ {
		row := h.doc.Doc.Tables()[1].InsertRowAfter(h.doc.Doc.Tables()[1].Rows()[id])
		for j := 0; j < 5; j++ {
			row.AddCell().AddParagraph()
		}

		row.Cells()[0].Paragraphs()[0].AddRun().AddText(strconv.Itoa(len(cars) + i + 1))
		row.Cells()[1].Paragraphs()[0].AddRun().AddText(h.data.Table.Cars[i].Brand)
		row.Cells()[2].Paragraphs()[0].AddRun().AddText(h.data.Table.Cars[i].Number)
		row.Cells()[3].Paragraphs()[0].AddRun().AddText(h.data.Table.Cars[i].FullName)
		row.Cells()[4].Paragraphs()[0].AddRun().AddText(h.data.Table.Cars[i].Telephone)
		id++
	}

	if err := h.doc.Doc.SaveToFile(h.doc.DocPath); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}
