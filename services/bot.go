package services

import (
	"fmt"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	id    string
	class string
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
	if mode == "/create" {
		if err := h.CreateBranch(chatID, s); err != nil {
			return fmt.Errorf("error in 'Create' branch: %s", err.Error())
		}
	} else {
		if err := h.ListBranch(chatID, s); err != nil {
			return fmt.Errorf("error in 'List' branch: %s", err.Error())
		}
	}

	return nil
}

func (h *Handler) CreateBranch(chatID int64, s string) error {
	if err := h.AddData(s); err != nil {
		return fmt.Errorf("failed to add data in 'Create' branch: %s", err.Error())
	}

	var msg tg.MessageConfig

	switch true {
	case s == "":
		return fmt.Errorf("s can't be empty")

	case s == "/create":
		msg = tg.NewMessage(chatID, "Сначала выбери, какой вид рапорта тебе нужен:")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Рапорт на вынос предметов", "/item-raport"),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Рапорт на въезд-выезд", "/car-raport"),
			),
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Рапорт и на вынос предметов, и на въезд-выезд", "/full-raport"),
			),
		)

	case s == "/car-raport" || s == "/item-raport" || s == "/full-raport":
		class = s
		msg = tg.NewMessage(chatID, "Сначала введи мероприятие, для которого тебе нужен рапорт, начиная со слов после _В связи с_. *Пример:* _редакторским просмотром фестиваля творчества \"Студенческая весна\"_.")

	case h.data.How == "" && h.data.Event != "":
		msg = tg.NewMessage(chatID, "Теперь выбери, каким образом ты будешь выносить предметы:")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Через КПП №1", "КПП №1"),
				tg.NewInlineKeyboardButtonData("Через гараж", "гаражный въезд"),
			),
		)

	case h.data.Date == "" && h.data.How != "":
		msg = tg.NewMessage(chatID, "Теперь введи дату в следующем формате: _дд.мм.гггг_. *Пример:* _31.12.2022_.")

	case h.data.Time == "" && h.data.Date != "" && h.data.Event != "":
		msg = tg.NewMessage(chatID, "Теперь введи время. *Пример:* _9:00 до 12:00_.")

	case (h.data.Table.Items == nil && h.data.Time != "" && class == "/item-raport") || (h.data.Table.Items == nil && h.data.Time != "" && class == "/full-raport"):
		msg = tg.NewMessage(chatID, "Теперь введи предметы, которые ты собираешься добавить в рапорт. Для рапорта нужны следующие параметры: наименование предмета и его количество. *Пример:* _Стул, 2_. Если у тебя *несколько предметов*, то пиши их так: _Стул, 2 | Стол, 1_.")

	case (h.data.Table.Cars == nil && class == "/car-raport") || (h.data.Table.Cars == nil && class == "/full-raport"):
		msg = tg.NewMessage(chatID, "Теперь введи данные автомобилей, которые ты собираешься добавить. Для рапорта нужны следующие параметры: марка автомобиля, его госномер, его ФИО, и его номер телефона. *Пример:* _Volkswagen Polo, А000ВС77, Иванов Иван Иванович, +78005553535_. Если у тебя *несколько автомобилей*, то пиши их так: _Volkswagen Polo, А000ВС77, Иванов Иван Иванович, +78005553535 | Kia Rio, А111ВС77, Александров Александр Александрович, +78005554545_.")

	case (class == "/item-raport" && h.data.Table.Items != nil && h.data.Table.Cars == nil) || (class == "/car-raport" && h.data.Table.Cars != nil && h.data.Table.Items == nil) || (class == "/full-raport" && h.data.Table.Cars != nil && h.data.Table.Items != nil):
		if err := h.CreateDocument(); err != nil {
			return fmt.Errorf("failed to create document: %s", err.Error())
		}

		if err := h.SendDocument(chatID); err != nil {
			return fmt.Errorf("failed to send document: %s", err.Error())
		}

		return nil

	default:
		msg = tg.NewMessage(chatID, "Я не могу обработать эти данные.")
	}

	msg.ParseMode = "markdown"

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send msg: %s", err.Error())
	}

	return nil
}

