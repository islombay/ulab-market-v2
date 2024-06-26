package auth_lib

import (
	"app/api/models"

	"github.com/google/uuid"
)

const (
	VerificationEmail = "email"
	VerificationPhone = "phone_number"
	CodeLength        = 6
)

var (
	RoleCourier = models.RoleModel{
		ID:          uuid.New().String(),
		Name:        "courier",
		Description: models.GetStringAddress("courier model. delivery"),
		Permissions: []models.PermissionModel{
			PermissionCanSeeOrderByID,
		},
	}
	RoleOwner = models.RoleModel{
		ID:          uuid.New().String(),
		Name:        "owner",
		Description: models.GetStringAddress("owner is the ceo of business"),
		Permissions: []models.PermissionModel{
			PermissionCanEditRole,

			PermissionAddIconToList,
			PermissionDeleteIconToList,

			PermissionCanEditProduct,
			PermissionCanAddProduct,
			PermissionCanDeleteProduct,

			PermissionBrandAdd,
			PermissionCanSeeOrderByID,
			PermissionBrandEdit,
			PermissionBrandDelete,

			PermissionCategoryAdd,
			PermissionCategoryEdit,
			PermissionCategoryDelete,

			PermissionCanSeeClients,

			PermissionIncomeAdd,
			PermissionIncomeSee,

			PermissionAddAdmin,
			PermissionDeleteAdmin,
			PermissionEditAdmin,

			PermissionAddCourier,
			PermissionEditCourier,
			PermissionDeleteCourier,

			PermissionCanViewRole,
			PermissionCanAttachRole,

			PermissionAddBranch,
			PermissionEditBranch,
			PermissionDeleteBranch,

			PermissionCanFinishOrder,
			PermissionCanCancelOrder,

			PermissionCreateStore,
			PermissionGetStoreByID,
			PermissionGetStoreList,
			PermissionUpdateStore,
			PermissionDeleteStore,
		},
	}
	RoleClient = models.RoleModel{
		ID:          uuid.New().String(),
		Name:        "client",
		Description: models.GetStringAddress("client model for any other role"),
		Permissions: []models.PermissionModel{
			PermissionAddToBasket,
			PermissionCanAddRole,
		},
	}
	RoleAdmin = models.RoleModel{
		ID:          uuid.New().String(),
		Name:        "admin",
		Description: models.GetStringAddress("admin is helper for owner"),
		Permissions: []models.PermissionModel{
			PermissionCanAddProduct,
			PermissionCanEditProduct,
			PermissionCanDeleteProduct,

			PermissionCategoryAdd,
			PermissionCategoryEdit,
			PermissionCategoryDelete,

			PermissionAddCourier,
			PermissionEditCourier,
			PermissionDeleteCourier,
			PermissionCanSeeClients,
			PermissionCanMakeOrderPicked,

			PermissionCreateStore,
			PermissionGetStoreByID,
			PermissionGetStoreList,
			PermissionUpdateStore,
			PermissionDeleteStore,

			PermissionCanSeeOrderByID,
		},
	}
	RoleSuper = models.RoleModel{
		ID:          uuid.New().String(),
		Name:        "super",
		Description: models.GetStringAddress("super administrator, in other words root"),
		Permissions: []models.PermissionModel{
			PermissionCanMigrateDown,

			PermissionCanAddRole,
			PermissionCanEditRole,
			PermissionCanDeleteRole,
			PermissionCanViewRole,
			PermissionCanAttachRole,

			PermissionAddIconToList,
			PermissionDeleteIconToList,
		},
	}
)

var (
	PermissionCanSeeClients = models.PermissionModel{ID: uuid.NewString(), Name: "can_see_clients"}
)

var (
	PermissionCanMakeOrderPicked = models.PermissionModel{ID: uuid.NewString(), Name: "can_make_order_picked"}
	PermissionCanFinishOrder     = models.PermissionModel{ID: uuid.NewString(), Name: "finish_order"}
	PermissionCanCancelOrder     = models.PermissionModel{ID: uuid.NewString(), Name: "cancel_order"}

	PermissionCanSeeOrderByID = models.PermissionModel{ID: uuid.NewString(), Name: "can_see_order_by_id"}
)

var (
	PermissionIncomeAdd = models.PermissionModel{ID: uuid.NewString(), Name: "can_add_income"}
	PermissionIncomeSee = models.PermissionModel{ID: uuid.NewString(), Name: "can_income_see"}
)

var (
	PermissionAddToBasket      = models.PermissionModel{ID: uuid.NewString(), Name: "can_add_to_basket"}
	PermissionRemoveFromBasket = models.PermissionModel{ID: uuid.NewString(), Name: "can_remove_from_basket"}
)

var (
	PermissionBrandAdd    = models.PermissionModel{ID: uuid.NewString(), Name: "can_add_brand"}
	PermissionBrandEdit   = models.PermissionModel{ID: uuid.NewString(), Name: "can_edit_brand"}
	PermissionBrandDelete = models.PermissionModel{ID: uuid.NewString(), Name: "can_delete_brand"}
)

