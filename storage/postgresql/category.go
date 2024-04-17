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

type CategoryRepo struct {
	db  *pgxpool.Pool
	log logs.LoggerInterface
}

func NewCategoryRepo(db *pgxpool.Pool, log logs.LoggerInterface) *CategoryRepo {
	return &CategoryRepo{db: db, log: log}
}

func (db *CategoryRepo) Create(ctx context.Context, m models.Category) error {
	q := `insert into category (id, name_uz, name_ru, parent_id, created_at) values ($1, $2, $3, $4, $5)`
	ra, err := db.db.Exec(ctx, q, m.ID, m.NameUz, m.NameRu, m.ParentID, m.CreatedAt)
	if err != nil {
		var pgcon *pgconn.PgError
		if errors.As(err, &pgcon) {
			if pgcon.Code == DuplicateKeyError {
				return storage.ErrAlreadyExists
			}
		}
		return err
	}

	if ra.RowsAffected() != 1 {
		return storage.ErrNotAffected
	}
	return nil
}

func (db *CategoryRepo) GetByID(ctx context.Context, id string) (*models.Category, error) {
	q := `select * from category where id = $1`
	var m models.Category
	if err := db.db.QueryRow(ctx, q, id).Scan(
		&m.ID,
		&m.NameUz,
		&m.NameRu,
		&m.Image,
		&m.ParentID,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *CategoryRepo) GetAll(ctx context.Context) ([]*models.Category, error) {
	q := `select * from category where parent_id is null and deleted_at is null`
	m := []*models.Category{}
	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var tmp models.Category
		if err := rows.Scan(
			&tmp.ID,
			&tmp.NameUz,
			&tmp.NameRu,
			&tmp.Image,
			&tmp.ParentID,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
			&tmp.DeletedAt,
		); err != nil {
			db.log.Error("could not scan category", logs.Error(err))
		}
		m = append(m, &tmp)
	}
	return m, nil
}

func (db *CategoryRepo) ChangeImage(ctx context.Context, cid, imageUrl string) error {
	q := `update category set image = $1 where id = $2`
	ra, err := db.db.Exec(ctx, q, imageUrl, cid)
	if err != nil {
		return err
	}

	if ra.RowsAffected() != 1 {
		return storage.ErrNotAffected
	}
	return nil
}

func (db *CategoryRepo) ChangeCategory(ctx context.Context, m models.Category) error {
	q := `update category set name_uz = $1, name_ru = $2, parent_id = $3, updated_at = now() where id = $4`
	ra, err := db.db.Exec(ctx, q, m.NameUz, m.NameRu, m.ParentID, m.ID)
	if err != nil {
		return err
	}

	if ra.RowsAffected() != 1 {
		return storage.ErrNotAffected
	}
	return nil
}

func (db *CategoryRepo) GetSubcategories(ctx context.Context, id string) ([]*models.Category, error) {
	q := `select * from category where parent_id = $1`
	var m []*models.Category
	row, _ := db.db.Query(ctx, q, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	for row.Next() {
		var tmp models.Category
		if err := row.Scan(
			&tmp.ID,
			&tmp.NameUz,
			&tmp.NameRu,
			&tmp.Image,
			&tmp.ParentID,
			&tmp.CreatedAt,
		); err != nil {
			db.log.Error("could not subcategory", logs.Error(err), logs.String("cid", id))
		} else {
			m = append(m, &tmp)
		}
	}
	return m, nil
}

func (db *CategoryRepo) DeleteCategory(ctx context.Context, id string) error {
	q := `update category set deleted_at = now() where id = $1`

	_, err := db.db.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}
