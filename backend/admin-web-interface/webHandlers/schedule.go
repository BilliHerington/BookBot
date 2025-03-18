package webHandlers

import (
	"awesomeProject/database/adminDbRequests"
	"awesomeProject/logs"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"time"
)

type ScheduleJSON struct {
	ID         int    `json:"id"`
	EmployeeID int    `json:"employee_id"`
	WorkDate   string `json:"work_date"`
	TimeStart  string `json:"time_start"`
	TimeEnd    string `json:"time_end"`
}

func AdminGetAllSchedule(c *gin.Context, db *sql.DB) {
	schedule, err := adminDbRequests.GetAllSchedule(db)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, schedule)
}

func AdminUploadFile(c *gin.Context, db *sql.DB) {
	file, _, err := c.Request.FormFile("schedule")
	if err != nil {
		logs.ErrorLogger.Println(err.Error())
		c.String(http.StatusBadRequest, "Failed to get form file: %s", err.Error())
		return
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			logs.ErrorLogger.Println(err.Error())
		}
	}(file)

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.String(http.StatusInternalServerError, "Failed to read CSV: %s", err.Error())
		return
	}

	tx, err := db.Begin()
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.String(http.StatusInternalServerError, "Failed to begin transaction: %s", err.Error())
		return
	}

	for i, record := range records {
		if len(record) < 4 {
			fmt.Printf("Skipping incomplete record at line %d: %v\n", i+1, record)
			continue
		}

		employeeID := record[0]
		workdate := record[1]
		starttime := record[2]
		endtime := record[3]

		// Преобразуем дату из формата DD.MM.YYYY в YYYY-MM-DD
		parsedDate, err := time.Parse("02.01.2006", workdate)
		if err != nil {
			fmt.Printf("Error parsing date at line %d: %v\n", i+1, err)
			if rbErr := tx.Rollback(); rbErr != nil {
				logs.ErrorLogger.Println(err)
				c.String(http.StatusInternalServerError, "Failed to parse date: %v, rollback error: %v", err, rbErr)
				return
			}
			logs.ErrorLogger.Println(err)
			c.String(http.StatusInternalServerError, "Failed to parse date: %v", err)
			return
		}
		formattedDate := parsedDate.Format("2006-01-02")

		_, err = tx.Exec(`INSERT INTO schedule (employee_id, work_date, start_time, end_time) VALUES ($1, $2, $3, $4)`, employeeID, formattedDate, starttime, endtime)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logs.ErrorLogger.Println(err)
				c.String(http.StatusInternalServerError, "Failed to insert record: %v, rollback error: %v", err, rbErr)
				return
			}
			logs.ErrorLogger.Println(err)
			c.String(http.StatusInternalServerError, "Failed to insert record: %v", err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		logs.ErrorLogger.Println(err)
		c.String(http.StatusInternalServerError, "Failed to commit transaction: %s", err.Error())
		return
	}

	c.String(http.StatusOK, "CSV data processed successfully")
}

func AdminRedactSchedule(c *gin.Context, db *sql.DB) {
	var schedule ScheduleJSON
	if err := c.BindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Преобразование строки времени в time.Time
	workDateParsed, startTime, endTime, err := schedule.ToTime(time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = adminDbRequests.RedactSchedule(db, schedule.ID, schedule.EmployeeID, workDateParsed, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule redacted successfully"})
}
func (s *ScheduleJSON) ToTime(workDate time.Time) (time.Time, time.Time, time.Time, error) {
	// Преобразование строки времени в time.Time
	dateStr := s.WorkDate
	startTimeStr := s.TimeStart
	endTimeStr := s.TimeEnd

	layout := "2006-01-02 15:04:05"

	// Преобразование строки времени в time.Time
	startTime, err := time.Parse(layout, dateStr+" "+startTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}
	endTime, err := time.Parse(layout, dateStr+" "+endTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}
	workDateParsed, err := time.Parse("2006-01-02", s.WorkDate)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}
	return workDateParsed, startTime, endTime, nil
}
func AdminDeleteSchedule(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	err := adminDbRequests.DeleteSchedule(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
}
