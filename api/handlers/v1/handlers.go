package handlersv1

import (
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/config"
	"app/pkg/logs"
	"app/pkg/smtp"
	"app/service"
	"app/storage"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	log       logs.LoggerInterface
	cfg       *config.Config
	storage   storage.StoreInterface
	smtp      smtp.SMTPInterface
	cache     storage.CacheInterface
	filestore storage.FileStorageInterface

	service service.IServiceManager
}

func NewHandler(
	log logs.LoggerInterface,
	cfg *config.Config,
	storage storage.StoreInterface,
	smtp smtp.SMTPInterface,
	cache storage.CacheInterface,
	filestore storage.FileStorageInterface,
	service service.IServiceManager,
) *Handlers {
	return &Handlers{
		log:       log,
		cfg:       cfg,
		storage:   storage,
		smtp:      smtp,
		cache:     cache,
		filestore: filestore,
		service:   service,
	}
}

func (v1 *Handlers) error(c *gin.Context, status status.Status) {
	switch code := status.Code; {
	case code >= 500:
		v1.log.Error("[-Server Error-]:",
			logs.Int("code", status.Code),
			logs.String("status", status.Message),
		)
	case code >= 400:
		v1.log.Error("[-Response-]:",
			logs.Int("code", status.Code),
			logs.String("status", status.Message),
		)
	}
	c.AbortWithStatusJSON(status.Code, models_v1.Response{
		Message: status.Message,
		Code:    status.Code,
	})
}

func (v1 *Handlers) response(c *gin.Context, code int, data interface{}) {
	v1.log.Info("[-Response-]:",
		logs.Int("code", code),
		logs.Any("url", c.Request.URL),
	)

	c.JSON(code, data)
}
