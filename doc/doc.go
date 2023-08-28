package doc

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"baliance.com/gooxml/document"
	"github.com/lukasjarosch/go-docx"
)

func (d *Data) NewDocX() (*Doc, error) {
	doc, err := docx.Open("docs/file.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to open doc: %s", err.Error())
	}

	return &Doc{
		DocName: "Рапорт." + d.Date + ".docx",
		DocPath: "docs/docs/Рапорт." + d.Date + ".docx",
		DocX:    doc,
		ReplaceMap: docx.PlaceholderMap{
			"дд.мм.гггг": d.Date,
			"xxx":        d.Event,
			"yyy":        d.How,
			"zzz":        d.Time,
		},
	}, nil
}

func NewDoc(name, path string) (*Doc, error) {
	doc, err := document.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open document: %s", err.Error())
	}

	return &Doc{
		Doc:     doc,
		DocName: name,
		DocPath: path,
	}, nil
}

func (d *Data) CreateDocument() (*Doc, error) {
	docx, err := d.NewDocX()
	if err != nil || docx == nil {
		return nil, fmt.Errorf("failed to create new doc: %s", err.Error())
	}

	if err := docx.DocX.ReplaceAll(docx.ReplaceMap); err != nil {
		return nil, fmt.Errorf("failed to replace letters: %s", err.Error())
	}

	if err := docx.DocX.WriteToFile(docx.DocPath); err != nil {
		return nil, fmt.Errorf("failed to write to file: %s", err.Error())
	}

	doc, err := document.Open(docx.DocPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open document: %s", err.Error())
	}

	for i := 0; i < len(d.Table.Items); i++ {
		row := doc.Tables()[1].InsertRowAfter(doc.Tables()[1].Rows()[i])
		for i := 0; i < 5; i++ {
			row.AddCell().AddParagraph()
		}

		row.Cells()[0].Paragraphs()[0].AddRun().AddText(strconv.Itoa(i + 1))
		row.Cells()[1].Paragraphs()[0].AddRun().AddText(d.Table.Items[i].Name)
		row.Cells()[2].Paragraphs()[0].AddRun().AddText(d.Table.Items[i].Count)
	}

	for i := len(d.Table.Items); i-len(d.Table.Items) < len(d.Table.Cars); i++ {
		row := doc.Tables()[1].InsertRowAfter(doc.Tables()[1].Rows()[i+2])

		row.AddCell().AddParagraph().AddRun().AddText(strconv.Itoa(i - len(d.Table.Items) + 1))
		row.AddCell().AddParagraph().AddRun().AddText(d.Table.Cars[i-len(d.Table.Items)].Brand)
		row.AddCell().AddParagraph().AddRun().AddText(d.Table.Cars[i-len(d.Table.Items)].Number)
		row.AddCell().AddParagraph().AddRun().AddText(d.Table.Cars[i-len(d.Table.Items)].FullName)
		row.AddCell().AddParagraph().AddRun().AddText(d.Table.Cars[i-len(d.Table.Items)].Telephone)
	}

	if err := doc.SaveToFile(docx.DocPath); err != nil {
		return nil, fmt.Errorf("failed to save new file: %s", err.Error())
	}

	return docx, nil
}

func GetListOfDocuments() ([]string, error) {
	matches, err := filepath.Glob("docs/docs/*.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to get list of files: %s", err.Error())
	}

	var docs []string
	for _, m := range matches {
		docs = append(docs, strings.TrimPrefix(m, "docs/docs/"))
	}

	return docs, nil
}

// func (d) DeleteDocument() error {
// 	lst, err := GetListOfDocuments()
// 	if err != nil {
// 		return fmt.Errorf("failed to get list of docs: %s", err.Error())
// 	}

// 	for _, name := range lst {
// 		if time.Now().Format("01.02.2006") == strings.TrimPrefix(name, "Рапорт.") && strconv.Itoa(time.Now().Hour()) == "23" && strconv.Itoa(time.Now().Minute()) == "59" {
// 			if err := os.Remove("docs/" + d.DocName); err != nil {
// 				return fmt.Errorf("failed to delete document: %s", err.Error())
// 			}
// 		}
// 	}

// 	return nil
// }

func (d *Doc) EditDate(date string) error {
	d.Doc.Paragraphs()[4].Runs()[0].Clear()
	d.Doc.Paragraphs()[4].Runs()[0].AddText(date)

	if err := os.Remove(d.DocPath); err != nil {
		return fmt.Errorf("failed to remove document: %s", err.Error())
	}

	d.DocName = "Рапорт." + date + ".docx"
	d.DocPath = "docs/docs/" + d.DocName

	if err := d.Doc.SaveToFile(d.DocPath); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}

func (d *Doc) GetListOfItems() ([][]string, error) {
	var lst [][]string
	for i := 1; strconv.Itoa(i) == d.Doc.Tables()[1].Rows()[i].Cells()[0].Paragraphs()[0].Runs()[0].Text(); i++ {
		var row []string
		for j := 0; j < 3; j++ {
			row = append(row, d.Doc.Tables()[1].Rows()[i].Cells()[j].Paragraphs()[0].Runs()[0].Text())
		}
		lst = append(lst, row)
	}

	return lst, nil
}

