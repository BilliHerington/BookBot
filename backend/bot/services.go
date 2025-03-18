package bot

import (
	"awesomeProject/database/botDbRequests"
	"awesomeProject/logs"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func IsAllServices(message string, context *UserContext, db *sql.DB) bool {
	if message == "Все услуги" {
		return true
	}
	return false
}
func HandleAllServices(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.DebugLogger.Println("User requested all services")
	allServices := botDbRequests.GetServices(db)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Список услуг")
	msg.ReplyMarkup = createKeyBoard(allServices, "Главное меню")
	bot.Send(msg)
}

func IsServiceSelection(message string, context *UserContext, db *sql.DB) bool {
	for _, serviceName := range botDbRequests.GetServices(db) {
		if message == serviceName {
			logs.DebugLogger.Printf("User selected service: %s", serviceName)
			context.SelectedService = serviceName
			return true
		}
	}
	return false
}
func HandleMasterSelection(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.DebugLogger.Printf("Handling master selection for service: %s", context.SelectedService)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите мастера")
	priceForServiceSlice, _ := botDbRequests.GetLevelAndCost(db, context.SelectedService)
	//logs.InfoLogger.Println("ТЕСТ ВЫВОД priceForServiceSlice", priceForServiceSlice)

	msg.ReplyMarkup = createKeyBoard(priceForServiceSlice, "Все услуги")
	bot.Send(msg)
}
