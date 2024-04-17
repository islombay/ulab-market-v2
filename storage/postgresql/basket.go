package postgresql

import (
	"app/api/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type BasketRepo struct {
	db *pgxpool.Pool
}

func NewBasketRepo(db *pgxpool.Pool) *BasketRepo {
	return &BasketRepo{db: db}
}

func (db *BasketRepo) Add(ctx context.Context, user_id, product_id string, quantity int, created_at time.Time) error {
	q := `insert into
		basket(user_id, product_id, quantity, created_at, deleted_at)
		values($1, $2, $3, $4, $5)`

	_, err := db.db.Exec(ctx, q, user_id, product_id, quantity, created_at, nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *BasketRepo) Get(ctx context.Context, user_id, product_id string) (*models.BasketModel, error) {
	q := `select * from basket where user_id = $1 and product_id = $2 and deleted_at is null`

	var tmp models.BasketModel
	err := db.db.QueryRow(ctx, q, user_id, product_id).Scan(
		&tmp.UserID,
		&tmp.ProductID,
		&tmp.Quantity,
		&tmp.CreatedAt,
		&tmp.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &tmp, nil
}

func (db *BasketRepo) GetAll(ctx context.Context, user_id string) ([]models.BasketModel, error) {
	q := `select * from basket where user_id = $1 and deleted_at is null`

	var res []models.BasketModel
	rows, _ := db.db.Query(ctx, q, user_id)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var tmp models.BasketModel
		if err := rows.Scan(
			&tmp.UserID,
			&tmp.ProductID,
			&tmp.Quantity,
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

func (db *BasketRepo) Delete(ctx context.Context, user_id, product_id string) error {
	//q := `delete from basket where user_id = $1 and product_id = $2`
	q := `update basket set deleted_at = coalesce(deleted_at, now()) where user_id = $1 and product_id = $2`
	_, err := db.db.Exec(ctx, q, user_id, product_id)
	if err != nil {
		return err
	}

	return nil
}
