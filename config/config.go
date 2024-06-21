package config

import (
	"os"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host           string
	Port           string
	DBName         string
	SSLMode        string
	MigrationsPath string
}

type FileStorageConfig struct {
	URL string
}

type RedisConfig struct {
	Host string
	Port string
}

type MediaConfig struct {
	CategoryPhotoMaxSize int64
	ProductPhotoMaxSize  int64
	ProductPhotoMaxCount int

	ProductVideoMaxSize  int64
	ProductVideoMaxCount int
}

type SMTPConfig struct {
	Email SMTPEmailConfig
}

type SMTPEmailConfig struct {
	Host        string
	Port        string
	SenderEmail string
}

type serverConfig struct {
	Host   string
	Port   string
	Public string
}

type Config struct {
	Server      serverConfig
	DB          DBConfig
	Redis       RedisConfig
	SMTP        SMTPConfig
	FileStorage FileStorageConfig
	Env         string
	Media       MediaConfig
}

const (
	LocalMode  = "local"
	ProdMode   = "prod"
	DockerMode = "docker"

	filePath = "config"
)

func Load() Config {
	if _, err := os.Stat(filePath); err != nil {
		panic("config file not found" + err.Error())
	}

	viper.AddConfigPath(filePath)
	env := os.Getenv("ENV")
	switch env {
	case LocalMode:
		viper.SetConfigName(LocalMode)

	case ProdMode:
		viper.SetConfigName(ProdMode)
	case DockerMode:
		viper.SetConfigName(DockerMode)
	default:
		viper.SetConfigName(ProdMode)
		env = ProdMode
	}
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	cfg := Config{
		Server: serverConfig{
			Host:   viper.GetString("server.host"),
			Port:   viper.GetString("server.port"),
			Public: viper.GetString("server.public_host"),
		},
		DB: DBConfig{
			Host:           viper.GetString("db.host"),
			Port:           viper.GetString("db.port"),
			DBName:         viper.GetString("db.dbname"),
			SSLMode:        viper.GetString("db.sslmode"),
			MigrationsPath: viper.GetString("db.migrations"),
		},
		Redis: RedisConfig{
			Host: viper.GetString("redis.host"),
			Port: viper.GetString("redis.port"),
		},
		SMTP: SMTPConfig{
			Email: SMTPEmailConfig{
				Host:        viper.GetString("smtp.email.host"),
				Port:        viper.GetString("smtp.email.port"),
				SenderEmail: viper.GetString("smtp.email.sender"),
			},
		},
		Env: env,
		FileStorage: FileStorageConfig{
			URL: viper.GetString("filestorage.url"),
		},
		Media: MediaConfig{
			CategoryPhotoMaxSize: 1024 * 1024 * viper.GetInt64("media.category_photo_max_size"),
			ProductPhotoMaxSize:  1024 * 1024 * viper.GetInt64("media.product_photo_max_size"),
			ProductPhotoMaxCount: viper.GetInt("media.product_max_photo_count"),

			ProductVideoMaxCount: viper.GetInt("media.product_max_video_count"),
			ProductVideoMaxSize:  viper.GetInt64("media.product_video_max_size"),
		},
	}
	return cfg
}
