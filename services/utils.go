package services

import (
	"fmt"
	"os"
	"strconv"
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

func (h *Handler) AddData(s string) error {
	switch true {
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

func (h *Handler) DeleteDocument() error {
	if time.Now().Format("01.02.2006") == h.data.Date && strconv.Itoa(time.Now().Hour()) == "23" && strconv.Itoa(time.Now().Minute()) == "59" {
		if err := os.Remove("docs/Рапорт." + h.data.Date + ".docx"); err != nil {
			return fmt.Errorf("failed to delete document: %s", err.Error())
		}
	}

	return nil
}

func CreateKeyboard(lst []string) tg.InlineKeyboardMarkup {
	key := tg.NewInlineKeyboardRow()
	for _, v := range lst {
		key = append(key, tg.NewInlineKeyboardButtonData(v, v))
	}

	return tg.NewInlineKeyboardMarkup(key)
}
