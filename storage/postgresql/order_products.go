package postgresql

import (
	"app/api/models"
	"app/pkg/logs"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderProductsRepo struct {
	db  *pgxpool.Pool
	log logs.LoggerInterface
}

func NewOrderProductRepo(db *pgxpool.Pool, log logs.LoggerInterface) *OrderProductsRepo {
	return &OrderProductsRepo{
		db:  db,
		log: log,
	}
}

func (db *OrderProductsRepo) Create(ctx context.Context, m []models.OrderProductModel) error {
	q := `insert into order_products(
                           id, order_id, product_id,
                           quantity, product_price
            ) values ($1, $2, $3, $4, $5)`

	tx, err := db.db.Begin(ctx)
	if err != nil {
		return err
	}
	for _, e := range m {
		if _, err := tx.Exec(ctx, q, e.ID, e.OrderID, e.ProductID, e.Quantity, e.ProductPrice); err != nil {
			tx.Rollback(ctx)
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return err
}

func (db *OrderProductsRepo) GetByID(ctx context.Context, id string) (*models.OrderProductModel, error) {
	q := `select
    	id, order_id, product_id,
    	quantity, product_price, total_price,
    	created_at, updated_at, deleted_at
    from order_products where id = $1`

	var res models.OrderProductModel

	if err := db.db.QueryRow(ctx, q, id).Scan(
		&res.ID, &res.OrderID, &res.ProductID,
		&res.Quantity, &res.ProductPrice, &res.TotalPrice,
		&res.CreatedAt, &res.UpdatedAt, &res.DeletedAt,
	); err != nil {
		return nil, err
	}

	return &res, nil
}

func (db *OrderProductsRepo) GetAll(ctx context.Context) ([]models.OrderProductModel, error) {
	q := `select
    	id, order_id, product_id,
    	quantity, product_price, total_price,
    	created_at, updated_at, deleted_at
    from order_products;`

	var res []models.OrderProductModel
	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var tmp models.OrderProductModel
		if err := rows.Scan(
			&tmp.ID, &tmp.OrderID, &tmp.ProductID,
			&tmp.Quantity, &tmp.ProductPrice, &tmp.TotalPrice,
			&tmp.CreatedAt, &tmp.UpdatedAt, &tmp.DeletedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}

	return res, nil
}

func (db *OrderProductsRepo) GetOrderProducts(ctx context.Context, order_id string) ([]models.OrderProductModel, error) {
	q := `select
			op.id, op.product_id, op.quantity, op.product_price, op.total_price,
			op.created_at, op.updated_at, op.deleted_at, p.articul, p.name_uz, p.name_ru
		from order_products as op
		join products as p on p.id = op.product_id
		where order_id = $1`

	var res []models.OrderProductModel
	rows, _ := db.db.Query(ctx, q, order_id)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var tmp models.OrderProductModel
		if err := rows.Scan(
			&tmp.ID, &tmp.ProductID,
			&tmp.Quantity, &tmp.ProductPrice,
			&tmp.TotalPrice, &tmp.CreatedAt,
			&tmp.UpdatedAt, &tmp.DeletedAt,
			&tmp.Aricul, &tmp.NameUz, &tmp.NameRu,
		); err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}

	return res, nil
}
