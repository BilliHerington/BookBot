package main

import (
	WI "awesomeProject/admin-web-interface"
	"awesomeProject/admin-web-interface/webHandlers"
	"context"
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"os/signal"
	"syscall"

	"awesomeProject/bot"
	"awesomeProject/database"
	"awesomeProject/logs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"time"
)

func main() {

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// chanel for catching system signals
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM) // catching "ctr + c", docker stop or another terminating

	//cwd, _ := os.Getwd()
	//logs.DebugLogger.Println("Current working directory:", cwd)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error loading .env file", err)
	}

	logs.InitLoggers()

	//running Database
	var db *sql.DB

	db, err = database.LoadDB()
	if err != nil {
		cancel()
	}

	// running HTTP-server
	go func() {
		err = runServer(db)
		if err != nil {
			logs.ErrorLogger.Printf("server error: %s", err)
			cancel()
		}
	}()

	// running Telegram-bot
	go func() {
		err = runTelegramBot(db)
		if err != nil {
			logs.ErrorLogger.Printf("telegram error: %s", err)
			cancel()
		}
	}()

	// waiting exit signal
	<-signChan
	logs.InfoLogger.Println("Terminating app...")
	cancel()

	time.Sleep(2 * time.Second)
	logs.InfoLogger.Printf("App terminated")

}

func runServer(db *sql.DB) error {

	ginMode := os.Getenv("GIN_MODE")
	gin.SetMode(ginMode)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	router.POST("/login", WI.Login)

	// Группа маршрутов, защищенных middleware для аутентификации
	authGroup := router.Group("/admin")
	authGroup.Use(WI.AuthMiddleware())

	authGroup.GET("/allServices", func(c *gin.Context) { webHandlers.AdminGetAllServices(c, db) })
	authGroup.POST("/addService", func(c *gin.Context) { webHandlers.AdminAddService(c, db) })
	authGroup.POST("/redactService", func(c *gin.Context) { webHandlers.AdminRedactService(c, db) })
	authGroup.DELETE("/deleteService/:id", func(c *gin.Context) { webHandlers.AdminDeleteService(c, db) })

	authGroup.GET("/allEmployees", func(c *gin.Context) { webHandlers.AdminGetAllEmployees(c, db) })
	authGroup.GET("/employeeLevels", func(c *gin.Context) { webHandlers.AdminGetEmployeeLevels(c, db) })
	authGroup.POST("/addEmployee", func(c *gin.Context) { webHandlers.AdminAddEmployee(c, db) })
	authGroup.POST("/redactEmployee", func(c *gin.Context) { webHandlers.AdminRedactEmployees(c, db) })
	authGroup.DELETE("/deleteEmployee/:id", func(c *gin.Context) { webHandlers.AdminDeleteEmployee(c, db) })

	authGroup.GET("/allAppointments", func(c *gin.Context) { webHandlers.AdminGetAllAppointments(c, db) })
	authGroup.DELETE("/deleteAppointment/:id", func(c *gin.Context) { webHandlers.AdminDeleteAppointment(c, db) })

	authGroup.GET("/allSchedules", func(c *gin.Context) { webHandlers.AdminGetAllSchedule(c, db) })
	authGroup.POST("/uploadNewSchedule", func(c *gin.Context) { webHandlers.AdminUploadFile(c, db) })
	authGroup.POST("/redactSchedule", func(c *gin.Context) { webHandlers.AdminRedactSchedule(c, db) })
	authGroup.DELETE("/deleteSchedule/:id", func(c *gin.Context) { webHandlers.AdminDeleteSchedule(c, db) })

	port := os.Getenv("LISTENING_PORT")

	logs.InfoLogger.Printf("Server starting on port %s...", port)

	err := router.Run(":" + port)
	if err != nil {
		return err
	}
	return nil
}

func runTelegramBot(db *sql.DB) error {
	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("BOTAPI"))
	if err != nil {
		logs.ErrorLogger.Println(err)
		return err
	}
	botAPI.Debug = false
	logs.InfoLogger.Printf("Authorized on Telegram account: %s", botAPI.Self.UserName)
	bot.RunBot(botAPI, db)
	return nil
}
