package adminDbRequests

import (
	"awesomeProject/database"
	"awesomeProject/logs"
	"database/sql"
	"time"
)

func GetAllSchedule(db *sql.DB) ([]database.AllScheduleData, error) {
	query := `SELECT schedule_id, employee_id, work_date, start_time, end_time FROM schedule ORDER BY work_date`
	rows, err := db.Query(query)
	if err != nil {
		logs.ErrorLogger.Println(err)
		return nil, err
	}
	defer rows.Close()

	scheduleDataList := []database.AllScheduleData{}
	for rows.Next() {
		scheduleData := database.AllScheduleData{}
		err = rows.Scan(&scheduleData.ScheduleID, &scheduleData.EmployeeID, &scheduleData.WorkDate, &scheduleData.StartTime, &scheduleData.EndTime)
		if err != nil {
			logs.ErrorLogger.Println(err)
			return nil, err
		}
		scheduleDataList = append(scheduleDataList, scheduleData)
	}
	return scheduleDataList, nil
}
func RedactSchedule(db *sql.DB, scheduleID int, employeeID int, workdate time.Time, starttime time.Time, endtime time.Time) error {
	query := `UPDATE schedule 
SET employee_id = $2, work_date = $3, start_time = $4, end_time = $5
WHERE schedule_id = $1`
	_, err := db.Exec(query, scheduleID, employeeID, workdate, starttime, endtime)
	if err != nil {
		logs.ErrorLogger.Println("Error updating schedule: ", err)
		return err
	}
	return nil
}
func DeleteSchedule(db *sql.DB, scheduleID string) error {
	query := `DELETE FROM schedule WHERE schedule_id = $1`
	_, err := db.Exec(query, scheduleID)
	if err != nil {
		logs.ErrorLogger.Println(err)
		return err
	}
	return nil
}
