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
