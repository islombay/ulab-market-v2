package storage

import (
	"app/api/models"
	"context"
	"fmt"
	"time"
)

var (
	ErrAlreadyExists = fmt.Errorf("already_exists")

	// ErrNotAffected is used for cheking 'RowsAffected' to
	// be equal 1 or more
	ErrNotAffected = fmt.Errorf("not_affected")

	ErrNoUpdate = fmt.Errorf("no_update")
)

type StoreInterface interface {
	Close()
	Role() RoleInterface
	User() UserInterface
	Category() CategoryInterface
	Brand() BrandInterface
	Product() ProductInterface
	Basket() BasketInterface
}

type BasketInterface interface {
	Add(ctx context.Context, user_id, product_id string, quantity int, created_at time.Time) error
	Get(ctx context.Context, user_id, product_id string) (*models.BasketModel, error)
	GetAll(ctx context.Context, user_id string) ([]models.BasketModel, error)
	Delete(ctx context.Context, user_id, product_id string) error
}

type ProductInterface interface {
	GetByArticul(ctx context.Context, articul string) (*models.Product, error)
	CreateProduct(ctx context.Context, m models.Product) error

	GetAll(ctx context.Context, query, catid, bid *string, req models.GetProductAllLimits) ([]*models.Product, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)

	DeleteProductByID(ctx context.Context, id string) error
	ChangeMainImage(ctx context.Context, id, url string, now time.Time) error

	DeleteProduct(ctx context.Context, id string) error

	CreateProductImageFile(ctx context.Context, id, pid, url string) error
	CreateProductVideoFile(ctx context.Context, id, pid, url string) error
	GetProductVideoFilesByID(ctx context.Context, id string) ([]models.ProductMediaFiles, error)
	GetProductImageFilesByID(ctx context.Context, id string) ([]models.ProductMediaFiles, error)
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
}

type CategoryInterface interface {
	Create(ctx context.Context, m models.Category) error
	GetByID(ctx context.Context, id string) (*models.Category, error)

	AddTranslation(ctx context.Context, m models.CategoryTranslation) error

	ChangeImage(ctx context.Context, cid, imageUrl string) error
	ChangeCategory(ctx context.Context, m models.Category) error

	GetTranslations(ctx context.Context, id string) ([]models.CategoryTranslation, error)
	GetSubcategories(ctx context.Context, id string) ([]*models.Category, error)
	GetAll(ctx context.Context) ([]*models.Category, error)

	DeleteCategory(ctx context.Context, id string) error
	DeleteTranslation(ctx context.Context, cid, lang string) error
}

type BrandInterface interface {
	Create(ctx context.Context, m models.Brand) error
	GetByID(ctx context.Context, id string) (*models.Brand, error)
	GetByName(ctx context.Context, name string) (*models.Brand, error)
	GetAll(ctx context.Context) ([]*models.Brand, error)

	Change(ctx context.Context, m models.Brand) error
	Delete(ctx context.Context, id string) error
}
