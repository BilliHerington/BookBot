package logs

import (
	"github.com/fatih/color"
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func InitLoggers() {
	// Создание цветных функций для логов
	infoColor := color.New(color.FgGreen).SprintFunc()
	errorColor := color.New(color.FgRed).SprintFunc()

	// Настройка логгеров с цветными префиксами
	InfoLogger = log.New(os.Stdout, infoColor("INFO: "), log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, errorColor("ERROR: "), log.Ldate|log.Ltime|log.Lshortfile)
}