func (d *Doc) GetListOfCars() ([][]string, error) {
	items, err := d.GetListOfItems()
	if err != nil {
		return nil, fmt.Errorf("failed to get list of items; %s", err.Error())
	}

	var id int
	var lst [][]string
	if len(d.Doc.Tables()[1].Rows())-len(items) > 3 {
		id = len(items) + 3
	} else {
		return nil, nil
	}

	for i := id; i < len(d.Doc.Tables()[1].Rows()); i++ {
		var row []string
		for j := 0; j < 5; j++ {
			row = append(row, d.Doc.Tables()[1].Rows()[i].Cells()[j].Paragraphs()[0].Runs()[0].Text())
		}
		lst = append(lst, row)
	}

	return lst, nil
}

func (d *Doc) EditItemRow(id *string, t *Table) error {
	for i := 1; strconv.Itoa(i) == d.Doc.Tables()[1].Rows()[i].Cells()[0].Paragraphs()[0].Runs()[0].Text(); i++ {
		if strconv.Itoa(i) == *id {
			d.Doc.Tables()[1].Rows()[i].Cells()[1].Paragraphs()[0].Runs()[0].ClearContent()
			d.Doc.Tables()[1].Rows()[i].Cells()[2].Paragraphs()[0].Runs()[0].ClearContent()
			d.Doc.Tables()[1].Rows()[i].Cells()[1].Paragraphs()[0].Runs()[0].AddText(t.Items[0].Name)
			d.Doc.Tables()[1].Rows()[i].Cells()[2].Paragraphs()[0].Runs()[0].AddText(t.Items[0].Count)
		}
	}
	*id = ""

	if err := d.Doc.SaveToFile(d.DocPath); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}

func (d *Doc) AddItemRow(t *Table) error {
	items, err := d.GetListOfItems()
	if err != nil {
		return fmt.Errorf("failed to get list of items: %s", err.Error())
	}

	id := len(items)

	for i := 0; i < len(t.Items); i++ {
		row := d.Doc.Tables()[1].InsertRowAfter(d.Doc.Tables()[1].Rows()[id])
		for j := 0; j < 5; j++ {
			row.AddCell().AddParagraph()
		}

		row.Cells()[0].Paragraphs()[0].AddRun().AddText(strconv.Itoa(id + 1))
		row.Cells()[1].Paragraphs()[0].AddRun().AddText(t.Items[i].Name)
		row.Cells()[2].Paragraphs()[0].AddRun().AddText(t.Items[i].Count)
		id++
	}

	if err := d.Doc.SaveToFile(d.DocPath); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}

func (d *Doc) EditCarRow(id *string, t *Table) error {
	items, err := d.GetListOfItems()
	if err != nil {
		return fmt.Errorf("failed to get list of items: %s", err.Error())
	}

	for i := len(items) + 2; i < len(d.Doc.Tables()[1].Rows()); i++ {
		if d.Doc.Tables()[1].Rows()[i].Cells()[0].Paragraphs()[0].Runs()[0].Text() == *id {
			for j := 1; j < 5; j++ {
				d.Doc.Tables()[1].Rows()[i].Cells()[j].Paragraphs()[0].Runs()[0].ClearContent()
			}

			d.Doc.Tables()[1].Rows()[i].Cells()[1].Paragraphs()[0].Runs()[0].AddText(t.Cars[0].Brand)
			d.Doc.Tables()[1].Rows()[i].Cells()[2].Paragraphs()[0].Runs()[0].AddText(t.Cars[0].Number)
			d.Doc.Tables()[1].Rows()[i].Cells()[3].Paragraphs()[0].Runs()[0].AddText(t.Cars[0].FullName)
			d.Doc.Tables()[1].Rows()[i].Cells()[4].Paragraphs()[0].Runs()[0].AddText(t.Cars[0].Telephone)
		}
	}
	*id = ""

	if err := d.Doc.SaveToFile(d.DocPath); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}

func (d *Doc) AddCarRow(t *Table) error {
	cars, err := d.GetListOfCars()
	if err != nil {
		return fmt.Errorf("failed to get list of cars: %s", err.Error())
	}

	id := len(d.Doc.Tables()[1].Rows()) - len(cars)
	if id == len(d.Doc.Tables()[1].Rows()) {
		s := d.Doc.Paragraphs()[6].Runs()[0].Text()
		if strings.Contains(s, "КПП №1") {
			d.Doc.Paragraphs()[6].RemoveRun(d.Doc.Paragraphs()[6].Runs()[0])
			s = strings.ReplaceAll(s, "КПП №1", "гаражный въезд")
			d.Doc.Paragraphs()[6].AddRun().AddText(s)
		}

		id -= 1
	}

	for i := 0; i < len(t.Cars); i++ {
		row := d.Doc.Tables()[1].InsertRowAfter(d.Doc.Tables()[1].Rows()[id])
		for j := 0; j < 5; j++ {
			row.AddCell().AddParagraph()
		}

		row.Cells()[0].Paragraphs()[0].AddRun().AddText(strconv.Itoa(len(cars) + i + 1))
		row.Cells()[1].Paragraphs()[0].AddRun().AddText(t.Cars[i].Brand)
		row.Cells()[2].Paragraphs()[0].AddRun().AddText(t.Cars[i].Number)
		row.Cells()[3].Paragraphs()[0].AddRun().AddText(t.Cars[i].FullName)
		row.Cells()[4].Paragraphs()[0].AddRun().AddText(t.Cars[i].Telephone)
		id++
	}

	if err := d.Doc.SaveToFile(d.DocPath); err != nil {
		return fmt.Errorf("failed to save edit file: %s", err.Error())
	}

	return nil
}
