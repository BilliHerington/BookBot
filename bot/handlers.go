package bot

import (
	"awesomeProject/database"
	"awesomeProject/logs"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

//var db *sql.DB

// ----------------------------------------START----------------------------------------
func IsMainMenu(message string, context *UserContext, db *sql.DB) bool {
	if message == "Главное меню" {
		return true
	}
	return false
}
func HandleMainMenu(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.InfoLogger.Println("User requested main menu")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Главное меню")
	msg.ReplyMarkup = mainMenuKeyboard()
	bot.Send(msg)
}

// ----------------------------------------ALL_SERVICES----------------------------------------
func IsAllServices(message string, context *UserContext, db *sql.DB) bool {
	if message == "Все услуги" {
		return true
	}
	return false
}
func HandleAllServices(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.InfoLogger.Println("User requested all services")
	allServices := database.GetServices(db)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Список услуг")
	msg.ReplyMarkup = createKeyBoard(allServices, "Главное меню")
	bot.Send(msg)
}

// ----------------------------------------SERVICE_SELECTION----------------------------------------
func IsServiceSelection(message string, context *UserContext, db *sql.DB) bool {
	for _, serviceName := range database.GetServices(db) {
		if message == serviceName {
			logs.InfoLogger.Printf("User selected service: %s", serviceName)
			context.SelectedService = serviceName
			return true
		}
	}
	return false
}
func HandleMasterSelection(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.InfoLogger.Printf("Handling master selection for service: %s", context.SelectedService)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите мастера")
	priceForServiceSlice, _ := database.GetLevelAndCost(db, context.SelectedService)
	msg.ReplyMarkup = createKeyBoard(priceForServiceSlice, "Все услуги")
	bot.Send(msg)
}

// ----------------------------------------MASTER_SELECTION----------------------------------------
func IsMasterSelection(message string, context *UserContext, db *sql.DB) bool {
	_, employLevelSlice := database.GetLevelAndCost(db, context.SelectedService)
	for _, employLevel := range employLevelSlice {
		userMessage := strings.Split(message, " ")[0]
		if userMessage == employLevel {
			logs.InfoLogger.Printf("User selected level: %s", employLevel)
			context.SelectedLevel = employLevel
			return true
		}
	}
	return false
}
func HandleDateSelection(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.InfoLogger.Printf("Handling date selection for level: %s", context.SelectedLevel)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите дату")
	msg.ReplyMarkup = createKeyBoard(database.GetDate(db, context.SelectedLevel), "Все услуги")
	bot.Send(msg)
}

// ----------------------------------------DATE_SELECTION----------------------------------------
func IsDateSelection(message string, context *UserContext, db *sql.DB) bool {
	for _, date := range database.GetDate(db, context.SelectedLevel) {
		if message == date {
			logs.InfoLogger.Printf("User selected date: %s", date)
			context.SelectedDate = date
			logs.InfoLogger.Printf("Fetching employee IDs for date: %s and level: %s", context.SelectedDate, context.SelectedLevel)
			selectedEmployeeIDSlice := []int{}
			for _, employeeID := range database.GetEmployeeIdByDate(db, context.SelectedDate, context.SelectedLevel) {
				selectedEmployeeIDSlice = append(selectedEmployeeIDSlice, employeeID)
			}
			context.EmployeeIDSlice = selectedEmployeeIDSlice
			logs.InfoLogger.Println("Fetching free time slots")
			freeTimeSlotsMap := make(map[int][]string)
			for _, employeeID := range context.EmployeeIDSlice {
				timeStartSlice, timeEndSlice := database.GetFreeTime(db, employeeID, context.SelectedDate, context.SelectedService, context.SelectedDate, true)
				timeSlotSlice := database.GetDuration(db, context.SelectedService, timeStartSlice, timeEndSlice)
				freeTimeSlotsMap[employeeID] = timeSlotSlice
			}
			freeTimeSlotsMapRev := database.GetTimeSlots(freeTimeSlotsMap)
			context.TimeSlotsMap = freeTimeSlotsMapRev
			keys := []string{}
			for k := range freeTimeSlotsMapRev {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			timeSlice := []string{}
			for _, key := range keys {
				timeSlice = append(timeSlice, key)
			}
			context.TimeSlotsSlice = timeSlice
			return true
		}
	}
	return false
}
func HandleTimeSlotSelection(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.InfoLogger.Println("Handling time slot selection")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите время")
	msg.ReplyMarkup = createKeyBoard(context.TimeSlotsSlice, "Все услуги")
	bot.Send(msg)
}

// ----------------------------------------TIME_SLOT_SELECTION----------------------------------------
func IsTimeSlotSelection(message string, context *UserContext, db *sql.DB) bool {
	for _, timeSlot := range context.TimeSlotsSlice {
		if message == timeSlot {
			logs.InfoLogger.Printf("User selected time: %s", timeSlot)
			context.SelectedTime = timeSlot
			employeeIDByTimeSLice, _ := context.TimeSlotsMap[timeSlot]
			rand.NewSource(time.Now().UnixNano())
			randomIndex := rand.Intn(len(employeeIDByTimeSLice))
			randomEmployeeID := employeeIDByTimeSLice[randomIndex]
			context.SelectedEmployeeId = randomEmployeeID
			return true
		}
	}
	return false
}
func HandleName(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	context.IsUserName = true
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите имя для создания записи")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(msg)
}

// ----------------------------------------NAME----------------------------------------
func IsName(message string, context *UserContext, db *sql.DB) bool {
	if context.IsUserName {
		return true
	}
	return false
}
func HandleUserName(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	context.UserName = message.Text
	logs.InfoLogger.Printf("User name: %s", context.UserName)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите номер телефона для создания записи")
	bot.Send(msg)
	context.IsNumber = true
	context.IsUserName = false
}

// ----------------------------------------NUMBER----------------------------------------
func IsNumber(message string, context *UserContext, db *sql.DB) bool {
	if context.IsNumber {
		return true
	}
	return false
}
func HandleContactNumber(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	context.IsNumber = false
	context.UserNumber = message.Text
	logs.InfoLogger.Printf("User number: %s", context.UserNumber)
	messageText := []string{"Записаться"}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Подтвердите создание")
	msg.ReplyMarkup = createKeyBoard(messageText, "Главное меню")
	bot.Send(msg)
}

// ----------------------------------------CREATE_APPOINTMENT----------------------------------------
func IsCreateAppointment(message string, context *UserContext, db *sql.DB) bool {
	if message == "Записаться" {
		return true
	}
	return false
}
func HandleCreateAppointment(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	err := database.СreateAppointments(db, context.SelectedService, context.SelectedEmployeeId, context.SelectedDate, context.SelectedTime, context.UserName, context.UserNumber)
	if err != nil {
		return
	} else {
		logs.InfoLogger.Printf("Create appointments successfully")
		masterName := database.GetEmployeeName(db, context.SelectedEmployeeId)
		messageText := fmt.Sprintf("Вы успешно записаны. \n Имя: %s\n Номер: %s\n Дата: %s\n Время: %s\n Услуга: %s\n Имя мастер: %s\n", context.UserName, context.UserNumber, context.SelectedDate, context.SelectedTime, context.SelectedService, masterName)
		msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
		msg.ReplyMarkup = mainMenuKeyboard()
		bot.Send(msg)
	}
}

// ---------------------------------------------------CreateAppointment-------------------------------------------
func isMyAppointments(message string, context *UserContext, db *sql.DB) bool {
	if message == "Мои записи" {
		return true
	}
	return false
}
func handleMyAppointments(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.InfoLogger.Printf("User requested Appointments")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите имя указанное при создании записи")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(msg)
	context.IsAppointmentName = true
}
func isMyAppointmentName(message string, context *UserContext, db *sql.DB) bool {
	if context.IsAppointmentName {
		context.IsAppointmentName = false
		context.IsAppointmentNumber = true
		return true
	}
	return false
}
func handleMyAppointmentName(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	context.AppointmentName = message.Text
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите номер указанный при создании записи")
	bot.Send(msg)
}
func isMyAppointmentNumber(message string, context *UserContext, db *sql.DB) bool {
	if context.IsAppointmentNumber {
		context.IsAppointmentNumber = false
		return true
	}
	return false
}

func handleUserAppointment(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	context.AppointmentNumber = message.Text
	userDataList := database.GetMyAppointments(db, context.AppointmentName, context.AppointmentNumber)
	keyboardText := []string{"Удалить запись"}
	if len(userDataList) > 0 {
		for _, userData := range userDataList {
			messageText := fmt.Sprintf("Номер записи:%d\nИмя: %s\nНомер: %s\nДата: %s\nВремя: %s\nУслуга: %s\nИмя мастера: %s\n", userData.AppointmentID, userData.UserName, userData.UserContact, userData.AppointmentDate, userData.AppointmentTime, userData.ServiceName, userData.EmployeeName)
			msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
			msg.ReplyMarkup = createKeyBoard(keyboardText, "Главное меню")
			bot.Send(msg)
		}
		logs.InfoLogger.Printf("User get Appointments successfully")
	} else {
		messageText := fmt.Sprintf("Не удалось найти записи по имени: %s и номеру: %s", context.UserName, context.UserNumber)
		msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
		msg.ReplyMarkup = createKeyBoard(keyboardText, "Главное меню")
		bot.Send(msg)
		logs.InfoLogger.Printf("User no have Appointments")
	}
}

// ---------------------------------------------------DeleteAppointment-------------------------------------------
func isDeleteAppointment(message string, context *UserContext, db *sql.DB) bool {
	if message == "Удалить запись" {
		return true
	}
	return false
}
func handleDeleteAppointmentMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Напишите номер записи, которую хоитите удалить")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(msg)
	context.IsAppontmentID = true

}
func isAppointmentID(message string, context *UserContext, db *sql.DB) bool {
	if context.IsAppontmentID {
		context.IsAppontmentID = false
		return true
	}
	return false
}
func handleDeleteAppointment(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	messageText, err := strconv.Atoi(message.Text)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	err2 := database.DeleteAppointments(db, messageText)
	if err != nil {
		logs.ErrorLogger.Fatal(err2)
	} else {
		botMessage := fmt.Sprintf("Запись %s успешно удалена", message.Text)
		msg := tgbotapi.NewMessage(message.Chat.ID, botMessage)
		msg.ReplyMarkup = mainMenuKeyboard()
		bot.Send(msg)
		logs.InfoLogger.Printf("Delete Appointment successfully")
	}

}

// ---------------------------------------------------OtherHandlers-------------------------------------------
func handleTextmessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	mu.Lock()
	context, exists := userContexts[message.Chat.ID]
	if !exists {
		context = &UserContext{}
		userContexts[message.Chat.ID] = context
	}
	mu.Unlock()

	logs.InfoLogger.Printf("Received message: %s", message.Text)
	for _, route := range routes {
		if route.Condition(message.Text, context, db) {
			route.Handler(bot, message, context, db)
			return
		}
	}
	logs.InfoLogger.Printf("Message not found in routes: %s", message.Text)
	handleUnknownCommand(bot, message)
}

func handleUnknownCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	logs.InfoLogger.Println("Handling unknown command")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я вас не понимаю или фунцкионал кнопки еще завершен")
	bot.Send(msg)
}

// обработка команды start
func handleStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Здравствуйте, чем могу помочь")
	msg.ReplyMarkup = mainMenuKeyboard()
	bot.Send(msg)
}

// отрисовка главного меню
func mainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Все услуги"),
			tgbotapi.NewKeyboardButton("Мои записи"),
		),
	)
}

// createKeyBoard создает клавиатуру с кнопками для каждого названия и добавляет кнопку возврата к более высокому уровню
func createKeyBoard(options []string, backOption string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton
	for _, option := range options {
		buttonRow := []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton(option),
		}
		rows = append(rows, buttonRow)
	}
	// Добавление кнопки возврата
	if backOption != "" {
		backButtonRow := []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton(backOption),
		}
		rows = append(rows, backButtonRow)
	}

	return tgbotapi.NewReplyKeyboard(rows...)
}
