package webHandlers

import (
	"awesomeProject/database/adminDbRequests"
	"awesomeProject/logs"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminGetAllEmployees(c *gin.Context, db *sql.DB) {
	employees, err := adminDbRequests.GetAllEmployees(db)
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, employees)
}
func AdminAddEmployee(c *gin.Context, db *sql.DB) {
	var employee struct {
		Name          string `json:"name"`
		Level         string `json:"level"`
		ContactNumber string `json:"contact_number"`
	}
	if err := c.BindJSON(&employee); err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := adminDbRequests.AddEmployee(db, employee.Name, employee.Level, employee.ContactNumber)
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add employee"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee added successfully"})
}
func AdminGetEmployeeLevels(c *gin.Context, db *sql.DB) {
	query := `SELECT DISTINCT level FROM employees`
	rows, err := db.Query(query)
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении уровней"})
		return
	}
	defer rows.Close()

	var levels []string
	for rows.Next() {
		var level string
		if err := rows.Scan(&level); err != nil {
			logs.ErrorLogger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных"})
			return
		}
		levels = append(levels, level)
	}

	c.JSON(http.StatusOK, gin.H{"levels": levels})

}
func AdminRedactEmployees(c *gin.Context, db *sql.DB) {
	var employee struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Level         string `json:"level"`
		ContactNumber string `json:"contact_number"`
	}
	if err := c.BindJSON(&employee); err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := adminDbRequests.RedactEmployee(db, employee.ID, employee.Name, employee.Level, employee.ContactNumber)
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee redacted successfully"})
}
func AdminDeleteEmployee(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	err := adminDbRequests.DeleteEmployee(db, id)
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}
