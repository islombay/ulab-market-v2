package models

type FavouriteModel struct {
	UserID    string `db:"user_id" json:"user_id,omitempty"`
	ProductID string `db:"product_id" json:"product_id"`
}
