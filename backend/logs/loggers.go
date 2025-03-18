package logs

import (
	"github.com/fatih/color"
	"io"
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

func InitLoggers() {

	debugEnabled := os.Getenv("DEBUG")

	// Создание цветных функций для логов
	infoColor := color.New(color.FgGreen).SprintFunc()
	errorColor := color.New(color.FgRed).SprintFunc()
	debugColor := color.New(color.FgCyan).SprintFunc()

	// Настройка логгеров с цветными префиксами
	InfoLogger = log.New(os.Stdout, infoColor("INFO: "), log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, errorColor("ERROR: "), log.Ldate|log.Ltime|log.Lshortfile)

	if debugEnabled == "true" {
		InfoLogger.Println("Debug logging enabled")
		DebugLogger = log.New(os.Stdout, debugColor("DEBUG: "), log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		InfoLogger.Println("Debug logging disabled. set DEBUG=true in environment, if needed")
		DebugLogger = log.New(io.Discard, "", 0) // Отключаем вывод
	}

}
