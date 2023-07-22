package services

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Start(chatID int64) error {
	msg := tg.NewMessage(chatID, "Привет! Для начала выбери, что ты хочешь сделать.")
	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Создать новый рапорт", "/create"),
			tg.NewInlineKeyboardButtonData("Редактировать ранее созданный рапорт", "/edit"),
		),
	)

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send 'start' msg: %s", err.Error())
	}

	return nil
}

func (h *Handler) Next(chatID int64, s string) error {
	switch true {
	case s == "":
		return fmt.Errorf("s can't be empty")
	case s == "/create":
		h.data.IsNew = true
	case s == "/edit":
		h.data.IsNew = false
	default:
		if err := h.AddData(s); err != nil {
			return fmt.Errorf("failed to add data: %s", err.Error())
		}
	}

	var msg tg.MessageConfig

	switch true {
	case !h.data.IsNew:
		lst, err := GetListOfDocuments()
		if err != nil {
			return fmt.Errorf("failed to get list of documents: %s", err.Error())
		}
		msg = tg.NewMessage(chatID, "Теперь выбери рапорт, в который ты хочешь внести изменения:")
		msg.ReplyMarkup = CreateKeyboard(lst)
		h.data.IsNew = true
	case h.data.Event == "":
		msg = tg.NewMessage(chatID, "Сначала введи мероприятие, для которого тебе нужен рапорт, начиная со слов после _В связи с_. *Пример:* _редакторским просмотром фестивался творчества \"Студенческая весна\"_.")
	case h.data.How == "":
		msg = tg.NewMessage(chatID, "Теперь выбери, каким образом ты будешь выносить предметы:")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Через КПП №1", "КПП №1"),
				tg.NewInlineKeyboardButtonData("Через гараж", "гаражный въезд"),
			),
		)
	case h.data.Date == "":
		msg = tg.NewMessage(chatID, "Теперь введи дату выноса в следующем формате: _дд.мм.гггг_. *Пример:* _31.12.2022_.")
	case h.data.Time == "":
		msg = tg.NewMessage(chatID, "Теперь введи время выноса. *Пример:* _9:00 до 12:00_.")
	case h.data.Count == 0:
		msg = tg.NewMessage(chatID, "Теперь введи количество видов предметов. *Пример:* если у нас 1 ящик, 2 стула и 1 стол: _3_, если у нас 3 стула, то: _1_.")
	case h.data.Items == nil:
		msg = tg.NewMessage(chatID, "Теперь введи предметы, которые ты собираешься выносить. Для рапорта нужны следующие параметры: наименование предмета и его количество. *Пример:* _Стул, 2_. Если у тебя *несколько предметов*, то пиши их так: _Стул, 2, Стол, 1_.")
	case h.data.Items != nil:
		if err := h.SendDocument(chatID); err != nil {
			return fmt.Errorf("failed to send document: %s", err.Error())
		}
		return nil
	default:
		msg = tg.NewMessage(chatID, "Я не могу обработать эти данные.")
	}
	msg.ParseMode = "markdown"

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send 'next' msg: %s", err.Error())
	}

	return nil
}

func (h *Handler) SendDocument(chatID int64) error {
	if err := h.CreateDocument(); err != nil {
		return fmt.Errorf("failed to create replace document: %s", err.Error())
	}

	msg := tg.NewDocument(chatID, tg.FilePath("docs/Рапорт."+h.data.Date+".docx"))
	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send 'create' msg: %s", err.Error())
	}
	if err := h.DeleteDocument(); err != nil {
		return fmt.Errorf("failed to check time: %s", err.Error())
	}

	return nil
}