func (h *Handler) ListBranch(chatID int64, s string) error {
	if err := h.AddData(s); err != nil {
		return fmt.Errorf("failed to add data in 'List' branch: %s", err.Error())
	}

	var msg tg.MessageConfig

	switch true {
	case s == "":
		return fmt.Errorf("s can't be empty")

	case s == "/list":
		msg = tg.NewMessage(chatID, "Теперь выбери рапорт:")
		docs, err := GetListOfDocuments()
		if err != nil {
			return fmt.Errorf("failed to get list of documents: %s", err.Error())
		}

		msg.ReplyMarkup = newKeyboard(docs, docs)

	case strings.HasSuffix(s, ".docx"):
		msg = tg.NewMessage(chatID, "Теперь выбери что ты хочешь сделать с выбранным рапортом:")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Получить рапорт", "/get"),
				tg.NewInlineKeyboardButtonData("Редактировать рапорт", "/edit"),
			),
		)

	case s == "/get":
		if err := h.SendDocument(chatID); err != nil {
			return fmt.Errorf("failed to send document: %s", err.Error())
		}

		return nil

	case s == "/edit":
		if err := h.SendEditMessage(chatID); err != nil {
			return fmt.Errorf("failed to send 'edit' msg: %s", err.Error())
		}

		return nil

	case s == "/date":
		msg = tg.NewMessage(chatID, "Теперь введи новую дату в следующем формате: _дд.мм.гггг_. *Пример:* _31.12.2022_.")

	case isDate(s) == "":
		if err := h.EditDate(); err != nil {
			return fmt.Errorf("failed to edit date in document: %s", err.Error())
		}

		if err := h.SendDocument(chatID); err != nil {
			return fmt.Errorf("failed to send document: %s", err.Error())
		}

		return nil

	case s == "/items":
		msg = tg.NewMessage(chatID, "Теперь выбери, что ты хочешь сделать со списком предметов:")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Изменить данные", "/edit item"),
				tg.NewInlineKeyboardButtonData("Добавить предметы(-ы)", "/add item"),
			),
		)

	case s == "/cars":
		msg = tg.NewMessage(chatID, "Теперь выбери, что ты хочешь сделать со списком автомобилей:")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Изменить данные", "/edit car"),
				tg.NewInlineKeyboardButtonData("Добавить автомобиль(-и)", "/add car"),
			),
		)

	case s == "/edit item":
		if err := h.SendItemMessage(chatID); err != nil {
			return fmt.Errorf("failed to send item table: %s", err.Error())
		}

		return nil

	case s == "/edit car":
		if err := h.SendAutoMessage(chatID); err != nil {
			return fmt.Errorf("failed to send auto table: %s", err.Error())
		}

		return nil

	case strings.HasPrefix(s, "item id: "):
		id = strings.TrimPrefix(s, "item id: ")
		msg = tg.NewMessage(chatID, "Теперь введи новые данные. Для рапорта нужны следующие параметры: наименование предмета и его количество. *Пример:* _Стул, 2_.")

	case strings.HasPrefix(s, "car id: "):
		id = strings.TrimPrefix(s, "car id: ")
		msg = tg.NewMessage(chatID, "Теперь введи новые данные. Для рапорта нужны следующие параметры: марка автомобиля, его госномер, ФИО владельца, и его номер телефона. *Пример:* _Volkswagen Polo, А000ВС77, Иванов Иван Иванович, +7800553535_. Если у тебя *несколько автомобилей*, то пиши их так: _Volkswagen Polo, А000ВС77, Иванов Иван Иванович, +78005553535 | Kia Rio, А111ВС77, Александров Александр Александрович, +78005554545_.")

	case s == "/add item":
		msg = tg.NewMessage(chatID, "Теперь введи предметы, которые ты собираешься добавить. Для рапорта нужны следующие параметры: наименование предмета и его количество. *Пример:* _Стул, 2_. Если у тебя *несколько предметов*, то пиши их так: _Стул, 2 | Стол, 1_.")

	case s == "/add car":
		msg = tg.NewMessage(chatID, "Теперь введи данные автомобилей, которые ты собираешься добавить. Для рапорта нужны следующие параметры: марка автомобиля, его госномер, его ФИО, и его номер телефона. *Пример:* _Volkswagen Polo-А000ВС77-Иванов Иван Иванович-+78005553535_. Если у тебя *несколько автомобилей*, то пиши их так: _Volkswagen Polo, А000ВС77, Иванов Иван Иванович, +78005553535 | Kia Rio, А111ВС77, Александров Александр Александрович, +78005554545_.")

	case h.data.Table.Cars != nil && id != "":
		if err := h.EditCarRow(id); err != nil {
			return fmt.Errorf("failed to edit car row in document: %s", err.Error())
		}

		if err := h.SendDocument(chatID); err != nil {
			return fmt.Errorf("failed to send document: %s", err.Error())
		}

		return nil

	case h.data.Table.Items != nil && id != "":

		if err := h.EditItemRow(id); err != nil {
			return fmt.Errorf("failed to edit item row in document: %s", err.Error())
		}

		if err := h.SendDocument(chatID); err != nil {
			return fmt.Errorf("failed to send document: %s", err.Error())
		}

		return nil

	case h.data.Table.Items != nil && id == "":
		if err := h.AddItemRow(); err != nil {
			return fmt.Errorf("failed to add item row in document: %s", err.Error())
		}

		if err := h.SendDocument(chatID); err != nil {
			return fmt.Errorf("failed to send document: %s", err.Error())
		}

		return nil

	case h.data.Table.Cars != nil && id == "":
		if err := h.AddCarRow(); err != nil {
			return fmt.Errorf("failed to add car row in document: %s", err.Error())
		}

		if err := h.SendDocument(chatID); err != nil {
			return fmt.Errorf("failed to send document: %s", err.Error())
		}

		return nil

	default:
		msg = tg.NewMessage(chatID, "Я не могу обработать эти данные.")
	}

	msg.ParseMode = "markdown"

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send 'List' msg: %s", err.Error())
	}

	return nil
}

