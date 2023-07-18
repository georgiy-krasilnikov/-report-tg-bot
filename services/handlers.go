package services

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	*tg.BotAPI
	data *Data
}

func New(botToken string) (*Handler, error) {
	bot, err := tg.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %s", err.Error())
	}

	return &Handler{
		bot,
		&Data{},
	}, nil
}

func (h *Handler) Run() error {
	upd := tg.NewUpdate(0)
	upd.Timeout = 60

	upds := h.GetUpdatesChan(upd)

	for u := range upds {
		if u.Message == nil && u.CallbackQuery != nil {
			fmt.Print(u.CallbackData())
			h.Delete(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID)
			if err := h.Next(u.CallbackQuery.Message.Chat.ID); err != nil {
				return fmt.Errorf("failed to call func 'next': %s", err.Error())
			}
		}

		if u.Message != nil {
			switch u.Message.Text {
			case "/start":
				if err := h.Start(u.Message.Chat.ID); err != nil {
					return fmt.Errorf("failed to call func 'next': %s", err.Error())
				}
			}
		}
	}

	return nil
}
