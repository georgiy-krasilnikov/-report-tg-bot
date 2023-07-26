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
	case s == "/create" || s == "/get" || s == "/list" || s == "/edit" || s == "":
		return nil
	case strings.Contains(s, "docx"):
		h.doc.DocName = s
	case h.data.Event == "":
		h.data.Event = s

	case h.data.How == "":
		h.data.How = s

	case h.data.Date == "":
		h.data.Date = s

	case h.data.Time == "":
		h.data.Time = s

	case h.data.Count == 0:
		c, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("failed to add data: %s", err.Error())
		}
		h.data.Count = c

	case h.data.Items == nil:
		d := strings.Split(s, ", ")
		for i := 0; i < len(d); i += 2 {
			h.data.Items, h.data.CountItems = append(h.data.Items, d[i]), append(h.data.CountItems, d[i+1])
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
