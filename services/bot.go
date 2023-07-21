package services

import (
	"fmt"
	"os"

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
	if err := h.Add(s); err != nil {
		return fmt.Errorf("failed to add data: %s", err.Error())
	}

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
	case h.data.Count == 0:
		msg = tg.NewMessage(chatID, "Теперь введи количество видов предметов. Пример, если у нас 1 ящик, 2 стула и 1 стол: 3, если у нас 3 стула, то: 1.")
	case h.data.Items == nil:
		msg = tg.NewMessage(chatID, "Теперь введи предметы, которые ты собираешься выносить. Для рапорта нужны следующие параметры: наименование предмета и его количество. Пример: Стул, 2. Если у тебя несколько предметов, то пиши их так: Стул, 2, Стол, 1.")
	case h.data.Items != nil:
		if err := h.Create(chatID); err != nil {
			return fmt.Errorf("failed to send document: %s", err.Error())
		}
		return nil
	default:
		msg = tg.NewMessage(chatID, "Я не могу обработать эти данные.")
	}

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send 'next' msg: %s", err.Error())
	}

	return nil
}

func (h *Handler) Create(chatID int64) error {
	if err := h.CreateDocument(); err != nil {
		return fmt.Errorf("failed to create replace document: %s", err.Error())
	}

	msg := tg.NewDocument(chatID, tg.FilePath("Рапорт."+h.data.Date+".docx"))
	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send 'create' msg: %s", err.Error())
	}

	if err := os.Remove("Рапорт." + h.data.Date + ".docx"); err != nil {
		return fmt.Errorf("failed to delete document: %s", err.Error())
	}

	return nil
}
