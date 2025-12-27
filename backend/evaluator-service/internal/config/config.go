package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServerCfg struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (s ServerCfg) Address() string {
	if s.Port == 0 {
		return ":8080"
	}
	if s.Host == "" {
		return fmt.Sprintf(":%d", s.Port)
	}
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type DBCfg struct {
	Driver  string `mapstructure:"driver"` // sqlite|mysql|postgres
	DSN     string `mapstructure:"dsn"`
	MaxOpen int    `mapstructure:"max_open"`
	MaxIdle int    `mapstructure:"max_idle"`
}

type StorageCfg struct {
	BaseDir   string `mapstructure:"base_dir"`
	Reports   string `mapstructure:"reports"`
	Resumes   string `mapstructure:"resumes"`
	TempDir   string `mapstructure:"temp_dir"`
	StaticDir string `mapstructure:"static_dir"`
}

type AICfg struct {
	Provider string `mapstructure:"provider"`
	APIKey   string `mapstructure:"api_key"`
	Model    string `mapstructure:"model"`
}

type ExportCfg struct {
	PDFEngine string `mapstructure:"pdf_engine"` // wkhtmltopdf|chromedp
}

type BatchCfg struct {
	Concurrency int `mapstructure:"concurrency"`
	MaxUploadMB int `mapstructure:"max_upload_mb"`
}

type CozeCfg struct {
	BaseURL    string `mapstructure:"base_url"`
	Token      string `mapstructure:"token"`
	WorkflowID string `mapstructure:"workflow_id"`
}

type CredentialsCfg struct {
	EncKey string `mapstructure:"enc_key"`
}

type PythonCfg struct {
	Path string `mapstructure:"path"` // Python 解释器路径
}

type GraduateCfg struct {
	APIUrl string `mapstructure:"api_url"` // 毕业设计后台 API 地址
}

type Config struct {
	Server      ServerCfg      `mapstructure:"server"`
	DB          DBCfg          `mapstructure:"db"`
	Storage     StorageCfg     `mapstructure:"storage"`
	AI          AICfg          `mapstructure:"ai"`
	Export      ExportCfg      `mapstructure:"export"`
	Batch       BatchCfg       `mapstructure:"batch"`
	Coze        CozeCfg        `mapstructure:"coze"`
	Credentials CredentialsCfg `mapstructure:"credentials"`
	Python      PythonCfg      `mapstructure:"python"`
	Graduate    GraduateCfg    `mapstructure:"graduate"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.SetEnvPrefix("RESUME")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)
	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	// 服务器配置
	v.SetDefault("server.port", 8090)

	// 数据库配置
	v.SetDefault("db.driver", "sqlite")
	v.SetDefault("db.dsn", "file:data/app.db?cache=shared&mode=rwc")

	// 存储配置
	v.SetDefault("storage.base_dir", "data")
	v.SetDefault("storage.reports", "reports")
	v.SetDefault("storage.resumes", "resumes")
	v.SetDefault("storage.temp_dir", "tmp")
	v.SetDefault("storage.static_dir", "static")

	// 导出配置
	v.SetDefault("export.pdf_engine", "wkhtmltopdf")

	// 批量处理配置
	v.SetDefault("batch.concurrency", 5)
	v.SetDefault("batch.max_upload_mb", 50)

	// AI 配置
	v.SetDefault("ai.model", "")

	// Coze 配置
	v.SetDefault("coze.base_url", "")
	v.SetDefault("coze.token", "")
	v.SetDefault("coze.workflow_id", "")

	// 凭据加密密钥 -
	v.SetDefault("credentials.enc_key", "")

	// Python 配置
	v.SetDefault("python.path", "")

	// 毕业设计后台配置
	v.SetDefault("graduate.api_url", "http://localhost:8084")
}
