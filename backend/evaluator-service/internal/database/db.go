package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"evaluator-service/internal/config"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.Config, log *logging.Logger) error {
	if err := ensureDirs(cfg); err != nil {
		return err
	}
	var (
		dialector gorm.Dialector
	)
	// only sqlite for now
	dialector = sqlite.Open(cfg.DB.DSN)
	gormCfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}
	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return err
	}
	// connection pool
	sqlDB, _ := db.DB()
	setupPool(sqlDB, cfg)
	if err := migrate(db); err != nil {
		return err
	}
	DB = db
	log.Info("db initialized", logging.KV("driver", cfg.DB.Driver))
	return nil
}

func setupPool(sqlDB *sql.DB, cfg *config.Config) {
	if sqlDB == nil {
		return
	}
	if cfg.DB.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(cfg.DB.MaxOpen)
	}
	if cfg.DB.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(cfg.DB.MaxIdle)
	}
	sqlDB.SetConnMaxLifetime(60 * time.Minute)
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Candidate{}, &models.Credential{}, &models.DingTalkConfig{}, &models.Position{})
}

func ensureDirs(cfg *config.Config) error {
	base := cfg.Storage.BaseDir
	paths := []string{base, filepath.Join(base, cfg.Storage.Reports), filepath.Join(base, cfg.Storage.Resumes), filepath.Join(base, cfg.Storage.TempDir)}
	for _, p := range paths {
		if err := os.MkdirAll(p, 0755); err != nil {
			return fmt.Errorf("create dir %s: %w", p, err)
		}
	}
	return nil
}
