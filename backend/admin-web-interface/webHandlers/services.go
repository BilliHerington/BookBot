package webHandlers

import (
	"awesomeProject/database/adminDbRequests"
	"awesomeProject/logs"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminGetAllServices(c *gin.Context, db *sql.DB) {

	logs.DebugLogger.Println("db contain in services:", db)

	services, err := adminDbRequests.GetAllServices(db)
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}
func AdminAddService(c *gin.Context, db *sql.DB) {
	var service struct {
		Name        string  `json:"Name"`
		Duration    string  `json:"Duration"`
		DefPrice    float64 `json:"DefPrice"`
		ProPrice    float64 `json:"ProPrice"`
		Description string  `json:"Description"`
	}
	if err := c.BindJSON(&service); err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := adminDbRequests.AddService(db, service.Name, service.Duration, service.DefPrice, service.ProPrice, service.Description)
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add service"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service added successfully"})
}
func AdminRedactService(c *gin.Context, db *sql.DB) {
	var service struct {
		ID          int     `json:"ID"`
		Name        string  `json:"Name"`
		Duration    string  `json:"Duration"`
		DefPrice    float64 `json:"DefPrice"`
		ProPrice    float64 `json:"ProPrice"`
		Description string  `json:"Description"`
	}
	if err := c.BindJSON(&service); err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := adminDbRequests.RedactService(db, service.ID, service.Name, service.Duration, service.DefPrice, service.ProPrice, service.Description)
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service redacted successfully"})
}

func AdminDeleteService(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	err := adminDbRequests.DeleteService(db, id)
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}
