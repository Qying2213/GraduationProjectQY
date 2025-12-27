package main

import (
	"log"
	"talent-service/handlers"
	"talent-service/models"

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

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	talentHandler := handlers.NewTalentHandler(db)

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
