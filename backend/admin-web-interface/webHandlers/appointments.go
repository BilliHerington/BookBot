package webHandlers

import (
	"awesomeProject/database/adminDbRequests"
	"awesomeProject/logs"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminGetAllAppointments(c *gin.Context, db *sql.DB) {
	appointments, err := adminDbRequests.GetAllAppointments(db)
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, appointments)
}
func AdminDeleteAppointment(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	err := adminDbRequests.DeleteAppointment(db, id)
	if err != nil {
		logs.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Appointment deleted successfully"})
}