func (h *Handler) SendEditMessage(chatID int64) error {
	msg := tg.NewDocument(chatID, tg.FilePath(h.doc.DocPath))
	msg.Caption = "Вот как твой рапорт выглядит сейчас. Теперь выбери, что ты хочешь редактировать:"
	msg.ReplyMarkup = newKeyboard([]string{"Дата", "Предметы", "Автомобили"}, []string{"/date", "/items", "/cars"})

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send msg with document: %s", err.Error())
	}

	return nil
}

func (h *Handler) SendItemMessage(chatID int64) error {
	items, err := h.GetListOfItems()
	if err != nil {
		return fmt.Errorf("failed to get rows from table: %s", err.Error())
	}

	var row []string
	var IDs []string

	for _, it := range items {
		row = append(row, fmt.Sprintf("%s | %s", it[1], it[2]))
		IDs = append(IDs, fmt.Sprintf("item id: %s", it[0]))
	}

	msg := tg.NewMessage(chatID, "Теперь выбери, строчку, которую хочешь изменить:")
	msg.ReplyMarkup = newKeyboard(row, IDs)

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send msg with items: %s", err.Error())
	}

	return nil
}

func (h *Handler) SendAutoMessage(chatID int64) error {
	cars, err := h.GetListOfCars()
	if err != nil {
		return fmt.Errorf("failed to get list of cars: %s", err.Error())
	} else if err == nil && cars == nil {
		if _, err := h.Send(tg.NewMessage(chatID, "В этом рапорте не указаны автомобили.")); err != nil {
			return fmt.Errorf("failed to send msg without cars: %s", err.Error())
		}

		if err := h.Start(chatID); err != nil {
			return fmt.Errorf("failed to send start msg: %s", err.Error())
		}
	}

	var row []string
	var IDs []string

	for _, c := range cars {
		row = append(row, fmt.Sprintf("%s | %s | %s | %s", c[1], c[2], c[3], c[4]))
		IDs = append(IDs, fmt.Sprintf("car id: %s", c[0]))
	}

	msg := tg.NewMessage(chatID, "Теперь выбери, строчку, которую хочешь изменить:")
	msg.ReplyMarkup = newKeyboard(row, IDs)

	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send msg with cars: %s", err.Error())
	}

	return nil
}

func (h *Handler) SendDocument(chatID int64) error {
	msg := tg.NewDocument(chatID, tg.FilePath(h.doc.DocPath))
	msg.Caption = "Вот твой рапорт 👆"
	if _, err := h.Send(msg); err != nil {
		return fmt.Errorf("failed to send msg with document: %s", err.Error())
	}

	return nil
}
