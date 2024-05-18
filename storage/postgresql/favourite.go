package postgresql

import (
	"app/api/models"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type FavouriteRepo struct {
	db *pgxpool.Pool
}

func NewFavouriteRepo(db *pgxpool.Pool) *FavouriteRepo {
	return &FavouriteRepo{db: db}
}

func (db *FavouriteRepo) Create(ctx context.Context, uid, pid string) error {
	q := `insert into
			favourite(user_id, product_id)
		values($1, $2)`

	_, err := db.db.Exec(ctx, q, uid, pid)
	return err
}

func (db *FavouriteRepo) Get(ctx context.Context, uid, pid string) (*models.FavouriteModel, error) {
	q := `select * from favourite where user_id = $1 and product_id = $2`

	var fav models.FavouriteModel
	err := db.db.QueryRow(ctx, q, uid, pid).Scan(&fav.UserID, &fav.ProductID)
	if err != nil {
		return nil, err
	}
	return &fav, nil
}

func (db *FavouriteRepo) GetAll(ctx context.Context, uid string) ([]models.FavouriteModel, error) {
	q := `select user_id, product_id from favourite where user_id = $1`

	row, _ := db.db.Query(ctx, q, uid)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var res []models.FavouriteModel

	for row.Next() {
		var tmp models.FavouriteModel
		if err := row.Scan(
			&tmp.UserID, &tmp.ProductID,
		); err != nil {
			return nil, err
		}

		res = append(res, tmp)
	}

	return res, nil
}
