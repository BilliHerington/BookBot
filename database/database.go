package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type ConfigDB struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}
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

func LoadConfigDB() (*ConfigDB, error) {
	file, err := os.Open("database/configDB.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &ConfigDB{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}
func LoadDB(config *ConfigDB) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Successfully connected to \"%s\"!\n", config.DBName)
	return db, nil
}
