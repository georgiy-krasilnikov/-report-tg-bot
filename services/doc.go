package services

import (
	"fmt"

	"github.com/lukasjarosch/go-docx"
)

func (h *Handler) NewDoc() error {
	doc, err := docx.Open("file.docx")
	if err != nil {
		return fmt.Errorf("failed to open doc: %s", err.Error())
	}

	h.doc = &Doc{
		DocX: doc,
		ReplaceMap: docx.PlaceholderMap{
			"дд.мм.гггг": h.data.Date,
			"xxx":        h.data.Event,
			"yyy":        h.data.How,
			"zzz":        h.data.Time,
			"x1":         h.data.Items[0],
			"y1":         h.data.CountItems[0]},
	}

	return nil
}
