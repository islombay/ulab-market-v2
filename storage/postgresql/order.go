package postgresql

import (
	"app/api/models"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"
	"strconv"
	"time"

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
	q := `insert into orders(
				id, payment_type, user_id,
				order_id, client_first_name,
				client_last_name, client_phone_number,
				client_comment, delivery_addr_lat,
				delivery_addr_long, delivery_addr_name
				)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := db.db.Exec(ctx, q,
		m.ID, m.PaymentType, m.UserID,
		m.OrderID, m.ClientFirstName,
		m.ClientLastName, m.ClientPhone,
		m.ClientComment, m.DeliveryAddrLat,
		m.DeliveryAddrLong,
		m.DeliveryAddrName,
	)
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
			created_at, updated_at, deleted_at,
			order_id, client_first_name,
			client_last_name, client_phone_number,
			client_comment, delivery_type,
			delivery_addr_lat, delivery_addr_long,
			picker_user_id, picked_at,
			delivering_user_id, delivered_at,
			delivery_addr_name
		from orders where id = $1`

	var res models.OrderModel
	var tmp_order_id int64
	if err := db.db.QueryRow(ctx, q, id).Scan(
		&res.ID, &res.UserID, &res.Status,
		&res.TotalPrice, &res.PaymentType,
		&res.CreatedAt, &res.UpdatedAt, &res.DeletedAt,
		&tmp_order_id, &res.ClientFirstName,
		&res.ClientLastName, &res.ClientPhone,
		&res.ClientComment, &res.DeliveryType,
		&res.DeliveryAddrLat, &res.DeliveryAddrLong,
		&res.PickerUserID, &res.PickedAt,
		&res.DeliverUserID, &res.DeliveredAt,
		&res.DeliveryAddrName,
	); err != nil {
		return nil, err
	}

	res.OrderID = strconv.FormatInt(tmp_order_id, 10)
	return &res, nil
}

func (db *OrderRepo) GetAll(ctx context.Context) ([]models.OrderModel, error) {
	q := `select
			id, user_id, status,
			total_price,payment_type,
			created_at, updated_at, deleted_at,
			order_id, client_first_name,
			client_last_name, client_phone_number,
			client_comment, delivery_type,
			delivery_addr_lat, delivery_addr_long,
			delivering_user_id, delivered_at,
			delivery_addr_name
		from orders
		where deleted_at is null
		order by created_at desc`

	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	res := []models.OrderModel{}

	for rows.Next() {
		var tmp_order_id int64
		var tmp models.OrderModel
		if err := rows.Scan(
			&tmp.ID, &tmp.UserID, &tmp.Status,
			&tmp.TotalPrice, &tmp.PaymentType,
			&tmp.CreatedAt, &tmp.UpdatedAt, &tmp.DeletedAt,
			&tmp_order_id, &tmp.ClientFirstName,
			&tmp.ClientLastName, &tmp.ClientPhone,
			&tmp.ClientComment, &tmp.DeliveryType,
			&tmp.DeliveryAddrLat, &tmp.DeliveryAddrLong,
			&tmp.DeliverUserID, &tmp.DeliveredAt,
			&tmp.DeliveryAddrName,
		); err != nil {
			return nil, err
		}
		tmp.OrderID = strconv.FormatInt(tmp_order_id, 10)
		res = append(res, tmp)
	}
	return res, nil
}

func (db *OrderRepo) GetArchived(ctx context.Context) ([]models.OrderModel, error) {
	q := `select
			id, user_id, status,
			total_price,payment_type,
			created_at, updated_at, deleted_at,
			order_id, client_first_name,
			client_last_name, client_phone_number,
			client_comment, delivery_type,
			delivery_addr_lat, delivery_addr_long,
			delivering_user_id, delivered_at,
			delivery_addr_name
		from orders where status in ('finished', 'canceled')`

	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	res := []models.OrderModel{}

	for rows.Next() {
		var tmp_order_id int64
		var tmp models.OrderModel
		if err := rows.Scan(
			&tmp.ID, &tmp.UserID, &tmp.Status,
			&tmp.TotalPrice, &tmp.PaymentType,
			&tmp.CreatedAt, &tmp.UpdatedAt, &tmp.DeletedAt,
			&tmp_order_id, &tmp.ClientFirstName,
			&tmp.ClientLastName, &tmp.ClientPhone,
			&tmp.ClientComment, &tmp.DeliveryType,
			&tmp.DeliveryAddrLat, &tmp.DeliveryAddrLong,
			&tmp.DeliverUserID, &tmp.DeliveredAt,
			&tmp.DeliveryAddrName,
		); err != nil {
			return nil, err
		}
		tmp.OrderID = strconv.FormatInt(tmp_order_id, 10)
		res = append(res, tmp)
	}
	return res, nil
}