// Role permissions
var (
	PermissionCanMigrateDown = models.PermissionModel{ID: uuid.New().String(), Name: "can_migrate_down"}

	PermissionCanAddRole    = models.PermissionModel{ID: uuid.New().String(), Name: "can_add_role"}
	PermissionCanEditRole   = models.PermissionModel{ID: uuid.New().String(), Name: "can_change_role"}
	PermissionCanDeleteRole = models.PermissionModel{ID: uuid.New().String(), Name: "can_delete_role"}
	PermissionCanViewRole   = models.PermissionModel{ID: uuid.NewString(), Name: "can_view_role"}

	PermissionCanAttachRole = models.PermissionModel{ID: uuid.NewString(), Name: "can_attach_role"}
)

// Product permissions
var (
	PermissionCanAddProduct    = models.PermissionModel{ID: uuid.New().String(), Name: "can_add_product"}
	PermissionCanEditProduct   = models.PermissionModel{ID: uuid.New().String(), Name: "can_edit_product"}
	PermissionCanDeleteProduct = models.PermissionModel{ID: uuid.New().String(), Name: "can_delete_product"}
)

// Category permissions
var (
	PermissionCategoryAdd    = models.PermissionModel{ID: uuid.New().String(), Name: "add_category"}
	PermissionCategoryEdit   = models.PermissionModel{ID: uuid.New().String(), Name: "edit_category"}
	PermissionCategoryDelete = models.PermissionModel{ID: uuid.New().String(), Name: "delete_category"}
)

var (
	PermissionAddBranch    = models.PermissionModel{ID: uuid.NewString(), Name: "add_branch"}
	PermissionEditBranch   = models.PermissionModel{ID: uuid.NewString(), Name: "edit_branch"}
	PermissionDeleteBranch = models.PermissionModel{ID: uuid.NewString(), Name: "delete_branch"}
)

var (
	PermissionAddAdmin    = models.PermissionModel{ID: uuid.New().String(), Name: "add_admin"}
	PermissionEditAdmin   = models.PermissionModel{ID: uuid.New().String(), Name: "edit_admin"}
	PermissionDeleteAdmin = models.PermissionModel{ID: uuid.New().String(), Name: "delete_admin"}
)

var (
	PermissionAddCourier    = models.PermissionModel{ID: uuid.New().String(), Name: "add_courier"}
	PermissionEditCourier   = models.PermissionModel{ID: uuid.New().String(), Name: "edit_courier"}
	PermissionDeleteCourier = models.PermissionModel{ID: uuid.New().String(), Name: "delete_courier"}
)

var (
	PermissionAddIconToList    = models.PermissionModel{ID: uuid.NewString(), Name: "can_add_icon_to_list"}
	PermissionDeleteIconToList = models.PermissionModel{ID: uuid.NewString(), Name: "can_delete_icon_to_list"}
)

// storage permission
var (
	PermissionCreateStore  = models.PermissionModel{ID: uuid.NewString(), Name: "can_add_storage"}
	PermissionGetStoreByID = models.PermissionModel{ID: uuid.NewString(), Name: "can_read_storage"}
	PermissionGetStoreList = models.PermissionModel{ID: uuid.NewString(), Name: "can_read_storage_list"}
	PermissionUpdateStore  = models.PermissionModel{ID: uuid.NewString(), Name: "can_edit_storage"}
	PermissionDeleteStore  = models.PermissionModel{ID: uuid.NewString(), Name: "can_delete_storage"}
)

var PermissionsList = []*models.PermissionModel{
	&PermissionCanMigrateDown,

	&PermissionAddIconToList,
	&PermissionDeleteIconToList,

	&PermissionAddToBasket,
	&PermissionRemoveFromBasket,

	&PermissionCanSeeClients,

	&PermissionIncomeAdd,
	&PermissionIncomeSee,

	&PermissionCanAddRole,
	&PermissionCanEditRole,
	&PermissionCanDeleteRole,
	&PermissionCanViewRole,
	&PermissionCanAttachRole,

	&PermissionCanAddProduct,
	&PermissionCanEditProduct,
	&PermissionCanDeleteProduct,

	&PermissionCategoryAdd,
	&PermissionCategoryEdit,
	&PermissionCategoryDelete,

	&PermissionAddAdmin,
	&PermissionEditAdmin,
	&PermissionDeleteAdmin,

	&PermissionAddCourier,
	&PermissionEditCourier,
	&PermissionDeleteCourier,

	&PermissionBrandEdit,
	&PermissionBrandDelete,
	&PermissionBrandAdd,

	&PermissionAddBranch,
	&PermissionEditBranch,
	&PermissionDeleteBranch,

	&PermissionCanFinishOrder,
	&PermissionCanCancelOrder,
	&PermissionCanMakeOrderPicked,
	&PermissionCanSeeOrderByID,

	&PermissionCreateStore,
	&PermissionGetStoreByID,
	&PermissionGetStoreList,
	&PermissionUpdateStore,
	&PermissionDeleteStore,
}

var RolesList = []*models.RoleModel{
	&RoleSuper,
	&RoleOwner,
	&RoleAdmin,
	&RoleClient,
	&RoleCourier,
}

func GetRolesList() *[]*models.RoleModel {
	return &RolesList
}

func GetPermissionsList() *[]*models.PermissionModel {
	return &PermissionsList
}
