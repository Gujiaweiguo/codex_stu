package handlers

import (
	"cool-todolist/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	service *models.TodoService
}

func NewTodoHandler(service *models.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required,min=1,max=200"`
	Order int    `json:"order"`
}

type UpdateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
	Order     *int    `json:"order"`
}

type ReorderRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// GetAll 获取所有任务
func (h *TodoHandler) GetAll(c *gin.Context) {
	todos, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

// Create 创建任务
func (h *TodoHandler) Create(c *gin.Context) {
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.service.Create(req.Title, req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// Update 更新任务
func (h *TodoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	if req.Order != nil {
		todo.Order = *req.Order
	}

	updated, err := h.service.Update(todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// Toggle 切换任务完成状态
func (h *TodoHandler) Toggle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	todo, err := h.service.Toggle(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Delete 删除任务
func (h *TodoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Reorder 重新排序任务
func (h *TodoHandler) Reorder(c *gin.Context) {
	var req ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Reorder(req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
