package postgresql

import (
	"app/api/models"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"
	"fmt"
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
	q := `insert into category (id, name, parent_id) values ($1, $2, $3)`
	fmt.Println(m.ParentID == nil)
	ra, err := db.db.Exec(ctx, q, m.ID, m.Name, m.ParentID)
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
		&m.Name,
		&m.Image,
		&m.ParentID,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *CategoryRepo) GetAll(ctx context.Context) ([]*models.Category, error) {
	q := `select * from category where parent_id is null`
	m := []*models.Category{}
	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var tmp models.Category
		if err := rows.Scan(
			&tmp.ID,
			&tmp.Name,
			&tmp.Image,
			&tmp.ParentID,
		); err != nil {
			db.log.Error("could not scan category", logs.Error(err))
		}
		m = append(m, &tmp)
	}
	return m, nil
}

func (db *CategoryRepo) AddTranslation(ctx context.Context, m models.CategoryTranslation) error {
	q := `insert into category_translation (category_id, name, language) values ($1, $2, $3)`
	ra, err := db.db.Exec(ctx, q, m.CategoryID, m.Name, m.LanguageCode)
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

func (db *CategoryRepo) DeleteTranslation(ctx context.Context, cid, lang string) error {
	q := `delete from category_translation where category_id = $1 and language = $2`
	if _, err := db.db.Exec(ctx, q, cid, lang); err != nil {
		return err
	}
	return nil
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
	q := `update category set name = $1, parent_id = $2 where id = $3`
	ra, err := db.db.Exec(ctx, q, m.Name, m.ParentID, m.ID)
	if err != nil {
		return err
	}

	if ra.RowsAffected() != 1 {
		return storage.ErrNotAffected
	}
	return nil
}

func (db *CategoryRepo) GetTranslations(ctx context.Context, id string) ([]models.CategoryTranslation, error) {
	q := `select * from category_translation where category_id = $1`
	r, _ := db.db.Query(ctx, q, id)
	if err := r.Err(); err != nil {
		return nil, err
	}
	defer r.Close()
	res := []models.CategoryTranslation{}
	for r.Next() {
		var tm models.CategoryTranslation
		if err := r.Scan(
			&tm.CategoryID, &tm.Name, &tm.LanguageCode,
		); err != nil {
			db.log.Error("could not scan translation", logs.Error(err))
		}
		res = append(res, tm)
	}

	return res, nil
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
			&tmp.Name,
			&tmp.Image,
			&tmp.ParentID,
		); err != nil {
			db.log.Error("could not subcategory", logs.Error(err), logs.String("cid", id))
		} else {
			m = append(m, &tmp)
		}
	}
	return m, nil
}

func (db *CategoryRepo) DeleteCategory(ctx context.Context, id string) error {
	q := `delete from category_translation where category_id = $1`
	tx, err := db.db.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	q = `delete from category where id = $1`
	_, err = tx.Exec(ctx, q, id)
	if err != nil {
		defer tx.Rollback(ctx)
		return err
	}
	defer tx.Commit(ctx)
	return nil
}
