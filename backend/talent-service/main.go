package main

import (
	"log"
	"talent-service/handlers"
	"talent-service/models"

	"common/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=qinyang dbname=talent_platform port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	if err := db.AutoMigrate(&models.Talent{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.SimpleOperationLog("talent-service"))

	talentHandler := handlers.NewTalentHandler(db)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "talent-service", "port": 8086})
	})

	api := r.Group("/api/v1/talents")
	{
		api.POST("", talentHandler.CreateTalent)
		api.GET("", talentHandler.ListTalents)
		api.GET("/search", talentHandler.SearchTalents)
		api.GET("/:id", talentHandler.GetTalent)
		api.PUT("/:id", talentHandler.UpdateTalent)
		api.DELETE("/:id", talentHandler.DeleteTalent)
	}

	log.Println("Talent service is running on :8086")
	if err := r.Run(":8086"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
