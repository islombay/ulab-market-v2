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
