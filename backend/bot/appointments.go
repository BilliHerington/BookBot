package bot

import (
	"awesomeProject/database/botDbRequests"
	"awesomeProject/logs"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

// ---------------------------------------------------CreateAppointment-------------------------------------------
func IsMyAppointments(message string, context *UserContext, db *sql.DB) bool {
	if message == "Мои записи" {
		return true
	}
	return false
}
func HandleMyAppointments(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.DebugLogger.Printf("User requested Appointments")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите имя указанное при создании записи")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(msg)
	context.IsAppointmentName = true
}
func IsMyAppointmentName(message string, context *UserContext, db *sql.DB) bool {
	if context.IsAppointmentName {
		context.IsAppointmentName = false
		context.IsAppointmentNumber = true
		return true
	}
	return false
}
func HandleMyAppointmentName(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	context.AppointmentName = message.Text
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите номер указанный при создании записи")
	bot.Send(msg)
}
func IsMyAppointmentNumber(message string, context *UserContext, db *sql.DB) bool {
	if context.IsAppointmentNumber {
		context.IsAppointmentNumber = false
		return true
	}
	return false
}

func HandleUserAppointment(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	context.AppointmentNumber = message.Text
	userDataList := botDbRequests.GetMyAppointments(db, context.AppointmentName, context.AppointmentNumber)
	keyboardText := []string{"Удалить запись"}
	if len(userDataList) > 0 {
		for _, userData := range userDataList {
			messageText := fmt.Sprintf("Номер записи:%d\nИмя: %s\nНомер: %s\nДата: %s\nВремя: %s\nУслуга: %s\nИмя мастера: %s\n", userData.AppointmentID, userData.UserName, userData.UserContact, userData.AppointmentDate, userData.AppointmentTime, userData.ServiceName, userData.EmployeeName)
			msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
			msg.ReplyMarkup = createKeyBoard(keyboardText, "Главное меню")
			bot.Send(msg)
		}
		logs.DebugLogger.Printf("User get Appointments successfully")
	} else {
		messageText := fmt.Sprintf("Не удалось найти записи по имени: %s и номеру: %s", context.UserName, context.UserNumber)
		msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
		msg.ReplyMarkup = createKeyBoard(keyboardText, "Главное меню")
		bot.Send(msg)
		logs.DebugLogger.Printf("User no have Appointments")
	}
}

// ---------------------------------------------------DeleteAppointment-------------------------------------------
func IsDeleteAppointment(message string, context *UserContext, db *sql.DB) bool {
	if message == "Удалить запись" {
		return true
	}
	return false
}
func HandleDeleteAppointmentMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Напишите номер записи, которую хоитите удалить")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(msg)
	context.IsAppontmentID = true

}
func IsAppointmentID(message string, context *UserContext, db *sql.DB) bool {
	if context.IsAppontmentID {
		context.IsAppontmentID = false
		return true
	}
	return false
}
func HandleDeleteAppointment(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	messageText, err := strconv.Atoi(message.Text)
	if err != nil {
		logs.ErrorLogger.Println(err)
	}
	err2 := botDbRequests.DeleteAppointments(db, messageText)
	if err != nil {
		logs.ErrorLogger.Println(err2)
	} else {
		botMessage := fmt.Sprintf("Запись %s успешно удалена", message.Text)
		msg := tgbotapi.NewMessage(message.Chat.ID, botMessage)
		msg.ReplyMarkup = mainMenuKeyboard()
		bot.Send(msg)
		logs.DebugLogger.Printf("Delete Appointment successfully")
	}

}
