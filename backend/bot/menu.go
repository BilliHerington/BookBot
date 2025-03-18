package bot

import (
	"awesomeProject/logs"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sync"
)

var userContexts = make(map[int64]*UserContext)
var mu sync.Mutex // для синхронизации доступа к userContexts

func IsMainMenu(message string, context *UserContext, db *sql.DB) bool {
	if message == "Главное меню" {
		return true
	}
	return false
}
func HandleMainMenu(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.DebugLogger.Println("User requested main menu")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Главное меню")
	msg.ReplyMarkup = mainMenuKeyboard()
	bot.Send(msg)
}

func HandleTextmessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	mu.Lock()
	context, exists := userContexts[message.Chat.ID]
	if !exists {
		context = &UserContext{}
		userContexts[message.Chat.ID] = context
	}
	mu.Unlock()
	//logs.DebugLogger.Println("test user ctx:", userContexts[message.Chat.ID])

	logs.DebugLogger.Printf("Received message: %s", message.Text)
	for _, route := range Routes {
		if route.Condition(message.Text, context, db) {
			route.Handler(bot, message, context, db)
			return
		}
	}
	logs.DebugLogger.Printf("Message not found in routes: %s", message.Text)
	handleUnknownCommand(bot, message)
}

func handleUnknownCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	logs.DebugLogger.Println("Handling unknown command")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я вас не понимаю или фунцкионал кнопки еще завершен")
	bot.Send(msg)
}

// обработка команды start
func HandleStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
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
