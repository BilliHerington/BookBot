package adminDbRequests

import (
	"awesomeProject/database"
	"awesomeProject/logs"
	"database/sql"
)

func GetAllServices(db *sql.DB) ([]database.AllServicesData, error) {
	query := `
SELECT * FROM services ORDER BY service_id`
	rows, err := db.Query(query)
	if err != nil {
		logs.ErrorLogger.Println(err)
		return nil, err
	}
	defer rows.Close()
	serviceDataList := []database.AllServicesData{}
	for rows.Next() {
		serviceData := database.AllServicesData{}
		err = rows.Scan(&serviceData.ServiceID, &serviceData.Name, &serviceData.Duration, &serviceData.DefaultPrice, &serviceData.ProPrice, &serviceData.Description)
		if err != nil {
			logs.ErrorLogger.Println("Failed to scan row: %v", err)
			return nil, err
		}
		serviceDataList = append(serviceDataList, serviceData)
	}
	return serviceDataList, nil
}
func AddService(db *sql.DB, name string, duration string, defprice float64, proprice float64, description string) error {

	query := `INSERT INTO services (name, duration, default_price, pro_price, description) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, name, duration, defprice, proprice, description)
	if err != nil {
		logs.ErrorLogger.Printf("Failed to execute query: %v", err)
		return err
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
		logs.ErrorLogger.Println(err)
		return err
	}
	return nil
}
