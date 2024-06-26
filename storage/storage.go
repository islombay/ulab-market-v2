package storage

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"context"
	"fmt"
	"time"
)

var (
	ErrAlreadyExists = fmt.Errorf("already_exists")

	// ErrNotAffected is used for cheking 'RowsAffected' to
	// be equal 1 or more
	ErrNotAffected = fmt.Errorf("not_affected")

	ErrNoUpdate         = fmt.Errorf("no_update")
	ErrInvalidEnumInput = fmt.Errorf("invalid_enum_input")
)

type StoreInterface interface {
	Close()
	Role() RoleInterface
	User() UserInterface
	Category() CategoryInterface
	Brand() BrandInterface
	Product() ProductInterface
	Basket() BasketInterface
	Icon() IconInterface
	Branch() BranchInterface
	Order() OrderI
	OrderProduct() OrderProductI
	Storage() StoragesInterface
	Favourite() FavouriteI
	Income() IncomeInterface
}

type OrderProductI interface {
	GetByID(ctx context.Context, id string) (*models.OrderProductModel, error)
	Create(ctx context.Context, m []models.OrderProductModel) error
	GetAll(ctx context.Context) ([]models.OrderProductModel, error)

	GetOrderProducts(ctx context.Context, order_id string) ([]models.OrderProductModel, error)
}

type OrderI interface {
	Create(ctx context.Context, m models.OrderModel) error
	Delete(ctx context.Context, id string) error
	ChangeStatus(ctx context.Context, id, status string) error

	GetAll(ctx context.Context, pagination models.Pagination, statuses []string) ([]models.OrderModel, int, error)
	GetByID(ctx context.Context, id string) (*models.OrderModel, error)
	GetNew(ctx context.Context, pagination models.Pagination, forCourier bool) ([]models.OrderModel, int, error)

	GetCourierActiveList(ctx context.Context, pagination models.Pagination, courier_id string) ([]models.OrderModel, int, error)
	GetCourierActiveListCount(ctx context.Context, courier_id string) (int, error)

	OrdersCount(ctx context.Context, user_id string) (int, error)

	MarkPicked(ctx context.Context, order_id, picker_id string, picked_at time.Time) error
	MarkPickedByCourier(ctx context.Context, order_id, courier_id string, picked_at time.Time) error
	MarkDelivered(ctx context.Context, order_id string) error
	MarkDelivering(ctx context.Context, order_id, courier_id string) error

	GetAllByClient(ctx context.Context, user_id string, pagination models.Pagination) ([]models.OrderModel, int, error)
	GetCourierOrders(ctx context.Context, user_id string, pagination models.Pagination) ([]models.OrderModel, int, error)
}

type BranchInterface interface {
	Create(ctx context.Context, m models.BranchModel) error
	GetByID(ctx context.Context, id string) (*models.BranchModel, error)
	GetByName(ctx context.Context, name string) (*models.BranchModel, error)
	GetAll(ctx context.Context) ([]*models.BranchModel, error)

	Change(ctx context.Context, m models.BranchModel) error
	Delete(ctx context.Context, id string) error
}

type IconInterface interface {
	GetIconByName(ctx context.Context, name string) (*models.IconModel, error)
	GetIconByID(ctx context.Context, id string) (*models.IconModel, error)
	AddIcon(ctx context.Context, m models.IconModel) error

	GetAll(ctx context.Context) ([]models.IconModel, error)
	Delete(ctx context.Context, id string) error
}

type BasketInterface interface {
	Add(ctx context.Context, user_id, product_id string, quantity int, created_at time.Time) error
	Get(ctx context.Context, user_id, product_id string) (*models.BasketModel, error)
	GetAll(ctx context.Context, user_id string) ([]models.BasketModel, error)
	Delete(ctx context.Context, user_id, product_id string) error
	DeleteAll(ctx context.Context, user_id string) error

	ChangeQuantity(ctx context.Context, pid, uid string, quantity int) error
}

type ProductInterface interface {
	GetByArticul(ctx context.Context, articul string) (*models.Product, error)
	CreateProduct(ctx context.Context, m models.Product) error

	GetAll(ctx context.Context, query, catid, bid *string, req models.GetProductAllLimits) ([]*models.Product, error)
	GetAllPagination(ctx context.Context, pagination models_v1.ProductPagination) ([]*models.Product, int, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)

	DeleteProductByID(ctx context.Context, id string) error
	ChangeMainImage(ctx context.Context, id, url string, now time.Time) error

	ChangeProductPrice(ctx context.Context, id string, price float32) error

	CreateProductImageFile(ctx context.Context, id, pid, url string) error
	CreateProductVideoFile(ctx context.Context, id, pid, url string) error
	GetProductVideoFilesByID(ctx context.Context, id string) ([]models.ProductMediaFiles, error)
	GetProductImageFilesByID(ctx context.Context, id string) ([]models.ProductMediaFiles, error)

	IncrementViewCount(ctx context.Context, id string) error

	Change(ctx context.Context, m *models.Product) error
}

