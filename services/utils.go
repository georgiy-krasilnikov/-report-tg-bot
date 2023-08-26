package services

import (
	"fmt"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

		case h.data.Event == "":
			h.data.Event = s
			if class == "/car-raport" || class == "/full-raport" {
				h.data.How = "гаражный въезд"
			}

		case h.data.How == "" && h.data.Event != "":
			h.data.How = s

		case h.data.Date == "" && h.data.How != "":
			h.data.Date = s

		case h.data.Time == "" && h.data.Date != "":
			h.data.Time = s

		case h.data.Table.Items == nil && h.data.Time != "" && strings.Count(strings.Split(s, " | ")[0], ", ") == 1:
			h.NewItems(strings.Split(s, " | "))

		case h.data.Table.Cars == nil && strings.Count(strings.Split(s, " | ")[0], ", ") == 3:
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

		case h.data.Date == "":
			h.data.Date = s

		case h.data.Table.Items == nil && strings.Count(strings.Split(s, " | ")[0], ", ") == 1:
			h.NewItems(strings.Split(s, " | "))

		case h.data.Table.Cars == nil && strings.Count(strings.Split(s, " | ")[0], ", ") == 3:
			h.NewCars(strings.Split(s, " | "))
		}
	}

	return nil
}

func newKeyboard(lst, data []string) tg.InlineKeyboardMarkup {
	var kbrd [][]tg.InlineKeyboardButton

	if strings.Count(lst[0], "|") == 1 || !strings.Contains(lst[0], "|") {
		size := 2
		for i := 0; i < len(lst); i += 2 {
			if len(lst)-i == 1 {
				size = 1
			}

			var btns []tg.InlineKeyboardButton
			for j := i; len(btns) < size; j++ {
				btns = append(btns, tg.NewInlineKeyboardButtonData(lst[j], data[j]))
			}

			kbrd = append(kbrd, btns)
		}
	} else if strings.Count(lst[0], "|") == 3 {
		for i := 0; i < len(lst); i++ {
			var btns = []tg.InlineKeyboardButton{tg.NewInlineKeyboardButtonData(lst[i], data[i])}
			kbrd = append(kbrd, btns)
		}
	}

	return tg.InlineKeyboardMarkup{InlineKeyboard: kbrd}
}
