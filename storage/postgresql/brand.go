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
	q := `update brands set name = $1, updated_at = now()`

	if m.Image != nil {
		q += fmt.Sprintf(", image = '%s'", *m.Image)
	}

	q += ` where id = $2`

	if _, err := db.db.Exec(ctx, q, m.Name, m.ID); err != nil {
		return err
	}
	return nil
}

func (db *BrandRepo) GetAll(ctx context.Context) ([]*models.Brand, error) {
	q := `select 
			id, name, image, created_at,
			updated_at, deleted_at
		from brands where deleted_at is null
		order by created_at desc`
	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	res := []*models.Brand{}

	for rows.Next() {
		var m models.Brand
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Image,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.DeletedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, &m)
	}
	return res, nil
}

func (db *BrandRepo) Delete(ctx context.Context, id string) error {
	q := `update brands set deleted_at = coalesce(deleted_at, now()) where id = $1`
	if _, err := db.db.Exec(ctx, q, id); err != nil {
		return err
	}

	return nil
}
