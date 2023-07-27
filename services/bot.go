package services

import (
	"fmt"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Start(chatID int64) error {
	msg := tg.NewMessage(chatID, "Привет! Для начала выбери, что ты хочешь сделать.")
	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Создать новый рапорт", "/create"),
			tg.NewInlineKeyboardButtonData("Выбрать рапорт из списка", "/list"),
		),
	)

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send 'start' msg: %s", err.Error())
	}

	return nil
}

func (h *Handler) Next(chatID int64, s string) error {
	if err := h.AddData(s); err != nil {
		return fmt.Errorf("failed to add data: %s", err.Error())
	}

	var msg tg.MessageConfig

	switch true {
	case s == "":
		return fmt.Errorf("s can't be empty")
	case s == "/list":
		msg = tg.NewMessage(chatID, "Теперь выбери рапорт:")
		kbrd, err := NewKeyboard()
		if err != nil {
			return fmt.Errorf("failed to create keyboard: %s", err.Error())
		}
		msg.ReplyMarkup = kbrd

	case strings.Contains(s, "docx"):
		msg = tg.NewMessage(chatID, "Теперь выбери что ты хочешь сделать с выбранным рапортом:")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Получить рапорт", "/get"),
				tg.NewInlineKeyboardButtonData("Редактировать рапорт", "/edit"),
			),
		)

	case s == "/get":
		msg := tg.NewDocument(chatID, tg.FilePath("docs/"+h.doc.DocName))
		msg.Caption = "Вот твой рапорт 👇"
		if _, err := h.Send(msg); err != nil {
			return fmt.Errorf("failed to send msg with document: %s", err.Error())
		}

	case s == "/edit":
		if err := h.SendEditMessage(chatID); err != nil {
			return fmt.Errorf("error in 'edit' func: %s", err.Error())
		}
		return nil

	case s == "/create":
		msg = tg.NewMessage(chatID, "Сначала введи мероприятие, для которого тебе нужен рапорт, начиная со слов после _В связи с_. *Пример:* _редакторским просмотром фестивался творчества \"Студенческая весна\"_.")

	case h.data.How == "" && h.data.Event != "":
		msg = tg.NewMessage(chatID, "Теперь выбери, каким образом ты будешь выносить предметы:")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Через КПП №1", "КПП №1"),
				tg.NewInlineKeyboardButtonData("Через гараж", "гаражный въезд"),
			),
		)

	case s == "/items":
		msg = tg.NewMessage(chatID, "Теперь выбери, что ты хочешь сделать со списком предметов:")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Удалить все предметы", "/delete"),
				tg.NewInlineKeyboardButtonData("Добавить предметы", "/add"),
			),
		)
	case s == "/add":

	case s == "/delete":

	case (h.data.Date == "" && h.data.How != "") || s == "/data":
		msg = tg.NewMessage(chatID, "Теперь введи дату выноса в следующем формате: _дд.мм.гггг_. *Пример:* _31.12.2022_.")

	case h.data.Time == "" && h.data.Date != "":
		msg = tg.NewMessage(chatID, "Теперь введи время выноса. *Пример:* _9:00 до 12:00_.")

	case (h.data.Count == 0 && h.data.Time != "") || s == "/add":
		msg = tg.NewMessage(chatID, "Теперь введи количество видов предметов. *Пример:* если у нас 1 ящик, 2 стула и 1 стол: _3_, если у нас 3 стула, то: _1_.")

	case h.data.Items == nil && h.data.Count != 0:
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

	msg := tg.NewDocument(chatID, tg.FilePath("docs/"+h.doc.DocName))
	msg.Caption = "Вот твой рапорт 👇"
	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send msg with document: %s", err.Error())
	}

	return nil
}

func (h *Handler) SendEditMessage(chatID int64) error {
	msg := tg.NewDocument(chatID, tg.FilePath("docs/"+h.doc.DocName))
	msg.Caption = "Вот какой твой рапорт выглядит сейчас. Теперь выбери, что ты хочешь редактировать:"
	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Дату", "/data"),
			tg.NewInlineKeyboardButtonData("Список предметов", "/items"),
		),
	)

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send msg with document: %s", err.Error())
	}

	return nil
}
