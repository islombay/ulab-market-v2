package api

import (
	"app/api/docs"
	"app/config"
	"app/pkg/logs"
	"app/pkg/smtp"
	"app/storage"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewApi(
	r *gin.Engine,
	cfg *config.Config,
	store storage.StoreInterface,
	log logs.LoggerInterface,
	smtp smtp.SMTPInterface,
	cache storage.CacheInterface,
	filestore storage.FileStorageInterface,
) {

	// @securityDefinitions.apikey ApiKeyAuth
	// @in header
	// @name Authorization

	docs.SwaggerInfo.Title = "E-commerce project v2-1"
	docs.SwaggerInfo.Description = "This is a sample server e-commerce server."
	docs.SwaggerInfo.Version = "1.0"

	r.Use(customCORSMiddleware())
	api := r.Group("/api")
	NewV1(api, cfg, store, log, smtp, cache, filestore)

	cfg.Env = os.Getenv("ENV")
	if cfg.Env == config.LocalMode {
		docs.SwaggerInfo.Host = "http://localhost:8123"
	} else if cfg.Env == config.ProdMode {
		docs.SwaggerInfo.Host = "https://ulab-market-backend.onrender.com/"
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		//ginSwagger.URL("swagger/doc.json"),
	))
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Acces-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Header("Access-Control-Allow-Headers", "Platform-Id, Content-Type, Content-Length, Accept-Encoding, X-CSF-TOKEN, Authorization, Cache-Control")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
