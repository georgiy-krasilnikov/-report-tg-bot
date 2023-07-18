package services

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Start(chatID int64) error {
	msg := tg.NewMessage(chatID, "Привет! 😄\n\nВыбери мероприятие, для которого тебе нужен рапорт:")
	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Дебют", "Дебют"),
			tg.NewInlineKeyboardButtonData("Неон", "Неон"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Студенческая весна", "Студенческая весна"),
			tg.NewInlineKeyboardButtonData("PERFOMANCE", "PERFOMANCE"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Мисс и мистер университет", "Мисс и мистер университет"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Другое мероприятие", "Другое"),
		),
	)

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send 'start' msg: %s", err.Error())
	}

	return nil
}

func (h *Handler) Next(chatID int64) error {
	msg := tg.NewMessage(chatID, "Привет! 😄\n\nВыбери мероприятие, для которого тебе нужен рапорт:")
	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Через КПП №1", "КПП №1"),
			tg.NewInlineKeyboardButtonData("Через гараж", "гаражный въезд"),
		),
	)

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send 'next' msg: %s", err.Error())
	}

	return nil
}
