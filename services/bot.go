package services

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Start(chatID int64) error {
	msg := tg.NewMessage(chatID, "Привет! 😄\n\nСначала введи то, для чего тебе нужен рапорт, начиная со слов после \"В связи с\". Пример: редакторским просмотром фестивался творчества \"Студенческая весна\".")
	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send 'start' msg: %s", err.Error())
	}

	return nil
}

func (h *Handler) Next(chatID int64, s string) error {
	if s == "" {
		return fmt.Errorf("s can't be empty")
	}
	h.Add(s)

	var msg tg.MessageConfig
	
	switch true {
	case h.data.How == "":
		msg = tg.NewMessage(chatID, "Теперь выбери, через что ты будешь выносить предметы:")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Через КПП №1", "КПП №1"),
				tg.NewInlineKeyboardButtonData("Через гараж", "гаражный въезд"),
			),
		)
	case h.data.Date == "":
		msg = tg.NewMessage(chatID, "Теперь введи дату выноса в следующем формате: дд.мм.гггг. Пример: 31.12.2022.")
	case h.data.Time == "":
		msg = tg.NewMessage(chatID, "Теперь введи время выноса. Пример: 9:00 до 12:00.")
	case h.data.Items == nil:
		msg = tg.NewMessage(chatID, "Теперь введи предметы, которые ты собираешься выносить. Для рапорта нужны следующие параметры: наименование предмета и его количество. Пример: Стул, 2.")
	default:
		msg = tg.NewMessage(chatID, "Мы сохранили данные.")
	}

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send 'next' msg: %s", err.Error())
	}

	return nil
}

func (h *Handler) Create(chatID int64) error {
	if err := h.NewDoc(); err != nil || h.doc == nil {
		return fmt.Errorf("failed to create new document: %s", err.Error())
	}
	fmt.Println(h.doc.ReplaceMap)

	if err := h.doc.DocX.ReplaceAll(h.doc.ReplaceMap); err != nil {
		return fmt.Errorf("failed to replace: %s", err.Error())
	}

	if err := h.doc.DocX.WriteToFile("replaced.docx"); err != nil {
		return fmt.Errorf("failed to write file: %s", err.Error())
	}

	return nil
}
