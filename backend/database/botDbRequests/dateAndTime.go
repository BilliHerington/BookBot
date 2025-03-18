package botDbRequests

import (
	"awesomeProject/logs"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
	//logs.InfoLogger.Println("ТЕСТ:", level)
	rows, err := db.Query(query, level)
	if err != nil {
		logs.ErrorLogger.Printf("row error: ", err)
	}
	defer rows.Close()
	var workDateSlice []string
	for rows.Next() {
		date_name := time.Time{}
		err = rows.Scan(&date_name)
		if err != nil {
			logs.ErrorLogger.Println(err)
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
			logs.ErrorLogger.Println(err)
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
			logs.ErrorLogger.Println(err)
		}
		formatedStartTime := startTime[strings.Index(startTime, "T")+1:]
		formatedEndtTime := endTime[strings.Index(endTime, "T")+1:]
		startTimeSlice = append(startTimeSlice, formatedStartTime)
		endTimeSlice = append(endTimeSlice, formatedEndtTime)
	}
	return startTimeSlice, endTimeSlice
}
func GetDuration(db *sql.DB, service string, startTime []string, endTime []string) []string {
	query := `SELECT duration FROM services WHERE name = $1`
	rows, err := db.Query(query, service)
	if err != nil {
		logs.ErrorLogger.Println(err)
	}
	defer rows.Close()

	var duration string
	if rows.Next() {
		err = rows.Scan(&duration)
		if err != nil {
			logs.ErrorLogger.Println(err)
		}
	}

	dur, err := ParseDuration(duration)
	if err != nil {
		logs.ErrorLogger.Println(err)
	}

	var timeSlots []string
	for i := 0; i < len(startTime); i++ {
		start, err2 := time.Parse("15:04:05Z", startTime[i])
		end, err3 := time.Parse("15:04:05Z", endTime[i])
		if err2 != nil || err3 != nil {
			logs.ErrorLogger.Println(err2, err3)
		}

		for t := start; t.Add(dur).Before(end) || t.Add(dur).Equal(end); t = t.Add(dur) {
			timeSlots = append(timeSlots, t.Format("15:04"))
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
