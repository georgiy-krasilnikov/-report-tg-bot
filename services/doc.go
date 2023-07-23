package services

import (
	"fmt"
	"os"
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
		return fmt.Errorf("failed to replace: %s", err.Error())
	}

	if err := h.doc.DocX.WriteToFile("docs/" + h.doc.DocName); err != nil {
		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	doc, err := document.Open("docs/" + h.doc.DocName)
	if err != nil {
		return fmt.Errorf("error opening document: %s", err.Error())
	}

	for i := 0; i < h.data.Count; i++ {
		row := doc.Tables()[1].InsertRowAfter(doc.Tables()[1].Rows()[i])
		for i := 0; i < 5; i++ {
			row.AddCell().AddParagraph()
		}
		row.Cells()[0].Paragraphs()[0].AddRun().AddText(strconv.Itoa(i + 1))
		row.Cells()[1].Paragraphs()[0].AddRun().AddText(h.data.Items[i])
		row.Cells()[2].Paragraphs()[0].AddRun().AddText(h.data.CountItems[i])
	}

	if err := doc.SaveToFile("docs/" + h.doc.DocName); err != nil {
		return fmt.Errorf("failed to save replaced file: %s", err.Error())
	}

	return nil
}

func GetListOfDocuments() ([]string, error) {
	m, err := filepath.Glob("docs/*.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to get list of names of files: %s", err.Error())
	}

	var lst []string
	for _, v := range m {
		lst = append(lst, strings.TrimPrefix(v, "docs/"))
	}

	return lst, nil
}

func (h *Handler) DeleteDocument() error {
	lst, err := GetListOfDocuments()
	if err != nil {
		return fmt.Errorf("failed to get list of docs: %s", err.Error())
	}

	for _, name := range lst {
		if time.Now().Format("01.02.2006") == strings.TrimPrefix(name, "Рапорт.") && strconv.Itoa(time.Now().Hour()) == "23" && strconv.Itoa(time.Now().Minute()) == "59" {
			if err := os.Remove("docs/" + h.doc.DocName); err != nil {
				return fmt.Errorf("failed to delete document: %s", err.Error())
			}
		}
	}

	return nil
}
