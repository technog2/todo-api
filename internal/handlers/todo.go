package handlers

import (
	"net/http"
	"strconv"
	"todo-api/internal/models"
	"todo-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TodoHandler struct{}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{}
}

// create new todo
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo.UserID = userID

	db := repository.GetDB()
	if err := db.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// get all todos
func (h *TodoHandler) GetTodos(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var todos []models.Todo
	db := repository.GetDB()

	// query parameters
	completed := c.Query("completed")
	query := db.Where("user_id = ?", userID)

	if completed != "" {
		isCompleted, _ := strconv.ParseBool(completed)
		query = query.Where("completed = ?", isCompleted)
	}

	if err := query.Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// update todo
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userID := getUserIDFromContext(c)
	todoID, _ := strconv.Atoi(c.Param("id"))

	var todo models.Todo
	db := repository.GetDB()

	// find todo
	if err := db.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// delete todo
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userID := getUserIDFromContext(c)
	todoID, _ := strconv.Atoi(c.Param("id"))

	db := repository.GetDB()
	result := db.Where("id = ? AND user_id = ?", todoID, userID).Delete(&models.Todo{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}

// help function to get user_id from token
func getUserIDFromContext(c *gin.Context) uint {
	claims, _ := c.Get("user")
	userClaims := claims.(jwt.MapClaims)
	return uint(userClaims["user_id"].(float64))
}
