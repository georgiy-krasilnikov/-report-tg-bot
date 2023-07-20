package services

import (
	"fmt"
	"strconv"

	"github.com/lukasjarosch/go-docx"
)

func (h *Handler) NewDoc() error {
	doc, err := docx.Open("file.docx")
	if err != nil {
		return fmt.Errorf("failed to open doc: %s", err.Error())
	}

	repMap := docx.PlaceholderMap{
		"дд.мм.гггг": h.data.Date,
		"xxx":        h.data.Event,
		"yyy":        h.data.How,
		"zzz":        h.data.Time}

	for i := 0; i < h.data.Count; i++ {
		repMap["x"+strconv.Itoa(i)] = h.data.Items[i]
		repMap["y"+strconv.Itoa(i)] = h.data.CountItems[i]
	}

	h.doc = &Doc{
		DocX:       doc,
		ReplaceMap: repMap,
	}

	return nil
}

func (h *Handler) CreateDocument() error {
	if err := h.NewDoc(); err != nil || h.doc == nil {
		return fmt.Errorf("failed to create new doc: %s", err.Error())
	}

	



	// if err := h.doc.DocX.ReplaceAll(h.doc.ReplaceMap); err != nil {
	// 	return fmt.Errorf("failed to replace: %s", err.Error())
	// }

	// if err := h.doc.DocX.WriteToFile("Рапорт." + h.data.Date + ".docx"); err != nil {
	// 	return fmt.Errorf("failed to write file: %s", err.Error())
	// }

	return nil
}
