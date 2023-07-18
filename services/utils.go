package services

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Delete(chatID int64, msgID int) error {
	_, err := h.Send(tg.NewDeleteMessage(chatID, msgID))
	if err != nil {
		return fmt.Errorf("failed to delete msg: %s", err.Error())
	}

	return nil
}
