package postgresql

import (
	"app/api/models"
	"app/storage"
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type branchRepo struct {
	db *pgxpool.Pool
}

func NewBranchRepo(db *pgxpool.Pool) *branchRepo {
	return &branchRepo{db: db}
}

func (db *branchRepo) Create(ctx context.Context, m models.BranchModel) error {
	q := `insert into branches (id, name, open_time, close_time) values ($1, $2, $3, $4)`
	r, err := db.db.Exec(ctx, q, m.ID, m.Name, m.OpenTime, m.CloseTime)
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

func (db *branchRepo) GetByID(ctx context.Context, id string) (*models.BranchModel, error) {
	q := `select
			id, name, created_at, updated_at, deleted_at, open_time, close_time
		from branches where id = $1`
	var m models.BranchModel
	if err := db.db.QueryRow(ctx, q, id).Scan(
		&m.ID,
		&m.Name,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
		&m.OpenTime,
		&m.CloseTime,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *branchRepo) GetByName(ctx context.Context, name string) (*models.BranchModel, error) {
	q := `select 
			id, name, created_at, updated_at, deleted_at, open_time, close_time
	 	from branches where name = $1 and deleted_at is null`
	var m models.BranchModel
	if err := db.db.QueryRow(ctx, q, name).Scan(
		&m.ID,
		&m.Name,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
		&m.OpenTime,
		&m.CloseTime,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *branchRepo) Change(ctx context.Context, m models.BranchModel) error {
	q := `update branches set
			name = $1,
			updated_at = now(),
			open_time = $2,
			close_time = $3
		where id = $4`
	if _, err := db.db.Exec(ctx, q, m.Name, m.OpenTime, m.CloseTime, m.ID); err != nil {
		return err
	}
	return nil
}

func (db *branchRepo) GetAll(ctx context.Context) ([]*models.BranchModel, error) {
	q := `select
			id, name, created_at, updated_at, deleted_at, open_time, close_time
		from branches where deleted_at is null`
	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	res := []*models.BranchModel{}

	for rows.Next() {
		var m models.BranchModel
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.DeletedAt,
			&m.OpenTime,
			&m.CloseTime,
		); err != nil {
			return nil, err
		}
		res = append(res, &m)
	}
	return res, nil
}

func (db *branchRepo) Delete(ctx context.Context, id string) error {
	q := `update branches set deleted_at = coalesce(deleted_at, now()) where id = $1`
	if _, err := db.db.Exec(ctx, q, id); err != nil {
		return err
	}

	return nil
}
