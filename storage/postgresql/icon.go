package postgresql

import (
	"app/api/models"
	"app/pkg/logs"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type IconRepo struct {
	db  *pgxpool.Pool
	log logs.LoggerInterface
}

func NewIconRepo(db *pgxpool.Pool, log logs.LoggerInterface) *IconRepo {
	return &IconRepo{
		db:  db,
		log: log,
	}
}

func (db *IconRepo) AddIcon(ctx context.Context, m models.IconModel) error {
	q := `insert into icons_list(id, name, url) values($1, $2,$3)`
	_, err := db.db.Exec(ctx, q, m.ID, m.Name, m.URL)
	if err != nil {
		return err
	}

	return nil
}

func (db *IconRepo) GetIconByID(ctx context.Context, id string) (*models.IconModel, error) {
	q := `select * from icons_list where id = $1 and deleted_at is null`

	var m models.IconModel
	if err := db.db.QueryRow(ctx, q, id).Scan(
		&m.ID,
		&m.Name, &m.URL,
		&m.CreatedAt, &m.DeletedAt, &m.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &m, nil
}

func (db *IconRepo) GetIconByName(ctx context.Context, name string) (*models.IconModel, error) {
	q := `select * from icons_list where name = $1 and deleted_at is null`

	var m models.IconModel
	if err := db.db.QueryRow(ctx, q, name).Scan(
		&m.ID,
		&m.Name, &m.URL,
		&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	); err != nil {
		return nil, err
	}

	return &m, nil
}

func (db *IconRepo) GetAll(ctx context.Context) ([]models.IconModel, error) {
	q := `select * from icons_list where deleted_at is null`

	var m []models.IconModel
	rows, _ := db.db.Query(ctx, q)

	for rows.Next() {
		var tmp models.IconModel
		if err := rows.Scan(
			&tmp.ID,
			&tmp.Name, &tmp.URL,
			&tmp.CreatedAt, &tmp.DeletedAt, &tmp.UpdatedAt,
		); err != nil {
			return nil, err
		}
		m = append(m, tmp)
	}

	return m, nil
}

func (db *IconRepo) Delete(ctx context.Context, id string) error {
	q := `update icons_list set deleted_at = coalesce(deleted_at, now()) where id =$1`

	tx, err := db.db.Begin(ctx)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, q, id); err != nil {
		return err
	}

	q = `update category set icon_id = null where icon_id = $1`
	if _, err := tx.Exec(ctx, q, id); err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}
