package bot

import (
	"awesomeProject/database/botDbRequests"
	"awesomeProject/logs"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func IsMasterSelection(message string, context *UserContext, db *sql.DB) bool {
	_, employLevelSlice := botDbRequests.GetLevelAndCost(db, context.SelectedService)
	for _, employLevel := range employLevelSlice {
		userMessage := strings.Split(message, " ")[0]
		if userMessage == employLevel {
			logs.DebugLogger.Printf("User selected level: %s", employLevel)
			context.SelectedLevel = employLevel
			return true
		}
	}
	return false
}
func HandleDateSelection(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.DebugLogger.Printf("Handling date selection for level: %s", context.SelectedLevel)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите дату")
	msg.ReplyMarkup = createKeyBoard(botDbRequests.GetDate(db, context.SelectedLevel), "Все услуги")
	bot.Send(msg)
}
