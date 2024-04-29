package postgresql

import (
	"app/api/models"
	"app/pkg/logs"
	"app/storage"
	"context"
	"fmt"
	"strings"
	"time"

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
                     outcome_price,
                     quantity, main_image,
                     category_id, brand_id,
                     status, created_at, updated_at
	) values(
	         $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
	);`
	rows, err := db.db.Exec(ctx, q,
		m.ID, &m.Articul, m.NameUz, m.NameRu,
		m.DescriptionUz, m.DescriptionRu,
		m.OutcomePrice,
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
		&m.OutcomePrice,
		&m.Quantity,
		&m.CategoryID, &m.BrandID,
		&m.Rating, &m.Status, &m.MainImage,
		&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *ProductRepo) GetByID(ctx context.Context, id string) (*models.Product, error) {
	q := `select
			id, articul, name_uz, name_ru,
			description_uz, description_ru,
			outcome_price, (
				select coalesce(sum(s.quantity), 0) from storage as s
				where s.product_id = $1
			) as quantity, category_id, brand_id,
			rating, status, main_image,
			created_at, updated_at, deleted_at
	from products where id = $1 and deleted_at is null`
	var m models.Product
	if err := db.db.QueryRow(ctx, q, id).Scan(
		&m.ID, &m.Articul,
		&m.NameUz, &m.NameRu,
		&m.DescriptionUz, &m.DescriptionRu,
		&m.OutcomePrice,
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
	q := `update products set deleted_at = coalesce(deleted_at, now()) where id = $1`
	if _, err := db.db.Exec(ctx, q, id); err != nil {
		return err
	}
	return nil
}

func (db *ProductRepo) GetProductImageFiles(ctx context.Context, id string) ([]models.ProductMediaFiles, error) {
	q := `select * from product_image_files where id = $1`
	rows, _ := db.db.Query(ctx, q, id)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()
	var all []models.ProductMediaFiles
	for rows.Next() {
		var tmp models.ProductMediaFiles
		if err := rows.Scan(
			&tmp.ID,
			&tmp.ProductID,
			&tmp.MediaFile,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
			&tmp.DeletedAt,
		); err != nil {
			return nil, err
		}
		all = append(all, tmp)
	}
	return all, nil
}

func (db *ProductRepo) GetProductVideoFiles(ctx context.Context, id string) ([]models.ProductMediaFiles, error) {
	q := `select * from product_video_files where id = $1`
	rows, _ := db.db.Query(ctx, q, id)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var all []models.ProductMediaFiles
	for rows.Next() {
		var tmp models.ProductMediaFiles
		if err := rows.Scan(
			&tmp.ID,
			&tmp.ProductID,
			&tmp.MediaFile,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
			&tmp.DeletedAt,
		); err != nil {
			return nil, err
		}
		all = append(all, tmp)
	}
	return all, nil
}

func (db *ProductRepo) GetAll(ctx context.Context, query, catid, bid *string, req models.GetProductAllLimits) ([]*models.Product, error) {
	q := `select
			id, articul, name_uz, name_ru,
			description_uz, description_ru,
			outcome_price, (
				select coalesce(sum(s.quantity), 0) from storage as s
				where s.product_id = $1
			) as quantity, category_id, brand_id,
			rating, status, main_image,
			created_at, updated_at, deleted_at
		from products`
	var args []interface{}
	var whereClause []string

	var (
		offset = " offset 0"
		limit  = " limit 10"
	)

	if catid != nil {
		whereClause = append(whereClause, fmt.Sprintf("category_id = $%d", len(args)+1))
		args = append(args, *catid)
	}
	if bid != nil {
		whereClause = append(whereClause, fmt.Sprintf("brand_id = $%d", len(args)+1))
		args = append(args, *bid)
	}
	if query != nil {
		whereClause = append(whereClause, fmt.Sprintf("(name_uz ilike $%d or name_ru ilike $%d or description_ru ilike $%d or description_uz ilike $%d)", len(args)+1, len(args)+1, len(args)+1, len(args)+1))
		args = append(args, "%"+*query+"%")
	}
	whereClause = append(whereClause, "deleted_at is null")

	if req.Offset > 0 {
		offset = fmt.Sprintf(" offset %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" limit %d", req.Limit)
	}

	if len(whereClause) > 0 {
		q += " where " + strings.Join(whereClause, " and ")
	}

	q += offset + limit

	fmt.Println(q)

	products := []*models.Product{}
	rows, _ := db.db.Query(ctx, q, args...)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Product
		if err := rows.Scan(
			&m.ID, &m.Articul,
			&m.NameUz, &m.NameRu,
			&m.DescriptionUz, &m.DescriptionRu,
			&m.OutcomePrice,
			&m.Quantity,
			&m.CategoryID, &m.BrandID,
			&m.Rating, &m.Status, &m.MainImage,
			&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
		); err != nil {
			return nil, err
		}
		imgFiles, err := db.GetProductImageFiles(ctx, m.ID)
		if err != nil {
			db.log.Error("could not load image files for product", logs.Error(err),
				logs.String("product_id", m.ID))
		} else {
			m.ImageFiles = imgFiles
		}

		vdFiles, err := db.GetProductVideoFiles(ctx, m.ID)
		if err != nil {
			db.log.Error("could not load video files for product", logs.Error(err),
				logs.String("product_id", m.ID))
		} else {
			m.VideoFiles = vdFiles
		}

		products = append(products, &m)
	}
	return products, nil
}

func (db *ProductRepo) GetProductImageFilesByID(ctx context.Context, id string) ([]models.ProductMediaFiles, error) {
	q := `select * from product_image_files where product_id = $1`
	res := []models.ProductMediaFiles{}
	rows, _ := db.db.Query(ctx, q, id)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var tmp models.ProductMediaFiles
		if err := rows.Scan(
			&tmp.ID, &tmp.ProductID, &tmp.MediaFile,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
			&tmp.DeletedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}

	return res, nil
}

func (db *ProductRepo) GetProductVideoFilesByID(ctx context.Context, id string) ([]models.ProductMediaFiles, error) {
	q := `select * from product_video_files where product_id = $1`
	res := []models.ProductMediaFiles{}
	rows, _ := db.db.Query(ctx, q, id)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var tmp models.ProductMediaFiles
		if err := rows.Scan(
			&tmp.ID, &tmp.ProductID, &tmp.MediaFile,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
			&tmp.DeletedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}

	return res, nil
}

func (db *ProductRepo) ChangeMainImage(ctx context.Context, id, url string, now time.Time) error {
	q := `update products set main_image = $1, updated_at = $3 where id = $2`
	_, err := db.db.Exec(ctx, q, url, id, now)
	if err != nil {
		return err
	}
	return nil
}

func (db *ProductRepo) ChangeProductPrice(ctx context.Context, id string, price float32) error {
	q := `update products set outcome_price = $1, updated_at = now() where id = $2`
	_, err := db.db.Exec(ctx, q, price, id)
	return err
}
