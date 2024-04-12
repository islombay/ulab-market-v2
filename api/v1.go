package api

import (
	handlersv1 "app/api/handlers/v1"
	"app/config"
	auth_lib "app/pkg/auth"
	"app/pkg/logs"
	"app/pkg/smtp"
	"app/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewV1(
	r *gin.RouterGroup,
	cfg *config.Config,
	store storage.StoreInterface,
	log logs.LoggerInterface,
	smtp smtp.SMTPInterface,
	cache storage.CacheInterface,
	filestore storage.FileStorageInterface,
) {
	// initialize v1 handler
	handler := handlersv1.NewHandler(log, cfg, store, smtp, cache, filestore)

	v1 := r.Group("/")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/login_admin", handler.LoginAdmin)
		auth.POST("/verify_code", handler.VerifyCode)

		//auth.POST("/register", handler.RegisterClient)
		//auth.POST("/change_password", handler.ChangePassword)
		//auth.POST("/request_code", handler.RequestCode)
	}

	owner := v1.Group("/owner")
	{
		ownerSuper := owner.Group("/").Use(handler.MiddlewareIsSuper())
		ownerSuper.POST("", handler.CreateOwner)
		ownerSuper.DELETE("/:id", handler.DeleteOwner)
		ownerSuper.PUT("", handler.ChangeOwner)
	}

	admin := v1.Group("/admin")
	{
		admin.POST("", handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionAddAdmin), handler.CreateAdmin)
		admin.DELETE("/:id", handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionDeleteAdmin), handler.DeleteAdmin)
		admin.PUT("", handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionEditAdmin), handler.ChangeAdmin)
	}

	category := v1.Group("/category")
	{
		category.POST("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCategoryAdd),
			handler.CreateCategory,
		)
		category.POST("/add_translation",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCategoryAdd),
			handler.AddCategoryTranslation,
		)

		category.POST("/change_image",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCategoryEdit),
			handler.ChangeCategoryImage,
		)

		category.PUT("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCategoryEdit),
			handler.ChangeCategory,
		)

		category.DELETE("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCategoryDelete),
			handler.DeleteCategory,
		)

		category.GET("/:id", handler.GetCategoryByID)
		category.GET("", handler.GetAllCategory)
	}

	brand := v1.Group("/brand")
	{
		brand.POST("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionBrandAdd),
			handler.CreateBrand,
		)

		brand.PUT("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionBrandEdit),
			handler.ChangeBrand,
		)

		brand.DELETE("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionBrandDelete),
			handler.DeleteBrand,
		)

		brand.GET("/:id", handler.GetBrandByID)
		brand.GET("", handler.GetAllBrand)
	}

	role := v1.Group("/roles")
	{
		role.GET("/role",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanViewRole),
			handler.GetAllRoles,
		) // handler to get all roles
		role.GET("/permission",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanViewRole),
			handler.GetAllPermissions,
		)
		role.POST("/attach",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanAttachRole),
			handler.AttachPermissionToRole,
		)
		role.DELETE("/attach",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanAttachRole),
			handler.DisAttachPermissionToRole,
		)
	}

	product := v1.Group("/product")
	{
		product.POST("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanAddProduct),
			handler.CreateProduct,
		)

		product.POST("/change_image",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanEditProduct),
			handler.ChangeProductMainImage,
		)

		product.DELETE("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanDeleteProduct),
			handler.DeleteProduct,
		)

		product.GET("", handler.GetAllProducts)
		product.GET("/:id", handler.GetProductByID)
	}

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ping": "pong"})
	})
}
