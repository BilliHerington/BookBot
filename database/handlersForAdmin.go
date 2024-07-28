package database

import (
	"awesomeProject/logs"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

//--------------------------------------------------------------SERVICES-----------------------------------------------------------------------------------

func GetAllServices(db *sql.DB) []AllServicesData {
	query := `
SELECT * FROM services ORDER BY service_id`
	rows, err := db.Query(query)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	defer rows.Close()
	serviceDataList := []AllServicesData{}
	for rows.Next() {
		serviceData := AllServicesData{}
		err = rows.Scan(&serviceData.ServiceID, &serviceData.Name, &serviceData.Duration, &serviceData.DefaultPrice, &serviceData.ProPrice, &serviceData.Description)
		if err != nil {
			logs.ErrorLogger.Fatal("Failed to scan row: %v", err)
		}
		serviceDataList = append(serviceDataList, serviceData)
	}
	return serviceDataList
}
func AddService(db *sql.DB, name string, duration string, defprice float64, proprice float64, description string) error {

	query := `INSERT INTO services (name, duration, default_price, pro_price, description) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, name, duration, defprice, proprice, description)
	if err != nil {
		logs.ErrorLogger.Fatalf("Failed to execute query: %v", err)
	}
	return nil
}
func RedactService(db *sql.DB, serviceID int, serviceName string, duration string, defprice float64, proprice float64, description string) error {
	query := `UPDATE services 
SET name = $2, duration = $3, default_price = $4, pro_price = $5, description = $6
WHERE service_id = $1`
	_, err := db.Exec(query, serviceID, serviceName, duration, defprice, proprice, description)
	if err != nil {
		return err
	}
	return nil
}
func DeleteService(db *sql.DB, serviceID string) error {
	query := `DELETE FROM services WHERE service_id = $1`
	_, err := db.Exec(query, serviceID)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	return nil
}

//--------------------------------------------------------------EMPLOYEES-----------------------------------------------------------------------------------

func GetAllEmployees(db *sql.DB) []AllEmployeesData {
	query := `
SELECT * FROM employees ORDER BY employee_id`
	rows, err := db.Query(query)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	defer rows.Close()
	employeeDataList := []AllEmployeesData{}
	for rows.Next() {
		employeeData := AllEmployeesData{}
		err = rows.Scan(&employeeData.EmployeeID, &employeeData.Name, &employeeData.Level, &employeeData.ContactNumber)
		if err != nil {
			logs.ErrorLogger.Fatal(err)
		}
		employeeDataList = append(employeeDataList, employeeData)
	}
	return employeeDataList
}
func AddEmployee(db *sql.DB, employeeName string, level string, contactnumber string) error {
	query := `INSERT INTO employees (name, level , contactnumber) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, employeeName, level, contactnumber)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	return nil
}
func RedactEmployee(db *sql.DB, employeeID int, employeeName string, level string, contactnumber string) error {
	query := `UPDATE employees 
SET name = $2, level = $3, contactnumber = $4
WHERE employee_id = $1`
	_, err := db.Exec(query, employeeID, employeeName, level, contactnumber)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	return nil
}
func DeleteEmployee(db *sql.DB, employeeID string) error {
	query := `DELETE FROM employees WHERE employee_id = $1`
	_, err := db.Exec(query, employeeID)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	return nil
}

// --------------------------------------------------------------APPOINTMENTS-----------------------------------------------------------------------------------
func GetAllAppointments(db *sql.DB) []AllAppointmentsData {
	query := `
SELECT appointment_id, client_name, contact_number, appointment_date, time_start, time_end, service_id, employee_id, status
FROM appointments ORDER BY appointment_id`
	rows, err := db.Query(query)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	defer rows.Close()
	userDataList := []AllAppointmentsData{}
	for rows.Next() {
		userData := AllAppointmentsData{}
		err = rows.Scan(&userData.AppointmentID, &userData.UserName, &userData.UserContact, &userData.AppointmentDate, &userData.TimeStart, &userData.TimeEnd, &userData.ServiceID, &userData.EmployeeID, &userData.Status)
		if err != nil {
			logs.ErrorLogger.Fatal("Failed to scan row: %v", err)
		}
		userDataList = append(userDataList, userData)
	}
	return userDataList
}
func DeleteAppointment(db *sql.DB, appointmentID string) error {
	query := `DELETE FROM appointments WHERE appointment_id = $1`
	_, err := db.Exec(query, appointmentID)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	return nil
}

// --------------------------------------------------------------SCHEDULES-----------------------------------------------------------------------------------
func GetAllSchedule(db *sql.DB) []AllScheduleData {
	query := `SELECT schedule_id, employee_id, work_date, start_time, end_time FROM schedule ORDER BY work_date`
	rows, err := db.Query(query)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	defer rows.Close()

	scheduleDataList := []AllScheduleData{}
	for rows.Next() {
		scheduleData := AllScheduleData{}
		err = rows.Scan(&scheduleData.ScheduleID, &scheduleData.EmployeeID, &scheduleData.WorkDate, &scheduleData.StartTime, &scheduleData.EndTime)
		if err != nil {
			logs.ErrorLogger.Fatal(err)
		}
		scheduleDataList = append(scheduleDataList, scheduleData)
	}
	return scheduleDataList
}
func RedactSchedule(db *sql.DB, scheduleID int, employeeID int, workdate time.Time, starttime time.Time, endtime time.Time) error {
	query := `UPDATE schedule 
SET employee_id = $2, work_date = $3, start_time = $4, end_time = $5
WHERE schedule_id = $1`
	_, err := db.Exec(query, scheduleID, employeeID, workdate, starttime, endtime)
	if err != nil {
		logs.ErrorLogger.Fatal("Error updating schedule: ", err)
		return err
	}
	return nil
}
func DeleteSchedule(db *sql.DB, scheduleID string) error {
	query := `DELETE FROM schedule WHERE schedule_id = $1`
	_, err := db.Exec(query, scheduleID)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	return nil
}
func ProcessCSV(db *sql.DB, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // Используем точку с запятой как разделитель

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read csv data: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}

	for i, record := range records {
		fmt.Printf("Raw record at line %d: %v\n", i+1, record) // Добавлен вывод для отладки
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
				return fmt.Errorf("parse error: %v, rollback error: %v", err, rbErr)
			}
			return err
		}
		formattedDate := parsedDate.Format("2006-01-02")

		fmt.Printf("Processing record at line %d: EmployeeID: %s, WorkDate: %s, StartTime: %s, EndTime: %s\n", i+1, employeeID, formattedDate, starttime, endtime)

		_, err = tx.Exec(`INSERT INTO schedule (employee_id, work_date, start_time, end_time) VALUES ($1, $2, $3, $4)`, employeeID, formattedDate, starttime, endtime)
		if err != nil {
			fmt.Printf("Error inserting record at line %d: %v\n", i+1, err)
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
			}
			return err
		}
		fmt.Printf("Successfully inserted record at line %d: %v\n", i+1, record)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}
