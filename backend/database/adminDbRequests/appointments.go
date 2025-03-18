package adminDbRequests

import (
	"awesomeProject/database"
	"awesomeProject/logs"
	"database/sql"
)

func GetAllAppointments(db *sql.DB) ([]database.AllAppointmentsData, error) {
	query := `
SELECT appointment_id, client_name, contact_number, appointment_date, time_start, time_end, service_id, employee_id, status
FROM appointments ORDER BY appointment_id`
	rows, err := db.Query(query)
	if err != nil {
		logs.ErrorLogger.Println(err)
		return nil, err
	}
	defer rows.Close()
	userDataList := []database.AllAppointmentsData{}
	for rows.Next() {
		userData := database.AllAppointmentsData{}
		err = rows.Scan(&userData.AppointmentID, &userData.UserName, &userData.UserContact, &userData.AppointmentDate, &userData.TimeStart, &userData.TimeEnd, &userData.ServiceID, &userData.EmployeeID, &userData.Status)
		if err != nil {
			logs.ErrorLogger.Println("Failed to scan row: %v", err)
			return nil, err
		}
		userDataList = append(userDataList, userData)
	}
	return userDataList, nil
}
func DeleteAppointment(db *sql.DB, appointmentID string) error {
	query := `DELETE FROM appointments WHERE appointment_id = $1`
	_, err := db.Exec(query, appointmentID)
	if err != nil {
		logs.ErrorLogger.Println(err)
		return err
	}
	return nil
}
