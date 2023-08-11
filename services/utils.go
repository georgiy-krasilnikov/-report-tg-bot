package services

import (
	"fmt"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) DeleteMessage(chatID int64, msgID int) error {
	_, err := h.Request(tg.NewDeleteMessage(chatID, msgID))
	if err != nil {
		return fmt.Errorf("failed to send request: %s", err.Error())
	}

	return nil
}

func (h *Handler) NewItems(s []string) {
	for i := 0; i < len(s); i++ {
		h.data.Table.Items = append(h.data.Table.Items, Item{strings.Split(s[i], ", ")[0], strings.Split(s[i], ", ")[1]})
	}
}

func (h *Handler) NewCars(s []string) {
	for i := 0; i < len(s); i++ {
		h.data.Table.Cars = append(h.data.Table.Cars, Car{strings.Split(s[i], ", ")[0], strings.Split(s[i], ", ")[1], strings.Split(s[i], ", ")[2], strings.Split(s[i], ", ")[3]})
	}
}

func (h *Handler) AddData(s string) error {
	if h.mood == "/create" {
		switch true {
		case s == "":
			return nil

		case h.data.Event == "" && isDate(s) != "" && isTime(s) != "":
			h.data.Event = s

		case h.data.How == "" && h.data.Event != "" && isDate(s) != "" && isTime(s) != "":
			h.data.How = s

		case h.data.Date == "" && h.data.How != "" && isDate(s) == "" && isTime(s) != "":
			h.data.Date = s

		case h.data.Time == "" && h.data.Date != "" && isDate(s) != "" && isTime(s) == "":
			h.data.Time = s

		case h.data.Table.Items == nil && h.data.Time != "" && isDate(s) != "" && isTime(s) != "":
			h.NewItems(strings.Split(s, " | "))

		case h.data.Table.Cars == nil && h.data.Table.Items != nil && isDate(s) != "" && isTime(s) != "":
			h.NewCars(strings.Split(s, " | "))
		}
	} else if h.mood == "/list" {
		switch true {
		case strings.HasPrefix(s, "/") || s == "" || strings.HasPrefix(s, "id: "):
			return nil

		case strings.HasSuffix(s, ".docx"):
			if err := h.NewDoc(s, "docs/"+s); err != nil {
				return fmt.Errorf("failed to assign doc to handler: %s", err.Error())
			}

		case h.data.Date == "" && isDate(s) == "" && isTime(s) != "":
			h.data.Date = s

		case h.data.Table.Items == nil && isDate(s) != "" && isTime(s) != "":
			h.NewItems(strings.Split(s, " | "))

		case h.data.Table.Cars == nil && isDate(s) != "" && isTime(s) != "":
			h.NewCars(strings.Split(s, " | "))
		}
	}

	return nil
}

func newKeyboard(lst, data []string) tg.InlineKeyboardMarkup {
	var kbrd [][]tg.InlineKeyboardButton
	size := 2

	for i := 0; i < len(lst); i += 2 {
		if len(lst)-i == 1 || len(lst[i]) > 36 {
			size = 1
		}

		var btns []tg.InlineKeyboardButton
		for j := i; len(btns) < size; j++ {
			btns = append(btns, tg.NewInlineKeyboardButtonData(lst[j], data[j]))
		}

		kbrd = append(kbrd, btns)
	}

	return tg.InlineKeyboardMarkup{InlineKeyboard: kbrd}
}

func isDate(s string) string {
	_, err := time.Parse("02.01.2006", s)
	if err != nil {
		return fmt.Sprintf("invalid date format: %s", err.Error())
	}

	return ""
}

func isTime(s string) string {
	_, err := time.Parse("15:04", s)
	if err != nil {
		return fmt.Sprintf("invalid time format: %s", err.Error())
	}

	return ""
}
