package database

import (
	"awesomeProject/logs"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetServices(db *sql.DB) []string {
	rows, err := db.Query("SELECT name FROM services")

	defer rows.Close()
	var serviceSlice []string

	for rows.Next() {
		service_name := ""
		err = rows.Scan(&service_name)
		if err != nil {
			logs.ErrorLogger.Fatal(err)
		}
		serviceSlice = append(serviceSlice, service_name)
	}
	return serviceSlice
}
func GetLevelAndCost(db *sql.DB, service string) ([]string, []string) {
	query := `
SELECT 
    e.level,
    CASE 
        WHEN e.level = 'Профессионал' THEN s.pro_price
        ELSE s.default_price
    END AS cost
FROM 
    employees e
CROSS JOIN 
    services s
WHERE 
    s.name = $1 -- Замените $1 на нужное название услуги
GROUP BY 
    e.level, s.pro_price, s.default_price;
`

	rows, err := db.Query(query, service)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	defer rows.Close()

	var levelAndCostSlice []string
	var employeeLevelSlice []string

	//var serviceCostByLevel float64
	for rows.Next() {
		employee_level := ""
		serviceCostByLevel := 0.0
		err = rows.Scan(&employee_level, &serviceCostByLevel)
		if err != nil {
			logs.ErrorLogger.Fatal(err)
		}
		resultString := fmt.Sprintf("%s %0.0f", employee_level, serviceCostByLevel)
		levelAndCostSlice = append(levelAndCostSlice, resultString)
		employeeLevelSlice = append(employeeLevelSlice, employee_level)
	}

	return levelAndCostSlice, employeeLevelSlice
}
func GetDate(db *sql.DB, level string) []string {
	query := `
SELECT DISTINCT 
    s.work_date
FROM 
    schedule s
JOIN 
    employees e ON s.employee_id = e.employee_id
WHERE 
    e.level = $1  -- Замените $1 на нужный уровень сотрудника
ORDER BY 
    s.work_date ASC;
`
	rows, err := db.Query(query, level)
	defer rows.Close()
	var workDateSlice []string
	for rows.Next() {
		date_name := time.Time{}
		err = rows.Scan(&date_name)
		if err != nil {
			logs.ErrorLogger.Fatal(err)
		}
		dateString := date_name.Format("2006-01-02")
		workDateSlice = append(workDateSlice, dateString)
	}
	return workDateSlice
}
func GetEmployeeIdByDate(db *sql.DB, date string, level string) []int {
	query := `
SELECT 
    s.employee_id
FROM 
    schedule s
JOIN 
    employees e ON s.employee_id = e.employee_id
WHERE 
    s.work_date = $1 -- Замените $1 на нужную дату работы
    AND e.level = $2  -- Замените $2 на нужный уровень сотрудника`
	rows, err := db.Query(query, date, level)
	defer rows.Close()
	var employeeIdSlice []int
	for rows.Next() {
		employeeId := 0
		err = rows.Scan(&employeeId)
		if err != nil {
			logs.ErrorLogger.Fatal(err)
		}
		employeeIdSlice = append(employeeIdSlice, employeeId)
	}
	return employeeIdSlice
}
func GetFreeTime(db *sql.DB, employee_id int, work_date string, service string, appointments_date string, status bool) ([]string, []string) {
	query := `
	WITH work_schedule AS (
    -- Рабочее время сотрудника на указанную дату
    SELECT 
        start_time,
        end_time
    FROM schedule
    WHERE employee_id = $1
      AND work_date = $2
),

	occupied_intervals AS (
    -- Занятые временные интервалы на указанную дату и услугу
    SELECT 
        a.time_start,
        a.time_end
    FROM 
        appointments a
    JOIN 
        services s ON a.service_id = s.service_id
    WHERE 
        a.employee_id = $1
        AND s.name = $3 -- Замените $3 на нужное название услуги
        AND a.appointment_date = $4
        AND a.status = $5
),

all_intervals AS (
    -- Объединяем все временные точки (начала и конца интервалов)
    SELECT start_time AS time_point FROM work_schedule
    UNION
    SELECT end_time AS time_point FROM work_schedule
    UNION
    SELECT time_start AS time_point FROM occupied_intervals
    UNION
    SELECT time_end AS time_point FROM occupied_intervals
),

ordered_intervals AS (
    -- Упорядочиваем все временные точки
    SELECT time_point
    FROM all_intervals
    ORDER BY time_point
),

free_intervals AS (
    -- Находим свободные интервалы между временными точками
    SELECT 
        LAG(time_point) OVER (ORDER BY time_point) AS start_time,
        time_point AS end_time
    FROM ordered_intervals
),

-- Исключение занятых интервалов и фильтрация свободных интервалов
available_intervals AS (
    SELECT 
        fi.start_time,
        fi.end_time
    FROM free_intervals fi
    JOIN work_schedule ws
      ON fi.start_time >= ws.start_time
     AND fi.end_time <= ws.end_time
    WHERE fi.start_time IS NOT NULL
      AND fi.end_time IS NOT NULL
      AND NOT EXISTS (
            SELECT 1
            FROM occupied_intervals oi
            WHERE (fi.start_time < oi.time_end AND fi.end_time > oi.time_start)
        )
)
-- Выборка окончательных свободных интервалов
SELECT 
    start_time,
    end_time
FROM available_intervals
WHERE start_time < end_time
ORDER BY start_time;`

	rows, err := db.Query(query, employee_id, work_date, service, appointments_date, status)
	defer rows.Close()
	var startTimeSlice []string
	var endTimeSlice []string
	for rows.Next() {
		startTime := ""
		endTime := ""
		err = rows.Scan(&startTime, &endTime)
		if err != nil {
			logs.ErrorLogger.Fatal(err)
		}
		formatedStartTime := startTime[strings.Index(startTime, "T")+1:]
		formatedEndtTime := endTime[strings.Index(endTime, "T")+1:]
		startTimeSlice = append(startTimeSlice, formatedStartTime)
		endTimeSlice = append(endTimeSlice, formatedEndtTime)
	}
	return startTimeSlice, endTimeSlice
}
func GetDuration(db *sql.DB, service string, startTime []string, endTime []string) []string {
	query := `
	SELECT duration	
	FROM services
	WHERE name = $1
`
	rows, err := db.Query(query, service)
	defer rows.Close()

	var duration string
	for rows.Next() {
		err = rows.Scan(&duration)
		if err != nil {
			logs.ErrorLogger.Fatal(err)
		}
	}
	dur, err := ParseDuration(duration)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}

	var timeSlots []string
	for i := 0; i < len(startTime); i++ {
		start, err2 := time.Parse("15:04:05Z07:00", startTime[i])
		if err2 != nil {
			logs.ErrorLogger.Fatal(err2)
		}
		end, err3 := time.Parse("15:04:05Z07:00", endTime[i])
		if err3 != nil {
			logs.ErrorLogger.Fatal(err3)
		}
		for t := start; t.Before(end); t = t.Add(dur) {
			formatedT := t.Format("15:04:05Z07:00")
			formatedT = formatedT[:len(formatedT)-4]
			timeSlots = append(timeSlots, formatedT)
		}
	}
	return timeSlots
}
func ParseDuration(duration string) (time.Duration, error) {
	parts := strings.Split(duration, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid duration format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}
func GetTimeSlots(freeTimeSlots map[int][]string) map[string][]int {
	reversedTimeSlotsRev := make(map[string][]int)
	for key, values := range freeTimeSlots {
		for _, value := range values {
			reversedTimeSlotsRev[value] = append(reversedTimeSlotsRev[value], key)
		}
	}
	return reversedTimeSlotsRev
}
func GetEmployeeName(db *sql.DB, employeeID int) string {
	query := `
	SELECT name FROM employees WHERE employee_id = $1`
	rows, err := db.Query(query, employeeID)
	defer rows.Close()
	employeeName := ""
	for rows.Next() {
		err = rows.Scan(&employeeName)
		if err != nil {
			logs.ErrorLogger.Fatal(err)
		}
	}
	return employeeName
}
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
			logs.ErrorLogger.Fatal(err)
		}
	}
	durationFormatted, err := ParseDuration(duration)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	appointmentsTimeFormatted, err := time.Parse("15:04", appointmentsTime)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
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
		logs.ErrorLogger.Fatal(err2)
	}
	return nil
}
func GetMyAppointments(db *sql.DB, userName string, userContact string) []UserAppointmentData {
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
		logs.ErrorLogger.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var userDataList []UserAppointmentData
	datetimeLayout := "2006-01-02T15:04:05Z" // Layout for parsing date and time with time zone from the database
	displayDateLayout := "02-01-2006"        // Desired date format for display
	displayTimeLayout := "15:04"             // Desired time format for display

	for rows.Next() {
		var userData UserAppointmentData
		var rawDatetime string
		var rawTime string
		err = rows.Scan(&userData.AppointmentID, &userData.UserName, &userData.UserContact, &rawDatetime, &rawTime, &userData.ServiceName, &userData.EmployeeName)
		if err != nil {
			logs.ErrorLogger.Fatalf("Failed to scan row: %v", err)
		}

		parsedDatetime, err := time.Parse(datetimeLayout, rawDatetime)
		if err != nil {
			logs.ErrorLogger.Fatalf("Failed to parse datetime: %v", err)
		}
		parsedTime, err := time.Parse(time.RFC3339, rawTime)
		if err != nil {
			logs.ErrorLogger.Fatalf("Failed to parse time: %v", err)
		}

		userData.AppointmentDate = parsedDatetime.Format(displayDateLayout)
		userData.AppointmentTime = parsedTime.Format(displayTimeLayout)

		userDataList = append(userDataList, userData)
	}

	if err = rows.Err(); err != nil {
		logs.ErrorLogger.Fatalf("Error occurred during iteration over rows: %v", err)
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
		logs.ErrorLogger.Fatal(err)
	}
	defer rows.Close()
	_, err = db.Exec(query, appointmentID)
	if err != nil {
		logs.ErrorLogger.Fatal(err)
	}
	return nil
}
