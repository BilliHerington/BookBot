package admin_web_interface

import (
	"awesomeProject/database"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// ---------------------------------------------------------SERVICES-----------------------------------------------------------------

func AdminGetAllServices(c *gin.Context, db *sql.DB) {
	services := database.GetAllServices(db)
	c.JSON(http.StatusOK, services)
}
func AdminAddService(c *gin.Context, db *sql.DB) {
	var service struct {
		Name        string  `json:"Name"`
		Duration    string  `json:"Duration"`
		DefPrice    float64 `json:"DefPrice"`
		ProPrice    float64 `json:"ProPrice"`
		Description string  `json:"Description"`
	}
	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := database.AddService(db, service.Name, service.Duration, service.DefPrice, service.ProPrice, service.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add service"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service added successfully"})
}
func AdminRedactService(c *gin.Context, db *sql.DB) {
	var service struct {
		ID          int     `json:"ID"`
		Name        string  `json:"Name"`
		Duration    string  `json:"Duration"`
		DefPrice    float64 `json:"DefPrice"`
		ProPrice    float64 `json:"ProPrice"`
		Description string  `json:"Description"`
	}
	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := database.RedactService(db, service.ID, service.Name, service.Duration, service.DefPrice, service.ProPrice, service.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service redacted successfully"})
}
func AdminDeleteService(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	err := database.DeleteService(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

// ---------------------------------------------------------EMPLOYEES-----------------------------------------------------------------

func AdminGetAllEmployees(c *gin.Context, db *sql.DB) {
	employees := database.GetAllEmployees(db)
	c.JSON(http.StatusOK, employees)
}
func AdminAddEmployee(c *gin.Context, db *sql.DB) {
	var employee struct {
		Name          string `json:"name"`
		Level         string `json:"level"`
		ContactNumber string `json:"contact_number"`
	}
	if err := c.BindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := database.AddEmployee(db, employee.Name, employee.Level, employee.ContactNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add employee"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee added successfully"})
}
func AdminRedactEmployees(c *gin.Context, db *sql.DB) {
	var employee struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Level         string `json:"level"`
		ContactNumber string `json:"contact_number"`
	}
	if err := c.BindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := database.RedactEmployee(db, employee.ID, employee.Name, employee.Level, employee.ContactNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee redacted successfully"})
}
func AdminDeleteEmployee(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	err := database.DeleteEmployee(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

// ---------------------------------------------------------APPOINTMENTS-----------------------------------------------------------------

func AdminGetAllAppointments(c *gin.Context, db *sql.DB) {
	appointments := database.GetAllAppointments(db)
	c.JSON(http.StatusOK, appointments)
}
func AdminDeleteAppointment(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	err := database.DeleteAppointment(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Appointment deleted successfully"})
}

// ---------------------------------------------------------SCHEDULES-----------------------------------------------------------------

func AdminGetAllSchedule(c *gin.Context, db *sql.DB) {
	schedule := database.GetAllSchedule(db)
	c.JSON(http.StatusOK, schedule)
}

//	func AdminUploadFile(c *gin.Context, db *sql.DB) {
//		file, err := c.FormFile("schedule")
//		if err != nil {
//			fmt.Printf("Error retrieving file: %v\n", err)
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving file"})
//			return
//		}
//
//		// Save the uploaded file to the server
//		filePath := "./" + file.Filename
//		err = c.SaveUploadedFile(file, filePath)
//		if err != nil {
//			fmt.Printf("Error saving file: %v\n", err)
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Error saving file"})
//			return
//		}
//
//		// Process the CSV file and insert data into the database
//		err = database.ProcessCSV(db, filePath)
//		if err != nil {
//			fmt.Printf("Error processing CSV: %v\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing CSV"})
//			return
//		}
//
//		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
//	}
func AdminUploadFile(c *gin.Context, db *sql.DB) {
	file, _, err := c.Request.FormFile("schedule")
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to get form file: %s", err.Error())
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read CSV: %s", err.Error())
		return
	}

	tx, err := db.Begin()
	if err != nil {
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
				c.String(http.StatusInternalServerError, "Failed to parse date: %v, rollback error: %v", err, rbErr)
				return
			}
			c.String(http.StatusInternalServerError, "Failed to parse date: %v", err)
			return
		}
		formattedDate := parsedDate.Format("2006-01-02")

		_, err = tx.Exec(`INSERT INTO schedule (employee_id, work_date, start_time, end_time) VALUES ($1, $2, $3, $4)`, employeeID, formattedDate, starttime, endtime)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				c.String(http.StatusInternalServerError, "Failed to insert record: %v, rollback error: %v", err, rbErr)
				return
			}
			c.String(http.StatusInternalServerError, "Failed to insert record: %v", err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
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

	err = database.RedactSchedule(db, schedule.ID, schedule.EmployeeID, workDateParsed, startTime, endTime)
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
	err := database.DeleteSchedule(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
}
