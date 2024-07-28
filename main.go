package main

import (
	WI "awesomeProject/admin-web-interface"
	"awesomeProject/bot"
	"awesomeProject/database"
	"awesomeProject/logs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sync"
	"time"
)

func main() {
	// Инициализация логгеров
	logs.InitLoggers()
	infoLogger := logs.InfoLogger
	errorLogger := logs.ErrorLogger

	//---------------------------------------------DataBase---------------------------------------------------
	configDB, err := database.LoadConfigDB()
	if err != nil {
		errorLogger.Fatal(err)
	}
	db, err := database.LoadDB(configDB)
	if err != nil {
		errorLogger.Fatal(err)
	}
	defer db.Close()

	//---------------------------------------------Interface---------------------------------------------------
	var wg sync.WaitGroup
	wg.Add(2) // Задаем количество горутин

	// Запуск сервера в отдельной горутине
	go func() {
		defer wg.Done()

		router := gin.Default()
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
		router.POST("/login", WI.Login)

		// Группа маршрутов, защищенных middleware для аутентификации
		authGroup := router.Group("/admin")
		authGroup.Use(WI.AuthMiddleware())

		authGroup.GET("/allServices", func(c *gin.Context) { WI.AdminGetAllServices(c, db) })
		authGroup.POST("/addService", func(c *gin.Context) { WI.AdminAddService(c, db) })
		authGroup.POST("/redactService", func(c *gin.Context) { WI.AdminRedactService(c, db) })
		authGroup.DELETE("/deleteService/:id", func(c *gin.Context) { WI.AdminDeleteService(c, db) })

		authGroup.GET("/allEmployees", func(c *gin.Context) { WI.AdminGetAllEmployees(c, db) })
		authGroup.POST("/addEmployee", func(c *gin.Context) { WI.AdminAddEmployee(c, db) })
		authGroup.POST("/redactEmployee", func(c *gin.Context) { WI.AdminRedactEmployees(c, db) })
		authGroup.DELETE("/deleteEmployee/:id", func(c *gin.Context) { WI.AdminDeleteEmployee(c, db) })

		authGroup.GET("/allAppointments", func(c *gin.Context) { WI.AdminGetAllAppointments(c, db) })
		authGroup.DELETE("/deleteAppointment/:id", func(c *gin.Context) { WI.AdminDeleteAppointment(c, db) })

		authGroup.GET("/allSchedules", func(c *gin.Context) { WI.AdminGetAllSchedule(c, db) })
		authGroup.POST("/uploadNewSchedule", func(c *gin.Context) { WI.AdminUploadFile(c, db) })
		authGroup.POST("/redactSchedule", func(c *gin.Context) { WI.AdminRedactSchedule(c, db) })
		authGroup.DELETE("/deleteSchedule/:id", func(c *gin.Context) { WI.AdminDeleteSchedule(c, db) })

		err = router.Run(":8080")
		if err != nil {
			errorLogger.Fatal(err)
		} else {
			infoLogger.Println("Server started")
		}
	}()

	//---------------------------------------------TelegramBot--------------------------------------
	configBot, err := bot.LoadConfigBot("bot/configBot.json")
	if err != nil {
		errorLogger.Fatalf("Error loading config: %v", err)
	}

	botAPI, err := tgbotapi.NewBotAPI(configBot.APIKey)
	if err != nil {
		errorLogger.Fatal(err)
	}
	botAPI.Debug = false
	infoLogger.Printf("Authorized on account %s", botAPI.Self.UserName)

	// Запуск бота
	go func() {
		defer wg.Done()
		bot.RunBot(botAPI, db)
	}()

	// Ожидание завершения всех горутин
	wg.Wait()
}
