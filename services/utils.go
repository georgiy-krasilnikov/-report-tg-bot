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
	if mode == "/create" {
		switch true {
		case s == "" || strings.HasPrefix(s, "/"):
			return nil

		case h.data.Event == "" && isDate(s) != "":
			h.data.Event = s
			if class == "/car-raport" || class == "/full-raport" {
				h.data.How = "гаражный въезд"
			}

		case h.data.How == "" && h.data.Event != "" && isDate(s) != "":
			h.data.How = s

		case h.data.Date == "" && h.data.How != "" && isDate(s) == "":
			h.data.Date = s

		case h.data.Time == "" && h.data.Date != "" && isDate(s) != "":
			h.data.Time = s

		case h.data.Table.Items == nil && h.data.Time != "" && isDate(s) != "" && strings.Count(strings.Split(s, " | ")[0], ", ") == 1:
			h.NewItems(strings.Split(s, " | "))

		case h.data.Table.Cars == nil && isDate(s) != "" && strings.Count(strings.Split(s, " | ")[0], ", ") == 3:
			h.NewCars(strings.Split(s, " | "))
		}
	} else if mode == "/list" {
		switch true {
		case strings.HasPrefix(s, "/") || s == "" || strings.Contains(s, "id: "):
			return nil

		case strings.HasSuffix(s, ".docx"):
			if err := h.NewDoc(s, "docs/"+s); err != nil {
				return fmt.Errorf("failed to assign doc to handler: %s", err.Error())
			}

		case h.data.Date == "" && isDate(s) == "":
			h.data.Date = s

		case h.data.Table.Items == nil && isDate(s) != "" && strings.Count(strings.Split(s, " | ")[0], ", ") == 1:
			h.NewItems(strings.Split(s, " | "))

		case h.data.Table.Cars == nil && isDate(s) != "" && strings.Count(strings.Split(s, " | ")[0], ", ") == 3:
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