type RoleInterface interface {
	GetPermissionByID(ctx context.Context, id string) (*models.PermissionModel, error)
	CreatePermission(ctx context.Context, m models.PermissionModel) error
	GetPermissionByName(ctx context.Context, name string) (*models.PermissionModel, error)

	CreateRole(ctx context.Context, m models.RoleModel) error
	GetRoleByID(ctx context.Context, id string) (*models.RoleModel, error)
	GetRoleByName(ctx context.Context, name string) (*models.RoleModel, error)

	GetRoles(ctx context.Context) ([]*models.RoleModel, error)
	GetPermissions(ctx context.Context) ([]*models.PermissionModel, error)

	Attach(ctx context.Context, rId, pId string) error
	Disattach(ctx context.Context, rId, pId string) error
	IsRolePermissionAttachExists(ctx context.Context, rId, pId string) (bool, error)

	GetRolePermissions(ctx context.Context, role_id string) ([]models.PermissionModel, error)
}

type UserInterface interface {
	CreateStaff(ctx context.Context, m models.Staff) error
	GetStaffByLogin(ctx context.Context, l string) (*models.Staff, error)
	GetStaffByRole(ctx context.Context, roleID string) ([]models.Staff, error)
	GetStaffByID(ctx context.Context, id string) (*models.Staff, error)

	DeleteStaff(ctx context.Context, id string) error
	ChangeStaff(ctx context.Context, m models.Staff) error
	ChangeStaffPassword(ctx context.Context, id, pwd string) error

	CreateClient(ctx context.Context, m models.Client) error
	GetClientByEmail(ctx context.Context, e string) (*models.Client, error)
	GetClientByPhone(ctx context.Context, p string) (*models.Client, error)
	GetClientByLogin(ctx context.Context, l string) (*models.Client, error)
	GetClientByID(ctx context.Context, id string) (*models.Client, error)

	GetClientList(ctx context.Context, pagination models.Pagination) ([]models.Client, int, error)

	UpdateClient(ctx context.Context, model models.ClientUpdate) error
}

type CategoryInterface interface {
	Create(ctx context.Context, m models.Category) error
	GetByID(ctx context.Context, id string) (*models.Category, error)

	ChangeImage(ctx context.Context, cid, imageUrl, iconURL *string) error
	ChangeCategory(ctx context.Context, m models.Category) error

	GetSubcategories(ctx context.Context, id string) ([]*models.Category, error)
	GetAll(ctx context.Context, pagination models.Pagination, onlySub bool) ([]*models.Category, int, error)
	GetByName(ctx context.Context, name string) (*models.Category, error)

	GetBrands(ctx context.Context, id string) ([]models.Brand, error)

	DeleteCategory(ctx context.Context, id string) error
}

type BrandInterface interface {
	Create(ctx context.Context, m models.Brand) error
	GetByID(ctx context.Context, id string) (*models.Brand, error)
	GetByName(ctx context.Context, name string) (*models.Brand, error)
	GetAll(ctx context.Context, pagination models.Pagination) ([]*models.Brand, int, error)

	Change(ctx context.Context, m models.Brand) error
	Delete(ctx context.Context, id string) error
}

type StoragesInterface interface {
	Create(context.Context, models_v1.CreateStorage) (models_v1.Storage, error)
	GetByID(context.Context, string) (models_v1.Storage, error)
	GetList(context.Context, models_v1.StorageRequest) (models_v1.StorageResponse, error)
	Update(context.Context, models_v1.UpdateStorage) (models_v1.Storage, error)
	Delete(context.Context, string) error
}

type FavouriteI interface {
	Create(context.Context, string, string) error
	Get(context.Context, string, string) (*models.FavouriteModel, error)
	GetAll(ctx context.Context, uid string) ([]models.FavouriteModel, error)
}

type IncomeInterface interface {
	Create(context.Context, models_v1.CreateIncome) (models_v1.Income, error)
	GetByID(context.Context, string) (models_v1.Income, error)
	GetList(context.Context, models_v1.IncomeRequest) (models_v1.IncomeResponse, error)

	CreateIncomeProduct(context.Context, models_v1.CreateIncomeProduct) (models_v1.IncomeProduct, error)
	GetProductsByIncomeID(context.Context, string) ([]models_v1.IncomeProduct, error)
	GetIncomeProductsList(context.Context, models_v1.IncomeProductRequest) (models_v1.IncomeProductResponse, error)
}
