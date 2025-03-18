package bot

import (
	"awesomeProject/database/botDbRequests"
	"awesomeProject/logs"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// ----------------------------------------NAME----------------------------------------
func IsName(message string, context *UserContext, db *sql.DB) bool {
	if context.IsUserName {
		return true
	}
	return false
}
func HandleUserName(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	context.UserName = message.Text
	logs.DebugLogger.Printf("User name: %s", context.UserName)
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
	logs.DebugLogger.Printf("User number: %s", context.UserNumber)
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
	err := botDbRequests.СreateAppointments(db, context.SelectedService, context.SelectedEmployeeId, context.SelectedDate, context.SelectedTime, context.UserName, context.UserNumber)
	if err != nil {
		return
	} else {
		logs.DebugLogger.Printf("Create appointments successfully")
		masterName := botDbRequests.GetEmployeeName(db, context.SelectedEmployeeId)
		messageText := fmt.Sprintf("Вы успешно записаны. \n Имя: %s\n Номер: %s\n Дата: %s\n Время: %s\n Услуга: %s\n Имя мастер: %s\n", context.UserName, context.UserNumber, context.SelectedDate, context.SelectedTime, context.SelectedService, masterName)
		msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
		msg.ReplyMarkup = mainMenuKeyboard()
		bot.Send(msg)
	}
}
