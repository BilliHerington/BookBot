package botDbRequests

import (
	"awesomeProject/database"
	"awesomeProject/logs"
	"database/sql"
	"time"
)

func СreateAppointments(db *sql.DB, serviceName string, employeeID int, appointmentsDate string, appointmentsTime string, clientName string, clientNumber string) error {
	query := `select service_id, duration
FROM services
WHERE name = $1 `
	rows, err := db.Query(query, serviceName)
	defer rows.Close()
	serviceID := 0
	duration := ""
	for rows.Next() {
		err = rows.Scan(&serviceID, &duration)
		if err != nil {
			logs.ErrorLogger.Println(err)
		}
	}
	durationFormatted, err := ParseDuration(duration)
	if err != nil {
		logs.ErrorLogger.Println(err)
	}
	appointmentsTimeFormatted, err := time.Parse("15:04", appointmentsTime)
	if err != nil {
		logs.ErrorLogger.Println(err)
	}
	resultTime := appointmentsTimeFormatted.Add(durationFormatted)
	resultTimeStr := resultTime.Format("15:04")
	query2 := `
	INSERT INTO appointments (client_name, service_id, employee_id, appointment_date, time_start, time_end, status, contact_number)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8)
`
	_, err2 := db.Exec(query2, clientName, serviceID, employeeID, appointmentsDate, appointmentsTime, resultTimeStr, true, clientNumber)
	if err2 != nil {
		logs.ErrorLogger.Println(err2)
	}
	return nil
}
func GetMyAppointments(db *sql.DB, userName string, userContact string) []database.UserAppointmentData {
	query := `
		SELECT 
			a.appointment_id,
			a.client_name, 
			a.contact_number,
			a.appointment_date,
			a.time_start,
			s.name AS service_name,
			e.name AS employee_name
		FROM 
			appointments a
		JOIN 
			services s ON a.service_id = s.service_id
		JOIN 
			employees e ON a.employee_id = e.employee_id
		WHERE 
			client_name = $1 AND contact_number = $2`

	rows, err := db.Query(query, userName, userContact)
	if err != nil {
		logs.ErrorLogger.Printf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var userDataList []database.UserAppointmentData
	datetimeLayout := "2006-01-02T15:04:05Z" // Layout for parsing date and time with time zone from the database
	displayDateLayout := "02-01-2006"        // Desired date format for display
	displayTimeLayout := "15:04"             // Desired time format for display

	for rows.Next() {
		var userData database.UserAppointmentData
		var rawDatetime string
		var rawTime string
		err = rows.Scan(&userData.AppointmentID, &userData.UserName, &userData.UserContact, &rawDatetime, &rawTime, &userData.ServiceName, &userData.EmployeeName)
		if err != nil {
			logs.ErrorLogger.Printf("Failed to scan row: %v", err)
		}

		parsedDatetime, err := time.Parse(datetimeLayout, rawDatetime)
		if err != nil {
			logs.ErrorLogger.Printf("Failed to parse datetime: %v", err)
		}
		parsedTime, err := time.Parse(time.RFC3339, rawTime)
		if err != nil {
			logs.ErrorLogger.Printf("Failed to parse time: %v", err)
		}

		userData.AppointmentDate = parsedDatetime.Format(displayDateLayout)
		userData.AppointmentTime = parsedTime.Format(displayTimeLayout)

		userDataList = append(userDataList, userData)
	}

	if err = rows.Err(); err != nil {
		logs.ErrorLogger.Printf("Error occurred during iteration over rows: %v", err)
	}
	return userDataList
}
func DeleteAppointments(db *sql.DB, appointmentID int) error {
	query := `
DELETE FROM appointments
WHERE appointment_id = $1;
`
	rows, err := db.Query(query, appointmentID)
	if err != nil {
		logs.ErrorLogger.Println(err)
	}
	defer rows.Close()
	_, err = db.Exec(query, appointmentID)
	if err != nil {
		logs.ErrorLogger.Println(err)
	}
	return nil
}
