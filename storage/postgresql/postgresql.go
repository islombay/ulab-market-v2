package postgresql

import (
	"app/config"
	"app/pkg/logs"
	"app/storage"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db            *pgxpool.Pool
	role          storage.RoleInterface
	user          storage.UserInterface
	category      storage.CategoryInterface
	brand         storage.BrandInterface
	product       storage.ProductInterface
	basket        storage.BasketInterface
	icon          storage.IconInterface
	branch        storage.BranchInterface
	order         storage.OrderI
	order_product storage.OrderProductI
	storages      storage.StoragesInterface
	income        storage.IncomeInterface

	favouite storage.FavouriteI

	log logs.LoggerInterface
}

const (
	DuplicateKeyError = "23505"
	InvalidEnumInput  = "22P02"
)

func NewPostgresStore(cfg config.DBConfig, log logs.LoggerInterface) (storage.StoreInterface, error) {
	conf, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%s sslmode=%s",
		cfg.Host,
		os.Getenv("DB_USER"),
		cfg.DBName,
		os.Getenv("DB_PWD"),
		cfg.Port,
		cfg.SSLMode,
	))
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		return nil, err
	}

	return &Store{
		db:       db,
		role:     NewRoleRepo(db, log),
		user:     NewUserRepo(db),
		category: NewCategoryRepo(db, log),
		brand:    NewBrandRepo(db),
		product:  NewProductRepo(db, log),
		income:   NewIncomeRepo(db, log),
	}, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Role() storage.RoleInterface {
	if s.role == nil {
		s.role = NewRoleRepo(s.db, s.log)
	}
	return s.role
}

func (s *Store) User() storage.UserInterface {
	if s.user == nil {
		s.user = NewUserRepo(s.db)
	}
	return s.user
}

func (s *Store) Category() storage.CategoryInterface {
	if s.category == nil {
		s.category = NewCategoryRepo(s.db, s.log)
	}
	return s.category
}

func (s *Store) Brand() storage.BrandInterface {
	if s.brand == nil {
		s.brand = NewBrandRepo(s.db)
	}
	return s.brand
}

func (s *Store) Product() storage.ProductInterface {
	if s.product == nil {
		s.product = NewProductRepo(s.db, s.log)
	}
	return s.product
}

func (s *Store) Basket() storage.BasketInterface {
	if s.basket == nil {
		s.basket = NewBasketRepo(s.db)
	}
	return s.basket
}

func (s *Store) Icon() storage.IconInterface {
	if s.icon == nil {
		s.icon = NewIconRepo(s.db, s.log)
	}
	return s.icon
}

func (s *Store) Branch() storage.BranchInterface {
	if s.branch == nil {
		s.branch = NewBranchRepo(s.db)
	}

	return s.branch
}

func (s *Store) Order() storage.OrderI {
	if s.order == nil {
		s.order = NewOrderRepo(s.db, s.log)
	}
	return s.order
}

func (s *Store) OrderProduct() storage.OrderProductI {
	if s.order_product == nil {
		s.order_product = NewOrderProductRepo(s.db, s.log)
	}
	return s.order_product
}

func (s *Store) Storage() storage.StoragesInterface {
	if s.storages == nil {
		s.storages = NewStorageRepo(s.db, s.log)
	}
	return s.storages
}

func (s *Store) Favourite() storage.FavouriteI {
	if s.favouite == nil {
		s.favouite = NewFavouriteRepo(s.db)
	}

	return s.favouite
}

func (s *Store) Income() storage.IncomeInterface {
	if s.income == nil {
		s.income = NewIncomeRepo(s.db, s.log)
	}
	return s.income

}
