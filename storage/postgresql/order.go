package postgresql

import (
	"app/api/models"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"
	"fmt"
	"strings"
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
				client_first_name,
				client_last_name, client_phone_number,
				client_comment, delivery_addr_lat,
				delivery_addr_long, delivery_addr_name,
				payment_card_type
				)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := db.db.Exec(ctx, q,
		m.ID, m.PaymentType, m.UserID,
		m.ClientFirstName,
		m.ClientLastName, m.ClientPhone,
		m.ClientComment, m.DeliveryAddrLat,
		m.DeliveryAddrLong,
		m.DeliveryAddrName, m.PaymentCardType,
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

	q = `update clients set name = $1, surname = $2 where id = $3`
	_, err = db.db.Exec(ctx, q, m.ClientFirstName, m.ClientLastName, m.UserID)
	if err != nil {
		db.log.Error("could not change the name of client in order", logs.Error(err), logs.String("uid", m.UserID),
			logs.String("name", *m.ClientFirstName), logs.String("surname", *m.ClientLastName))
	}
	return nil
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
			delivery_addr_name, delivering_user_id,
			payment_card_type
		from orders where id = $1`

	var res models.OrderModel
	if err := db.db.QueryRow(ctx, q, id).Scan(
		&res.ID, &res.UserID, &res.Status,
		&res.TotalPrice, &res.PaymentType,
		&res.CreatedAt, &res.UpdatedAt, &res.DeletedAt,
		&res.OrderID, &res.ClientFirstName,
		&res.ClientLastName, &res.ClientPhone,
		&res.ClientComment, &res.DeliveryType,
		&res.DeliveryAddrLat, &res.DeliveryAddrLong,
		&res.PickerUserID, &res.PickedAt,
		&res.DeliverUserID, &res.DeliveredAt,
		&res.DeliveryAddrName, &res.DeliverUserID,
		&res.PaymentCardType,
	); err != nil {
		return nil, err
	}

	return &res, nil
}

func (db *OrderRepo) GetAll(ctx context.Context, pagination models.Pagination, statuses []string) ([]models.OrderModel, int, error) {
	whereClause := "where deleted_at is null"

	if len(statuses) > 0 {
		inClause := "'" + strings.Join(statuses, "', '") + "'"
		whereClause += fmt.Sprintf(" and status in (%s)", inClause)
	}

	q := fmt.Sprintf(`select
			id, user_id, status,
			total_price,payment_type,
			created_at, updated_at, deleted_at,
			order_id, client_first_name,
			client_last_name, client_phone_number,
			client_comment, delivery_type,
			delivery_addr_lat, delivery_addr_long,
			delivering_user_id, delivered_at,
			delivery_addr_name, delivering_user_id,
			payment_card_type,
			(select count(*) from orders %s) as total_count
		from orders %s
		order by created_at desc`, whereClause, whereClause)

	q += fmt.Sprintf(" limit %d offset %d", pagination.Limit, pagination.Offset)

	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	var count int

	res := []models.OrderModel{}

	for rows.Next() {
		var tmp models.OrderModel
		if err := rows.Scan(
			&tmp.ID, &tmp.UserID, &tmp.Status,
			&tmp.TotalPrice, &tmp.PaymentType,
			&tmp.CreatedAt, &tmp.UpdatedAt, &tmp.DeletedAt,
			&tmp.OrderID, &tmp.ClientFirstName,
			&tmp.ClientLastName, &tmp.ClientPhone,
			&tmp.ClientComment, &tmp.DeliveryType,
			&tmp.DeliveryAddrLat, &tmp.DeliveryAddrLong,
			&tmp.DeliverUserID, &tmp.DeliveredAt,
			&tmp.DeliveryAddrName, &tmp.DeliverUserID,
			&tmp.PaymentCardType,
			&count,
		); err != nil {
			return nil, 0, err
		}
		res = append(res, tmp)
	}

	return res, count, nil
}

func (db *OrderRepo) GetNew(ctx context.Context, pagination models.Pagination, forCourier bool) ([]models.OrderModel, int, error) {
	whereClause := "deleted_at is null"

	if forCourier {
		whereClause += ` and status in ('picked')`
		whereClause += ` and delivering_user_id is null`
	} else {
		whereClause += ` and status in ('in_process')`
	}

	q := fmt.Sprintf(`select
			id, user_id, status,
			total_price,payment_type,
			created_at, updated_at, deleted_at,
			order_id, client_first_name,
			client_last_name, client_phone_number,
			client_comment, delivery_type,
			delivery_addr_lat, delivery_addr_long,
			delivery_addr_name, delivering_user_id,
			payment_card_type,
			(
				select count(*) from orders where %s
			)
		from orders
		where %s order by created_at desc`, whereClause, whereClause)

	if !forCourier {
		q += fmt.Sprintf(" limit %d offset %d", pagination.Limit, pagination.Offset)
	}

	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	res := []models.OrderModel{}

	var count int

	for rows.Next() {
		var tmp models.OrderModel
		if err := rows.Scan(
			&tmp.ID, &tmp.UserID, &tmp.Status,
			&tmp.TotalPrice, &tmp.PaymentType,
			&tmp.CreatedAt, &tmp.UpdatedAt, &tmp.DeletedAt,
			&tmp.OrderID, &tmp.ClientFirstName,
			&tmp.ClientLastName, &tmp.ClientPhone,
			&tmp.ClientComment, &tmp.DeliveryType,
			&tmp.DeliveryAddrLat, &tmp.DeliveryAddrLong,
			&tmp.DeliveryAddrName, &tmp.DeliverUserID,
			&tmp.PaymentCardType, &count,
		); err != nil {
			return nil, 0, err
		}
		res = append(res, tmp)
	}
	return res, count, nil
}

func (db *OrderRepo) GetCourierActiveListCount(ctx context.Context, courier_id string) (int, error) {
	q := fmt.Sprintf(
		`select count(*) from orders
		where deleted_at is null and delivering_user_id = '%s' and
		status in ('in_process', 'picked', 'delivering')`,
		courier_id,
	)

	var count int
	if err := db.db.QueryRow(ctx, q).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (db *OrderRepo) GetCourierActiveList(ctx context.Context, pagination models.Pagination, courier_id string) ([]models.OrderModel, int, error) {
	whereClause := fmt.Sprintf(
		`deleted_at is null and delivering_user_id = '%s' and
		status in ('in_process', 'picked', 'delivering')`,
		courier_id,
	)
	q := fmt.Sprintf(`select
			id, user_id, status,
			total_price,payment_type,
			created_at, updated_at, deleted_at,
			order_id, client_first_name,
			client_last_name, client_phone_number,
			client_comment, delivery_type,
			delivery_addr_lat, delivery_addr_long,
			delivery_addr_name,
			payment_card_type,
			(
				select count(*) from orders where %s
			)
		from orders
		where %s
		order by created_at desc
		limit %d offset %d
		`, whereClause, whereClause, pagination.Limit, pagination.Offset,
	)

	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	res := []models.OrderModel{}

	var count int

	for rows.Next() {
		var tmp models.OrderModel
		if err := rows.Scan(
			&tmp.ID, &tmp.UserID, &tmp.Status,
			&tmp.TotalPrice, &tmp.PaymentType,
			&tmp.CreatedAt, &tmp.UpdatedAt, &tmp.DeletedAt,
			&tmp.OrderID, &tmp.ClientFirstName,
			&tmp.ClientLastName, &tmp.ClientPhone,
			&tmp.ClientComment, &tmp.DeliveryType,
			&tmp.DeliveryAddrLat, &tmp.DeliveryAddrLong,
			&tmp.DeliveryAddrName,
			&tmp.PaymentCardType, &count,
		); err != nil {
			return nil, 0, err
		}
		tmp.IsDeliveringByCourier = true
		res = append(res, tmp)
	}
	return res, count, nil
}

func (db *OrderRepo) OrdersCount(ctx context.Context, user_id string) (int, error) {
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

func (db *OrderRepo) MarkDelivering(ctx context.Context, order_id, courier_id string) error {
	q := `update orders set
			updated_at = now(),
			status = 'delivering'
		where id = $1 and deleted_at is null and delivering_user_id = $2`

	cmdTag, err := db.db.Exec(ctx, q, order_id, courier_id)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return storage.ErrNotAffected
	}

	return nil
}

func (db *OrderRepo) MarkPickedByCourier(ctx context.Context, order_id, courier_id string, picked_at time.Time) error {
	q := `update orders set
			updated_at = $1,
			delivering_user_id = $2
		where id = $3 and deleted_at is null`

	_, err := db.db.Exec(ctx, q, picked_at, courier_id, order_id)

	return err
}

func (db *OrderRepo) GetAllByClient(ctx context.Context, user_id string, pagination models.Pagination) ([]models.OrderModel, int, error) {
	q := fmt.Sprintf(`select
			id, client_first_name, order_id,
			client_last_name, client_phone_number,
			client_comment, status, total_price, payment_type,
			delivery_type, delivery_addr_lat, delivery_addr_long,
			delivery_addr_name, created_at, payment_card_type, delivered_at,
			(	select
					count(*)
				from orders
				where deleted_at is null and user_id = $1
			) as total_count
		from orders
		where deleted_at is null and user_id = $1
		order by created_at desc
		limit %d offset %d`, pagination.Limit, pagination.Offset)

	row, _ := db.db.Query(ctx, q, user_id)
	if row.Err() != nil {
		return nil, 0, row.Err()
	}

	var res []models.OrderModel

	var count int

	for row.Next() {
		var tmp models.OrderModel
		if err := row.Scan(
			&tmp.ID, &tmp.ClientFirstName, &tmp.OrderID,
			&tmp.ClientLastName, &tmp.ClientPhone,
			&tmp.ClientComment, &tmp.Status, &tmp.TotalPrice,
			&tmp.PaymentType, &tmp.DeliveryType, &tmp.DeliveryAddrLat,
			&tmp.DeliveryAddrLong, &tmp.DeliveryAddrName, &tmp.CreatedAt,
			&tmp.PaymentCardType, &tmp.DeliveredAt,
			&count,
		); err != nil {
			return nil, 0, err
		}

		res = append(res, tmp)
	}

	return res, count, nil
}

func (db *OrderRepo) MarkDelivered(ctx context.Context, order_id string) error {
	q := `update orders set
			status = 'finished',
			updated_at = now(),
			delivered_at = now()
		where id = $1`

	_, err := db.db.Exec(ctx, q, order_id)
	if err != nil {
		return err
	}

	return nil
}

func (db *OrderRepo) GetCourierOrders(ctx context.Context, user_id string, pagination models.Pagination) ([]models.OrderModel, int, error) {
	q := fmt.Sprintf(`select
			id, client_first_name, order_id,
			client_last_name, status, total_price,
			payment_type, delivery_type, delivery_addr_name,
			created_at, payment_card_type,
			(	select
					count(*)
				from orders
				where deleted_at is null and delivering_user_id = $1
			) as total_count
		from orders
		where deleted_at is null and delivering_user_id = $1
		order by delivered_at desc
		limit %d offset %d`, pagination.Limit, pagination.Offset)

	row, _ := db.db.Query(ctx, q, user_id)
	if row.Err() != nil {
		return nil, 0, row.Err()
	}

	var res []models.OrderModel

	var count int

	for row.Next() {
		var tmp models.OrderModel
		if err := row.Scan(
			&tmp.ID, &tmp.ClientFirstName, &tmp.OrderID,
			&tmp.ClientLastName, &tmp.Status, &tmp.TotalPrice,
			&tmp.PaymentType, &tmp.DeliveryType,
			&tmp.DeliveryAddrName, &tmp.CreatedAt,
			&tmp.PaymentCardType,
			&count,
		); err != nil {
			return nil, 0, err
		}

		res = append(res, tmp)
	}

	return res, count, nil
}
