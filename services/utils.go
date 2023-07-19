package services

import (
	"fmt"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Delete(chatID int64, msgID int) error {
	_, err := h.Request(tg.NewDeleteMessage(chatID, msgID))
	if err != nil {
		return fmt.Errorf("failed to send request: %s", err.Error())
	}

	return nil
}

func (h *Handler) Add(s string) {
	switch true {
	case h.data.Event == "":
		h.data.Event = s
	case h.data.How == "":
		h.data.How = s
	case h.data.Data == "":
		h.data.Data = s
		fmt.Println(s)
	case h.data.Items == nil:
		h.data.Items, h.data.CountItems = append(h.data.Items, strings.Split(s, ", ")[0]), append(h.data.Items, strings.Split(s, ", ")[1])
	}
}
