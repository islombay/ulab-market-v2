package postgresql

import (
	"app/api/models"
	"app/pkg/logs"
	"app/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductRepo struct {
	db  *pgxpool.Pool
	log logs.LoggerInterface
}

func NewProductRepo(db *pgxpool.Pool, log logs.LoggerInterface) *ProductRepo {
	return &ProductRepo{
		db:  db,
		log: log,
	}
}

func (db *ProductRepo) CreateProduct(ctx context.Context, m models.Product) error {
	q := `insert into products(
                     id, articul,
                     name_uz, name_ru,
                     description_uz, description_ru,
                     income_price, outcome_price,
                     quantity, main_image,
                     category_id, brand_id,
                     status, created_at, updated_at
	) values(
	         $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
	);`
	rows, err := db.db.Exec(ctx, q,
		m.ID, &m.Articul, m.NameUz, m.NameRu,
		m.DescriptionUz, m.DescriptionRu,
		m.IncomePrice, m.OutcomePrice,
		m.Quantity, m.MainImage,
		m.CategoryID, m.BrandID,
		m.Status, m.CreatedAt, m.UpdatedAt,
	)

	if err != nil {
		return err
	}
	if rows.RowsAffected() != 1 {
		return storage.ErrNotAffected
	}

	return nil
}

func (db *ProductRepo) GetByArticul(ctx context.Context, articul string) (*models.Product, error) {
	q := `select * from products where articul = $1 and deleted_at is null`
	var m models.Product
	if err := db.db.QueryRow(ctx, q, articul).Scan(
		&m.ID, &m.Articul,
		&m.NameUz, &m.NameRu,
		&m.DescriptionUz, &m.DescriptionRu,
		&m.IncomePrice, &m.OutcomePrice,
		&m.Quantity,
		&m.CategoryID, &m.BrandID,
		&m.Rating, &m.Status, &m.MainImage,
		&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *ProductRepo) CreateProductImageFile(ctx context.Context, id, pid, url string) error {
	q := `insert into
	product_image_files (id, product_id, media_file)
	values (
    	$1, $2,$3
	)`
	res, err := db.db.Exec(ctx, q, id, pid, url)
	if err != nil {
		return err
	}

	if res.RowsAffected() != 1 {
		return storage.ErrNotAffected
	}
	return nil
}

func (db *ProductRepo) CreateProductVideoFile(ctx context.Context, id, pid, url string) error {
	q := `insert into
	product_video_files (id, product_id, media_file)
	values (
    	$1, $2,$3
	)`
	res, err := db.db.Exec(ctx, q, id, pid, url)
	if err != nil {
		return err
	}

	if res.RowsAffected() != 1 {
		return storage.ErrNotAffected
	}
	return nil
}

func (db *ProductRepo) DeleteProductByID(ctx context.Context, id string) error {
	q := `delete from products where id = $1`
	if _, err := db.db.Exec(ctx, q, id); err != nil {
		return err
	}
	return nil
}
