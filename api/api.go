package api

import (
	_ "app/api/docs"
	"app/config"
	"app/pkg/logs"
	"app/pkg/smtp"
	"app/storage"

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

	// @title E-commerce project v2
	// @host localhost:8123
	// @securityDefinitions.apikey ApiKeyAuth
	// @in header
	// @name Authorization

	r.Use(customCORSMiddleware())
	api := r.Group("/api")
	NewV1(api, cfg, store, log, smtp, cache, filestore)

	url := ginSwagger.URL("swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
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
