package api

import (
	handlersv1 "app/api/handlers/v1"
	"app/config"
	auth_lib "app/pkg/auth"
	"app/pkg/logs"
	"app/pkg/smtp"
	"app/service"
	"app/storage"
	"fmt"
	"net/http"
	"os"

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
	service service.IServiceManager,
) {
	// initialize v1 handler
	handler := handlersv1.NewHandler(log, cfg, store, smtp, cache, filestore, service)

	v1 := r.Group("/")

	go service.Notify().Courier.HandleNotifications()

	super := v1.Group("/super")
	{
		super.GET("migrate-down",
			handler.MiddlewareIsSuper(),
			handler.SuperMigrateDown,
		)
	}

	auth := v1.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/login_admin", handler.LoginAdmin)
		auth.POST("/login_courier", handler.LoginCourier)
		auth.POST("/verify_code", handler.VerifyCode)

		auth.GET("/tgbot", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"url": os.Getenv("auth_tg_bot")})
		})
	}

	client := v1.Group("/client")
	{
		client.GET("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanSeeClients),
			handler.GetClientList,
		)

		client.GET("/getme",
			handler.MiddlewareIsClient(),
			handler.ClientGetMe,
		)

		client.PUT("/me",
			handler.MiddlewareIsClient(),
			handler.ClientUpdate,
		)
	}

	income := v1.Group("/income")
	{
		income.POST("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionIncomeAdd),
			handler.CreateIncome,
		)
		income.GET("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionIncomeSee),
			handler.GetIncomeList,
		)
		income.GET("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionIncomeSee),
			handler.GetIncomeByID,
		)
	}

	branches := v1.Group("/branch")
	{
		branches.POST("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionAddBranch),
			handler.AddBranch,
		)

		branches.PUT("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionEditBranch),
			handler.ChangeBranch,
		)

		branches.DELETE("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionDeleteBranch),
			handler.DeleteBranch,
		)

		branches.GET("/:id", handler.GetBranchByID)
		branches.GET("", handler.GetAllBranches)
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

	courier := v1.Group("/courier")
	{
		courier.POST("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionAddCourier),
			handler.CreateCourier,
		)
		courier.DELETE("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionDeleteCourier),
			handler.DeleteCourier,
		)

		courier.GET("/order/ws",
			handler.MiddlewareIsCourier(),
			handler.CourierOrdersRealTimeConnection,
		)
	}

	category := v1.Group("/category")
	{
		category.POST("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCategoryAdd),
			handler.CreateCategory,
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
		category.GET("/:id/brand", handler.GetCategoryBrands)

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

		role.POST("/role",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanAddRole),
			handler.CreateNewRole,
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

		product.PUT("/change_price",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanEditProduct),
			handler.ChangeProductPrice,
		)

		product.DELETE("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanDeleteProduct),
			handler.DeleteProduct,
		)

		product.POST("/add_image_files",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanEditProduct),
			handler.AddProductImageFiles,
		)

		product.POST("/add_video_files",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanEditProduct),
			handler.AddProductVideoFiles,
		)

		product.PUT("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanEditProduct),
			handler.EditProduct,
		)

		product.GET("", handler.GetAllProducts)
		product.GET("_pagin", handler.GetAllProductsPagination)
		product.GET("/:id", handler.GetProductByID)
		product.GET("/admin/:id", handler.MiddlewareIsStaff(), handler.GetProductByIDAdmin)
	}

	basket := v1.Group("/basket")
	{
		basket.POST("",
			handler.MiddlewareIsClient(),
			handler.AddToBasket,
		)

		basket.GET("",
			handler.MiddlewareIsClient(),
			handler.GetBasket,
		)

		basket.PUT("",
			handler.MiddlewareIsClient(),
			handler.ChangeBasket,
		)

		basket.DELETE("",
			handler.MiddlewareIsClient(),
			handler.DeleteFromBasket,
		)

		basket.DELETE("/all",
			handler.MiddlewareIsClient(),
			handler.DeleteAllBasket,
		)
	}

	iconsList := v1.Group("/icon")
	{
		iconsList.POST("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionAddIconToList),
			handler.AddIconToList,
		)

		iconsList.DELETE("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionDeleteIconToList),
			handler.DeleteIcon,
		)

		iconsList.GET("", handler.GetIconsAll)
		iconsList.GET("/:id", handler.GetIconByID)
	}

	order := v1.Group("/order")
	{
		order.POST("",
			handler.MiddlewareIsClient(),
			handler.CreateOrder,
		)

		order.POST("/finish/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanFinishOrder),
			handler.OrderFinish,
		)

		order.POST("/cancel/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanCancelOrder),
			handler.OrderFinish,
		)

		order.GET("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanSeeOrderByID),
			handler.GetOrderByID,
		)

		order.GET("",
			handler.MiddlewareIsStaff(),
			handler.GetOrderAll,
		)

		order.GET("/product/:id", handler.GetOrderProduct)
		order.GET("/product", handler.GetOrderProductAll)

		order.GET("/new",
			handler.MiddlewareIsStaff(),
			handler.GetNewOrdersList,
		)

		order.GET("/courier",
			handler.MiddlewareIsCourier(),
			handler.GetAvailableOrdersCourier,
		)

		order.GET("/picked/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCanMakeOrderPicked),
			handler.OrderPicked,
		)

		order.GET("/picked_deliver/:id",
			handler.MiddlewareIsCourier(),
			handler.OrderMarkPickedByCourier,
		)

		order.GET("/myorders",
			handler.MiddlewareIsClient(),
			handler.ClientOrders,
		)

		order.GET("/delivered/:id",
			handler.MiddlewareIsCourier(),
			handler.OrderDelivered,
		)

		order.GET("/courier/myorders",
			handler.MiddlewareIsCourier(),
			handler.CourierOrdersGetAll,
		)

		order.POST("/courier/start",
			handler.MiddlewareIsCourier(),
			handler.CourierStartDeliver,
		)
	}

	storeTable := v1.Group("/storage")
	{
		storeTable.POST("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionCreateStore),
			handler.CreateStorage)

		storeTable.GET("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionGetStoreByID),
			handler.GetStorageByID)

		storeTable.GET("",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionGetStoreList),
			handler.GetStorageList)

		storeTable.PUT("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionUpdateStore),
			handler.UpdateStorage)

		storeTable.DELETE("/:id",
			handler.MiddlewareStaffPermissionCheck(auth_lib.PermissionDeleteStore),
			handler.DeleteStorage)
	}

	favourite := v1.Group("/favourite")
	{
		favourite.POST("",
			handler.MiddlewareIsClient(),
			handler.AddToFavourite,
		)

		favourite.DELETE("/:productID",
			handler.MiddlewareIsClient(),
			handler.DeleteFromFavourite,
		)

		favourite.GET("",
			handler.MiddlewareIsClient(),
			handler.GetAllFavourite,
		)
	}

	v1.GET("/ping", func(c *gin.Context) {
		fmt.Println("entered")
		c.JSON(http.StatusOK, gin.H{"ping": "pong"})
	})
}
