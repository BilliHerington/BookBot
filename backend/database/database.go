package database

import (
	"awesomeProject/logs"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type UserAppointmentData struct {
	AppointmentID   int
	UserName        string
	UserContact     string
	AppointmentDate string
	AppointmentTime string
	ServiceName     string
	EmployeeName    string
}
type AllAppointmentsData struct {
	AppointmentID   int
	UserName        string
	UserContact     string
	AppointmentDate string
	TimeStart       string
	TimeEnd         string
	ServiceID       int
	EmployeeID      int
	Status          bool
}
type AllServicesData struct {
	ServiceID    int
	Name         string
	Duration     string
	DefaultPrice float64
	ProPrice     float64
	Description  string
}
type AllEmployeesData struct {
	EmployeeID    int
	Name          string
	Level         string
	ContactNumber string
}
type AllScheduleData struct {
	ScheduleID int
	EmployeeID int
	WorkDate   time.Time
	StartTime  time.Time
	EndTime    time.Time
}

func LoadDB() (*sql.DB, error) {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("DBNAME")
	sslmode := os.Getenv("SSLMODE")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logs.ErrorLogger.Println("error sql open")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logs.ErrorLogger.Println("error ping")
		return nil, err
	}

	logs.InfoLogger.Printf("Successfully connected to database, name:\"%s\" ", dbname)
	return db, nil
}
