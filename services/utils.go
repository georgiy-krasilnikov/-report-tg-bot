package services

import (
	"fmt"
	"strconv"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) DeleteMessage(chatID int64, msgID int) error {
	_, err := h.Request(tg.NewDeleteMessage(chatID, msgID))
	if err != nil {
		return fmt.Errorf("failed to send request: %s", err.Error())
	}

	return nil
}

func (h *Handler) AddData(s string) error {
	switch true {
	case strings.Contains(s, "/") || s == "":
		return nil
	case strings.Contains(s, "docx"):
		h.doc.DocName = s
		h.doc.DocPath = "docs/" + s
	case h.data.Event == "" && strings.Contains(s, "редакторским"):
		h.data.Event = s

	case h.data.How == "" && h.data.Event != "":
		h.data.How = s

	case (h.data.Date == "" && h.data.How != "") || strings.Contains(s, ".20"):
		h.data.Date = s

	case h.data.Time == "" && h.data.Date != "" && strings.Contains(s, ":"):
		h.data.Time = s

	case h.data.Table.ItemsNumber == 0:
		n, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("failed to add data: %s", err.Error())
		}
		h.data.Table.ItemsNumber = n

	case h.data.Table.Items == nil && h.data.Table.ItemsNumber != 0:
		d := strings.Split(s, ", ")
		for i := 0; i < len(d); i += 2 {
			h.data.Table.Items, h.data.Table.CountItems = append(h.data.Table.Items, d[i]), append(h.data.Table.CountItems, d[i+1])
		}
	}

	return nil
}

func NewKeyboard() (tg.InlineKeyboardMarkup, error) {
	lst, err := GetListOfDocuments()
	if err != nil {
		return tg.InlineKeyboardMarkup{}, fmt.Errorf("failed to get list of documents: %s", err.Error())
	}

	var btns []tg.InlineKeyboardButton
	for _, v := range lst {
		btns = append(btns, tg.NewInlineKeyboardButtonData(v, v))
	}

	return tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(btns...)), nil
}
