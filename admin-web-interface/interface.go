package admin_web_interface

import (
	"awesomeProject/logs"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

var jwtKey = []byte("your_secret_key")

type ScheduleJSON struct {
	ID         int    `json:"id"`
	EmployeeID int    `json:"employee_id"`
	WorkDate   string `json:"work_date"`
	TimeStart  string `json:"time_start"`
	TimeEnd    string `json:"time_end"`
}

// Структура для хранения учетных данных
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Загрузка учетных данных из файла
func loadCredentials() (Credentials, error) {
	file, err := os.Open("admin-web-interface/configAdmin.json")
	if err != nil {
		logs.ErrorLogger.Printf("Error opening file: %v", err)
		return Credentials{}, err
	}
	defer file.Close()

	var creds Credentials
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&creds); err != nil {
		logs.ErrorLogger.Printf("Error decoding JSON: %v", err)
		return Credentials{}, err
	}

	return creds, nil
}

func Login(c *gin.Context) {
	var inputCredentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&inputCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Загрузите учетные данные из файла
	creds, err := loadCredentials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load credentials"})
		return
	}

	// Проверьте учетные данные
	if inputCredentials.Username != creds.Username || inputCredentials.Password != creds.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": inputCredentials.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "request does not contain an access token"})
			c.Abort()
			return
		}

		// Ожидаем формат "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // если не было префикса Bearer
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		username := claims["username"].(string)
		c.Set("username", username)

		c.Next()
	}
}
