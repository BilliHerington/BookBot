package bot

import (
	"awesomeProject/database/botDbRequests"
	"awesomeProject/logs"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"sort"
	"time"
)

// ----------------------------------------DATE_SELECTION----------------------------------------
func IsDateSelection(message string, context *UserContext, db *sql.DB) bool {
	for _, date := range botDbRequests.GetDate(db, context.SelectedLevel) {
		if message == date {
			logs.DebugLogger.Printf("User selected date: %s", date)
			context.SelectedDate = date
			logs.DebugLogger.Printf("Fetching employee IDs for date: %s and level: %s", context.SelectedDate, context.SelectedLevel)
			selectedEmployeeIDSlice := []int{}
			for _, employeeID := range botDbRequests.GetEmployeeIdByDate(db, context.SelectedDate, context.SelectedLevel) {
				selectedEmployeeIDSlice = append(selectedEmployeeIDSlice, employeeID)
			}
			context.EmployeeIDSlice = selectedEmployeeIDSlice
			logs.DebugLogger.Println("Fetching free time slots")
			freeTimeSlotsMap := make(map[int][]string)
			for _, employeeID := range context.EmployeeIDSlice {
				timeStartSlice, timeEndSlice := botDbRequests.GetFreeTime(db, employeeID, context.SelectedDate, context.SelectedService, context.SelectedDate, true)
				timeSlotSlice := botDbRequests.GetDuration(db, context.SelectedService, timeStartSlice, timeEndSlice)
				freeTimeSlotsMap[employeeID] = timeSlotSlice
			}
			freeTimeSlotsMapRev := botDbRequests.GetTimeSlots(freeTimeSlotsMap)
			context.TimeSlotsMap = freeTimeSlotsMapRev
			keys := []string{}
			for k := range freeTimeSlotsMapRev {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			timeSlice := []string{}
			for _, key := range keys {
				timeSlice = append(timeSlice, key)
			}
			context.TimeSlotsSlice = timeSlice
			return true
		}
	}
	return false
}
func HandleTimeSlotSelection(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	logs.DebugLogger.Println("Handling time slot selection")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите время")
	msg.ReplyMarkup = createKeyBoard(context.TimeSlotsSlice, "Все услуги")
	bot.Send(msg)
}

// ----------------------------------------TIME_SLOT_SELECTION----------------------------------------
func IsTimeSlotSelection(message string, context *UserContext, db *sql.DB) bool {
	for _, timeSlot := range context.TimeSlotsSlice {
		if message == timeSlot {
			logs.DebugLogger.Printf("User selected time: %s", timeSlot)
			context.SelectedTime = timeSlot
			employeeIDByTimeSLice, _ := context.TimeSlotsMap[timeSlot]
			rand.NewSource(time.Now().UnixNano())
			randomIndex := rand.Intn(len(employeeIDByTimeSLice))
			randomEmployeeID := employeeIDByTimeSLice[randomIndex]
			context.SelectedEmployeeId = randomEmployeeID
			return true
		}
	}
	return false
}
func HandleName(bot *tgbotapi.BotAPI, message *tgbotapi.Message, context *UserContext, db *sql.DB) {
	context.IsUserName = true
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите имя для создания записи")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(msg)
}
