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

func (h *Handler) NewItems(s []string) {
	for i := 0; i < h.data.Table.ItemsNumber; i++ {
		h.data.Table.Items = append(h.data.Table.Items, Item{strings.Split(s[i], ", ")[0], strings.Split(s[i], ", ")[1]})
	}
}

func (h *Handler) NewCars(s []string) {
	for i := 0; i < h.data.Table.CarsNumber; i++ {
		h.data.Table.Cars = append(h.data.Table.Cars, Car{strings.Split(s[i], ", ")[0], strings.Split(s[i], ", ")[1], strings.Split(s[i], ", ")[2], strings.Split(s[i], ", ")[3]})
	}
}

func (h *Handler) AddData(s string) error {
	switch true {
	case strings.Contains(s, "/") || s == "":
		return nil
		
	case strings.Contains(s, "docx"):
		h.doc.DocName = s
		h.doc.DocPath = "docs/" + s
		
	case h.data.Event == "" && !isNumber(s):
		h.data.Event = s

	case h.data.How == "" && h.data.Event != "" && !isNumber(s):
		h.data.How = s

	case (h.data.Date == "" && h.data.How != "") || strings.Contains(s, ".20"):
		h.data.Date = s

	case h.data.Time == "" && h.data.Date != "" && strings.Contains(s, ":"):
		h.data.Time = s

	case h.data.Table.ItemsNumber == 0 && isNumber(s):
		n, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("failed to add ItemsNumber: %s", err.Error())
		}

		h.data.Table.ItemsNumber = n

	case h.data.Table.Items == nil && h.data.Table.ItemsNumber != 0:
		h.NewItems(strings.Split(s, " | "))

	case h.data.Table.Items != nil && isNumber(s):
		n, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("failed to add CarsNumber: %s", err.Error())
		}
		h.data.Table.CarsNumber = n

	case h.data.Table.Cars == nil && h.data.Table.CarsNumber != 0:
		h.NewCars(strings.Split(s, " | "))
	}

	return nil
}

func newKeyboard(lst []string, data []string) tg.InlineKeyboardMarkup {
	var btns []tg.InlineKeyboardButton
	for _, v := range lst {
		btns = append(btns, tg.NewInlineKeyboardButtonData(v, v))
	}

	return tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(btns...))
}

func isNumber(s string) bool {
	if strings.HasPrefix(s, "0") {
		return false
	}

	for _, v := range s {
		if v < '0' || v > '9' {
			return false
		}
	}

	return true
}
