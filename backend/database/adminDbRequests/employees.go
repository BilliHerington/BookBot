package adminDbRequests

import (
	"awesomeProject/database"
	"awesomeProject/logs"
	"database/sql"
)

func GetAllEmployees(db *sql.DB) ([]database.AllEmployeesData, error) {
	query := `
SELECT * FROM employees ORDER BY employee_id`
	rows, err := db.Query(query)
	if err != nil {
		logs.ErrorLogger.Println(err)
		return nil, err
	}
	defer rows.Close()
	employeeDataList := []database.AllEmployeesData{}
	for rows.Next() {
		employeeData := database.AllEmployeesData{}
		err = rows.Scan(&employeeData.EmployeeID, &employeeData.Name, &employeeData.Level, &employeeData.ContactNumber)
		if err != nil {
			logs.ErrorLogger.Println(err)
			return nil, err
		}
		employeeDataList = append(employeeDataList, employeeData)
	}
	return employeeDataList, nil
}
func AddEmployee(db *sql.DB, employeeName string, level string, contactnumber string) error {
	query := `INSERT INTO employees (name, level , contact_number) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, employeeName, level, contactnumber)
	if err != nil {
		logs.ErrorLogger.Printf("Failed to execute query: %v", err)
		return err
	}
	return nil
}
func RedactEmployee(db *sql.DB, employeeID int, employeeName string, level string, contactnumber string) error {
	query := `UPDATE employees 
SET name = $2, level = $3, contact_number = $4
WHERE employee_id = $1`
	_, err := db.Exec(query, employeeID, employeeName, level, contactnumber)
	if err != nil {
		logs.ErrorLogger.Println(err)
		return err
	}
	return nil
}
func DeleteEmployee(db *sql.DB, employeeID string) error {
	query := `DELETE FROM employees WHERE employee_id = $1`
	_, err := db.Exec(query, employeeID)
	if err != nil {
		logs.ErrorLogger.Println(err)
		return err
	}
	return nil
}
