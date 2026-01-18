package main

import (
	"cool-todolist/database"
	"cool-todolist/handlers"
	"cool-todolist/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// 初始化服务和处理器
	todoService := models.NewTodoService(database.DB)
	todoHandler := handlers.NewTodoHandler(todoService)

	// 注册路由
	api := r.Group("/api")
	{
		todos := api.Group("/todos")
		{
			todos.GET("", todoHandler.GetAll)
			todos.POST("", todoHandler.Create)
			todos.PUT("/:id", todoHandler.Update)
			todos.PATCH("/:id/toggle", todoHandler.Toggle)
			todos.DELETE("/:id", todoHandler.Delete)
			todos.POST("/reorder", todoHandler.Reorder)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
