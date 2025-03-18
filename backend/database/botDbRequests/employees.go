package botDbRequests

import (
	"awesomeProject/logs"
	"database/sql"
)

func GetEmployeeName(db *sql.DB, employeeID int) string {
	query := `
	SELECT name FROM employees WHERE employee_id = $1`
	rows, err := db.Query(query, employeeID)
	defer rows.Close()
	employeeName := ""
	for rows.Next() {
		err = rows.Scan(&employeeName)
		if err != nil {
			logs.ErrorLogger.Println(err)
		}
	}
	return employeeName
}
