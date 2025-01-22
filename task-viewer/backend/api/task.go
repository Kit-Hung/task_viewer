package api

import (
	"net/http"
	"strconv"
	"task-viewer/service"

	"github.com/gin-gonic/gin"
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	Dependencies []int  `json:"dependencies"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// RegisterTaskRoutes 注册任务相关路由
func RegisterTaskRoutes(r *gin.Engine) {
	r.GET("/api/tasks", GetTasks)
	r.POST("/api/tasks", CreateTask)
	r.GET("/api/tasks/:id", GetTaskDetail)
	r.PUT("/api/tasks/:id", UpdateTask)
	r.DELETE("/api/tasks/:id", DeleteTask)
	r.POST("/api/tasks/:id/dependencies", AddDependency)
	r.DELETE("/api/tasks/:id/dependencies/:depId", RemoveDependency)
}

// GetTasks 获取所有任务
func GetTasks(c *gin.Context) {
	tasks, err := service.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dependencies, err := service.GetAllDependencies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":        tasks,
		"dependencies": dependencies,
	})
}

// CreateTask 创建任务
func CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := service.CreateTask(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 添加依赖关系
	for _, depID := range req.Dependencies {
		err = service.AddDependency(depID, task.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, task)
}

// GetTaskDetail 获取任务详情
func GetTaskDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := service.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	upstream, downstream, err := service.GetTaskDependencies(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task":                    task,
		"upstream_dependencies":   upstream,
		"downstream_dependencies": downstream,
	})
}

// UpdateTask 更新任务
func UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := service.UpdateTask(id, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask 删除任务
func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	err = service.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddDependency 添加依赖关系
func AddDependency(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		SourceID int `json:"source_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = service.AddDependency(req.SourceID, targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// RemoveDependency 移除依赖关系
func RemoveDependency(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	sourceID, err := strconv.Atoi(c.Param("depId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dependency ID"})
		return
	}

	err = service.RemoveDependency(sourceID, targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
