package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"taskflow/internal/domain/entities"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoHandler struct {
	todoUseCase service.TodoService
}

func NewTodoHandler(todoUseCase service.TodoService) *TodoHandler {
	return &TodoHandler{
		todoUseCase: todoUseCase,
	}
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req entities.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	todo, err := h.todoUseCase.CreateTodo(c.Request.Context(), &req)
	if err != nil {
		// Check for validation errors
		if errors.Is(err, errors.New("description is required")) ||
			errors.Is(err, errors.New("due date must be in the future")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	todo, err := h.todoUseCase.GetTodo(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) ListTodos(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	todos, err := h.todoUseCase.ListTodos(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(todos),
		},
	})
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var req entities.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	todo, err := h.todoUseCase.UpdateTodo(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, errors.New("todo not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	err = h.todoUseCase.DeleteTodo(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("todo not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
