package postgresql

import (
	"app/api/models"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderRepo struct {
	db  *pgxpool.Pool
	log logs.LoggerInterface
}

func NewOrderRepo(db *pgxpool.Pool, log logs.LoggerInterface) *OrderRepo {
	return &OrderRepo{
		db:  db,
		log: log,
	}
}

func (db *OrderRepo) Create(ctx context.Context, m models.OrderModel) error {
	q := `insert into orders(id, payment_type, user_id)
			values ($1, $2, $3)`

	_, err := db.db.Exec(ctx, q, m.ID, m.PaymentType, m.UserID)
	if err != nil {
		var pgcon *pgconn.PgError
		if errors.As(err, &pgcon) {
			if pgcon.Code == InvalidEnumInput {
				return storage.ErrInvalidEnumInput
			}
		}
		return err
	}
	return err
}

func (db *OrderRepo) Delete(ctx context.Context, id string) error {
	q := `update orders set deleted_at = now() where id = $1`
	_, err := db.db.Exec(ctx, q, id)
	return err
}

func (db *OrderRepo) ChangeStatus(ctx context.Context, id, status string) error {
	q := `update orders set status = $1, updated_at=now() where id = $2`
	if _, err := db.db.Exec(ctx, q, status, id); err != nil {
		var pgcon *pgconn.PgError
		if errors.As(err, &pgcon) {
			if pgcon.Code == InvalidEnumInput {
				return storage.ErrInvalidEnumInput
			}
		}
		return err
	}
	return nil
}

func (db *OrderRepo) GetByID(ctx context.Context, id string) (*models.OrderModel, error) {
	q := `select
			id, user_id, status,
			total_price,payment_type,
			created_at, updated_at, deleted_at
		from orders where id = $1`

	var res models.OrderModel
	if err := db.db.QueryRow(ctx, q, id).Scan(
		&res.ID, &res.UserID, &res.Status,
		&res.TotalPrice, &res.PaymentType,
		&res.CreatedAt, &res.UpdatedAt, &res.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &res, nil
}
