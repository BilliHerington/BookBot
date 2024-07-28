package bot

import (
	"awesomeProject/logs"
	_ "awesomeProject/logs"
	"database/sql"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"sync"
)

type ConfigBot struct {
	APIKey string `json:"API_KEY"`
}

// загружаем конфиг для Бота (апи кей)
func LoadConfigBot(filename string) (*ConfigBot, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer bytes.Close()

	config := &ConfigBot{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

// фунция для запуска бота
func RunBot(bot *tgbotapi.BotAPI, db *sql.DB) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				handleStartCommand(bot, update.Message)
			}
		} else {
			handleTextmessage(bot, update.Message, db)
		}
	}
}

// Тут будет context
type UserContext struct {
	SelectedService     string
	SelectedLevel       string
	SelectedDate        string
	EmployeeIDSlice     []int
	TimeSlotsMap        map[string][]int
	TimeSlotsSlice      []string
	SelectedTime        string
	SelectedEmployeeId  int
	IsUserName          bool
	IsNumber            bool
	UserName            string
	UserNumber          string
	IsCreateAppointment bool
	IsMyAppointments    bool
	IsAppointmentName   bool
	IsAppointmentNumber bool
	AppointmentName     string
	AppointmentNumber   string
	IsAppontmentID      bool
}

var userContexts = make(map[int64]*UserContext)
var mu sync.Mutex // для синхронизации доступа к userContexts

type Route struct {
	Condition func(message string, context *UserContext, db *sql.DB) bool
	Handler   func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB)
}

var routes = []Route{
	{Condition: IsMainMenu, Handler: HandleMainMenu},
	{Condition: IsAllServices, Handler: HandleAllServices},
	{Condition: IsServiceSelection, Handler: HandleMasterSelection},
	{Condition: IsMasterSelection, Handler: HandleDateSelection},
	{Condition: IsDateSelection, Handler: HandleTimeSlotSelection},
	{Condition: IsTimeSlotSelection, Handler: HandleName},
	{Condition: IsName, Handler: HandleUserName},
	{Condition: IsNumber, Handler: HandleContactNumber},
	{Condition: IsCreateAppointment, Handler: HandleCreateAppointment},
	{Condition: isMyAppointments, Handler: handleMyAppointments},
	{Condition: isMyAppointmentName, Handler: handleMyAppointmentName},
	{Condition: isMyAppointmentNumber, Handler: handleUserAppointment},
	{Condition: isDeleteAppointment, Handler: handleDeleteAppointmentMessage},
	{Condition: isAppointmentID, Handler: handleDeleteAppointment},
}
