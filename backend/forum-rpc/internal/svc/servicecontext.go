package svc

import (
	"fmt"

	"common/database"
	"forum-rpc/internal/config"
	"forum-rpc/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		c.Database.Host,
		c.Database.User,
		c.Database.Name,
		c.Database.Port,
		c.Database.SSLMode,
	)
	if c.Database.Password != "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
			c.Database.Host,
			c.Database.User,
			c.Database.Password,
			c.Database.Name,
			c.Database.Port,
			c.Database.SSLMode,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	if err := database.InitRedis(database.RedisConfig{
		Addr:     fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	}); err != nil {
		logx.Infof("failed to connect redis (forum cache degraded): %v", err)
	}

	if err := db.AutoMigrate(&model.ForumBoard{}, &model.ForumPost{}); err != nil {
		panic(fmt.Sprintf("failed to migrate forum tables: %v", err))
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
