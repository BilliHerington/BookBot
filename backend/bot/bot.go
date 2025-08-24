package bot

import (
	"awesomeProject/logs"
	_ "awesomeProject/logs"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"sync"
)

// фунция для запуска бота
func RunBot(bot *tgbotapi.BotAPI, db *sql.DB) {

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err := bot.GetUpdatesChan(u)
		if err != nil {
			logs.ErrorLogger.Println(err)
		}

		for update := range updates {
			if update.Message == nil {
				continue
			}
			HandleUpdates(bot, &update, db)
		}
	}()
	if os.Getenv("DEBUG") == "true" {
		// создаём отдельный тестовый роутер для /fakeUpdate
		testRouter := gin.Default()
		testRouter.POST("/fakeUpdate", func(c *gin.Context) {
			var update tgbotapi.Update
			if err := c.BindJSON(&update); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			//fmt.Println("db contain:", db)
			// передаём update в реальный BotAPI
			HandleUpdates(bot, &update, db)
			c.JSON(200, gin.H{"status": "ok"})
		})

		// запускаем тестовый сервер на другом порту (например 8081)
		go testRouter.Run(":8081")
	}

	logs.InfoLogger.Printf("Telegram bot initialized and fakeUpdate test endpoint ready on :8081")
}

func HandleUpdates(bot *tgbotapi.BotAPI, update *tgbotapi.Update, db *sql.DB) {

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			HandleStartCommand(bot, update.Message)
		}
	} else {
		HandleTextMessage(bot, update.Message, db)
	}
}

// Тут будет context
type UserContext struct {
	mu sync.RWMutex

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

type Route struct {
	Condition func(message string, context *UserContext, db *sql.DB) bool
	Handler   func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB)
}

var Routes = []Route{
	{Condition: IsMainMenu, Handler: HandleMainMenu},
	{Condition: IsAllServices, Handler: HandleAllServices},
	{Condition: IsServiceSelection, Handler: HandleMasterSelection},
	{Condition: IsMasterSelection, Handler: HandleDateSelection},
	{Condition: IsDateSelection, Handler: HandleTimeSlotSelection},
	{Condition: IsTimeSlotSelection, Handler: HandleName},
	{Condition: IsName, Handler: HandleUserName},
	{Condition: IsNumber, Handler: HandleContactNumber},
	{Condition: IsCreateAppointment, Handler: HandleCreateAppointment},
	{Condition: IsMyAppointments, Handler: HandleMyAppointments},
	{Condition: IsMyAppointmentName, Handler: HandleMyAppointmentName},
	{Condition: IsMyAppointmentNumber, Handler: HandleUserAppointment},
	{Condition: IsDeleteAppointment, Handler: HandleDeleteAppointmentMessage},
	{Condition: IsAppointmentID, Handler: HandleDeleteAppointment},
}
