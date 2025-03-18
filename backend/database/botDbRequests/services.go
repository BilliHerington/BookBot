package botDbRequests

import (
	"awesomeProject/logs"
	"database/sql"
)

func GetServices(db *sql.DB) []string {
	//logs.DebugLogger.Printf("Checking DB before query: %+v", db)

	rows, err := db.Query("SELECT name FROM services")

	defer rows.Close()
	var serviceSlice []string

	for rows.Next() {
		service_name := ""
		err = rows.Scan(&service_name)
		if err != nil {
			logs.ErrorLogger.Println(err)
		}
		serviceSlice = append(serviceSlice, service_name)
	}

	return serviceSlice
}
