package botDbRequests

import (
	"awesomeProject/logs"
	"database/sql"
	"fmt"
)

func GetLevelAndCost(db *sql.DB, service string) ([]string, []string) {
	query := `
SELECT 
    e.level,
    CASE 
        WHEN e.level = 'профессионал' THEN s.pro_price
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
		logs.ErrorLogger.Println(err)
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
			logs.ErrorLogger.Println(err)
		}
		resultString := fmt.Sprintf("%s %0.0f", employee_level, serviceCostByLevel)
		levelAndCostSlice = append(levelAndCostSlice, resultString)
		employeeLevelSlice = append(employeeLevelSlice, employee_level)
	}

	return levelAndCostSlice, employeeLevelSlice
}