func (db *OrderRepo) GetActive(ctx context.Context) ([]models.OrderModel, error) {
	q := `select
			id, user_id, status,
			total_price,payment_type,
			created_at, updated_at, deleted_at,
			order_id, client_first_name,
			client_last_name, client_phone_number,
			client_comment, delivery_type,
			delivery_addr_lat, delivery_addr_long,
			delivering_user_id, delivered_at,
			delivery_addr_name
		from orders where status in ('in_process', 'picked', 'delivering')`

	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	res := []models.OrderModel{}

	for rows.Next() {
		var tmp_order_id int64
		var tmp models.OrderModel
		if err := rows.Scan(
			&tmp.ID, &tmp.UserID, &tmp.Status,
			&tmp.TotalPrice, &tmp.PaymentType,
			&tmp.CreatedAt, &tmp.UpdatedAt, &tmp.DeletedAt,
			&tmp_order_id, &tmp.ClientFirstName,
			&tmp.ClientLastName, &tmp.ClientPhone,
			&tmp.ClientComment, &tmp.DeliveryType,
			&tmp.DeliveryAddrLat, &tmp.DeliveryAddrLong,
			&tmp.DeliverUserID, &tmp.DeliveredAt,
			&tmp.DeliveryAddrName,
		); err != nil {
			return nil, err
		}
		tmp.OrderID = strconv.FormatInt(tmp_order_id, 10)
		res = append(res, tmp)
	}
	return res, nil
}

func (db *OrderRepo) GetNew(ctx context.Context, forCourier bool) ([]models.OrderModel, error) {
	q := `select
			id, user_id, status,
			total_price,payment_type,
			created_at, updated_at, deleted_at,
			order_id, client_first_name,
			client_last_name, client_phone_number,
			client_comment, delivery_type,
			delivery_addr_lat, delivery_addr_long,
			delivery_addr_name
		from orders
		where status in ('in_process') and deleted_at is null
		order by created_at desc`

	if forCourier {
		q += ` and delivering_user_id is null`
	}

	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	res := []models.OrderModel{}

	for rows.Next() {
		var tmp_order_id int64
		var tmp models.OrderModel
		if err := rows.Scan(
			&tmp.ID, &tmp.UserID, &tmp.Status,
			&tmp.TotalPrice, &tmp.PaymentType,
			&tmp.CreatedAt, &tmp.UpdatedAt, &tmp.DeletedAt,
			&tmp_order_id, &tmp.ClientFirstName,
			&tmp.ClientLastName, &tmp.ClientPhone,
			&tmp.ClientComment, &tmp.DeliveryType,
			&tmp.DeliveryAddrLat, &tmp.DeliveryAddrLong,
			&tmp.DeliveryAddrName,
		); err != nil {
			return nil, err
		}
		tmp.OrderID = strconv.FormatInt(tmp_order_id, 10)
		res = append(res, tmp)
	}
	return res, nil
}

func (db *OrderRepo) GetUserOrdersCount(ctx context.Context, user_id string) (int, error) {
	q := `select count(*) from orders where user_id = $1 and deleted_at is null`

	var res int
	err := db.db.QueryRow(ctx, q, user_id).Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (db *OrderRepo) MarkPicked(ctx context.Context, order_id, picker_id string, picked_at time.Time) error {
	q := `update orders set
			updated_at = $1, picked_at = $1,
			picker_user_id = $2, status = 'picked'
		where id = $3 and deleted_at is null`

	_, err := db.db.Exec(ctx, q, picked_at, picker_id, order_id)

	return err
}
