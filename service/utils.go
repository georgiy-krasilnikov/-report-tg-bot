package service

import (
	"fmt"
	"strings"

	"report-bot/doc"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewItems(s []string) []doc.Item {
	var items []doc.Item
	for i := 0; i < len(s); i++ {
		items = append(items, doc.Item{
			Name:  strings.Split(s[i], ", ")[0],
			Count: strings.Split(s[i], ", ")[1],
		})
	}

	return items
}

func NewCars(s []string) []doc.Car {
	var cars []doc.Car
	for i := 0; i < len(s); i++ {
		cars = append(cars, doc.Car{
			Brand:     strings.Split(s[i], ", ")[0],
			Number:    strings.Split(s[i], ", ")[1],
			FullName:  strings.Split(s[i], ", ")[2],
			Telephone: strings.Split(s[i], ", ")[3],
		})
	}

	return cars
}

func (h *Handler) AddData(s string) error {
	if Mode == "/create" {
		switch true {
		case s == "" || strings.HasPrefix(s, "/"):
			return nil

		case h.data.Event == "":
			h.data.Event = s
			if Class == "/car-raport" || Class == "/full-raport" {
				h.data.How = "гаражный въезд"
			}

		case h.data.How == "" && h.data.Event != "":
			h.data.How = s

		case h.data.Date == "" && h.data.How != "":
			h.data.Date = s

		case h.data.Time == "" && h.data.Date != "":
			h.data.Time = s

		case h.data.Table.Items == nil && h.data.Time != "" && strings.Count(strings.Split(s, " | ")[0], ", ") == 1:
			h.data.Table.Items = NewItems(strings.Split(s, " | "))

		case h.data.Table.Cars == nil && strings.Count(strings.Split(s, " | ")[0], ", ") == 3:
			h.data.Table.Cars = NewCars(strings.Split(s, " | "))
		}
	} else if Mode == "/list" {
		switch true {
		case strings.HasPrefix(s, "/") || s == "" || strings.Contains(s, "id: "):
			return nil

		case strings.HasSuffix(s, ".docx"):
			doc, err := doc.NewDoc(s, "docs/"+s)
			if err != nil {
				return fmt.Errorf("failed to assign doc to handler: %s", err.Error())
			}

			h.doc = doc

		case h.data.Date == "" && strings.Contains(s, "202"):
			h.data.Date = s

		case h.data.Table.Items == nil && strings.Count(strings.Split(s, " | ")[0], ", ") == 1:
			h.data.Table.Items = NewItems(strings.Split(s, " | "))

		case h.data.Table.Cars == nil && strings.Count(strings.Split(s, " | ")[0], ", ") == 3:
			h.data.Table.Cars = NewCars(strings.Split(s, " | "))
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
