package main
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

type User struct {
	ID uint `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type Todo struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Title string `gorm:"not null" json:"title"`
	Completed bool `gorm:"default:false" json:"completed"`
	UserID uint `gorm:"not null" json:"user_id"`
}

type CreateTodoInput struct {
	Title string `json:"title" binding:"required"`
}

type UpdateTodoInput struct {
	Title string `json:"title"`
	Completed *bool `json:"completed"`
}

var db *gorm.DB

func getJWTKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("Critical error of safety: variable JWT_SECRET is not specified")
	}
	return []byte(secret)
}

func main() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

			var err error
			for i := 0; i < 5; i++ {
				db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
				if err == nil {
					break
				}
				log.Println("Waiting database...")
				time.Sleep(3 * time.Second)
			}
			if err != nil {
				log.Fatalf("Error of connecting to database: %v", err)
			}
			err = db.AutoMigrate(&User{}, &Todo{})
			if err != nil {
				log.Fatalf("Error of automigration: %v", err)
			}

			r := gin.Default()
			r.POST("/register", Register)
			r.POST("/login",Login)

			protected := r.Group("/api")
			protected.Use(AuthMiddleware())
			{
				protected.POST("/todos", CreateTodo)
				protected.GET("/todos", GetTodos)
				protected.PUT("/todos/:id", UpdateTodo)
				protected.DELETE("/todos/:id", DeleteTodo)
			}

			port := os.Getenv("APP_PORT")
			if port == "" {
				port = "8080"
			}

			log.Printf("Запуск сервера на порту %s", port)
			err = r.Run(":" + port)
			if err != nil {
				log.Fatalf("Ошибка запуска сервера: %v", err)
			}
}

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect format for data"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error of hashing password"})
		return
	}
	user := User{Username: input.Username, Password: string(hashedPassword)}
	err = db.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registrated successfully"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect format of data"})
		return
	}
	var user User
	err = db.Where("username = ?", input.Username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect login or password"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect login or password"})
		return
	}

	token:= jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(getJWTKey())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error of generating token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization") // where token lay in the http request
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not present"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect format of header Authorization"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return getJWTKey(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not valid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect structure of token"})
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Identificator of user is not present"})
			c.Abort()
			return
		}

		c.Set("userID", uint(userIDFloat))
		c.Next()
	}
}

func CreateTodo(c *gin.Context) {
	var input CreateTodoInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User is not identificated"})
		return
	}

	todo := Todo{
		Title: input.Title,
		UserID: userID.(uint),
	}

	err = db.Create(&todo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error of saving task"})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

func GetTodos(c *gin.Context) {
	var todos []Todo
	userID, _ := c.Get("userID")
	db.Where("user_id = ?", userID.(uint)).Find(&todos)
	c.JSON(http.StatusOK, todos)
}

func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	var todo Todo

	err := db.Where("id = ? AND user_id = ?", id, userID.(uint)).First(&todo).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task is not found"})
		return
	}

	var input UpdateTodoInput
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != "" {
		todo.Title = input.Title
	}
	if input.Completed != nil {
		todo.Completed = *input.Completed
	}

	err = db.Save(&todo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error of updating task"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	result := db.Where("id = ? AND user_id = ?", id, userID.(uint)).Delete(&Todo{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error of deleting"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task isnt found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task successfully deleted"})
}