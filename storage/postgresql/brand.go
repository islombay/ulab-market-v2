package postgresql

import (
	"app/api/models"
	"app/storage"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BrandRepo struct {
	db *pgxpool.Pool
}

func NewBrandRepo(db *pgxpool.Pool) storage.BrandInterface {
	return &BrandRepo{db: db}
}

func (db *BrandRepo) Create(ctx context.Context, m models.Brand) error {
	q := `insert into brands (id, name, image) values ($1, $2, $3)`
	r, err := db.db.Exec(ctx, q, m.ID, m.Name, m.Image)
	if err != nil {
		var pgcon *pgconn.PgError
		if errors.As(err, &pgcon) {
			if pgcon.Code == DuplicateKeyError {
				return storage.ErrAlreadyExists
			}
		}
		return err
	}

	if r.RowsAffected() != 1 {
		return storage.ErrNotAffected
	}
	return nil
}

func (db *BrandRepo) GetByID(ctx context.Context, id string) (*models.Brand, error) {
	q := `select
			id, name, image, created_at,
			updated_at, deleted_at
		from brands where id = $1`
	var m models.Brand
	if err := db.db.QueryRow(ctx, q, id).Scan(
		&m.ID,
		&m.Name,
		&m.Image,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *BrandRepo) GetByName(ctx context.Context, name string) (*models.Brand, error) {
	q := `select 
			id, name, image, created_at,
			updated_at, deleted_at 
		from brands where name = $1 and deleted_at is null limit 1`
	var m models.Brand
	if err := db.db.QueryRow(ctx, q, name).Scan(
		&m.ID,
		&m.Name,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *BrandRepo) Change(ctx context.Context, m models.Brand) error {
	q := `update brands set updated_at = now()`

	if m.Image != nil {
		q += fmt.Sprintf(", image = '%s'", *m.Image)
	}

	if m.Name != "" {
		q += fmt.Sprintf(", name = '%s'", m.Name)
	}

	q += ` where id = $1`

	if _, err := db.db.Exec(ctx, q, m.ID); err != nil {
		var pgcon *pgconn.PgError
		if errors.As(err, &pgcon) {
			if pgcon.Code == DuplicateKeyError {
				return storage.ErrAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (db *BrandRepo) GetAll(ctx context.Context, pagination models.Pagination) ([]*models.Brand, int, error) {
	whereClause := `deleted_at is null`

	if pagination.Query != "" {
		pagination.Query = "'%" + pagination.Query + "%'"
		whereClause += " and name ilike " + pagination.Query
	}

	q := fmt.Sprintf(`select 
			id, name, image, created_at,
			updated_at, deleted_at,
			(
				select count(*) from brands
				where %s
			)
		from brands where %s
		order by created_at desc
		limit %d offset %d`, whereClause, whereClause, pagination.Limit, pagination.Offset)
	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	res := []*models.Brand{}

	var count int

	for rows.Next() {
		var m models.Brand
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Image,
			&m.CreatedAt,
			&m.UpdatedAt,
			m.DeletedAt,
			&count,
		); err != nil {
			return nil, 0, err
		}
		res = append(res, &m)
	}
	return res, count, nil
}

func (db *BrandRepo) Delete(ctx context.Context, id string) error {
	q := `update brands set deleted_at = coalesce(deleted_at, now()) where id = $1`
	if _, err := db.db.Exec(ctx, q, id); err != nil {
		return err
	}

	return nil
}
