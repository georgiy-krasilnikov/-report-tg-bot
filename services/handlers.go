package services

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	*tg.BotAPI
	data *Data
	doc  *Doc
}

func New(botToken string) (*Handler, error) {
	bot, err := tg.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %s", err.Error())
	}

	return &Handler{
		bot,
		&Data{},
		&Doc{},
	}, nil
}

func (h *Handler) Run() error {
	upd := tg.NewUpdate(0)
	upd.Timeout = 60
	h.Debug = true
	upds := h.GetUpdatesChan(upd)

	for u := range upds {
		switch true {
		case u.Message != nil && u.Message.Text == "/start":
			h.data = &Data{}
			if err := h.Start(u.Message.Chat.ID); err != nil {
				return fmt.Errorf("error in 'start' func: %s", err.Error())
			}

		case u.Message == nil && u.CallbackQuery != nil:
			if err := h.Next(u.CallbackQuery.Message.Chat.ID, u.CallbackData()); err != nil {
				return fmt.Errorf("error in 'next' func: %s", err.Error())
			}
			if err := h.DeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID); err != nil {
				return fmt.Errorf("failed to delete msg: %s", err.Error())
			}

		case u.Message != nil && u.Message.Text != "/start":
			if err := h.Next(u.Message.Chat.ID, u.Message.Text); err != nil {
				return fmt.Errorf("error in 'next' func: %s", err.Error())
			}
			if err := h.DeleteMessage(u.Message.Chat.ID, u.Message.MessageID-1); err != nil {
				return fmt.Errorf("failed to delete msg: %s", err.Error())
			}
		}
	}

	return nil
}
