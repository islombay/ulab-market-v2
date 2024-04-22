package postgresql

import (
	models_v1 "app/api/models/v1"
	"app/pkg/helper"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type storageRepo struct {
	db  *pgxpool.Pool
	log logs.LoggerInterface
}

func NewStorageRepo(db *pgxpool.Pool, log logs.LoggerInterface) *storageRepo {
	return &storageRepo{db: db, log: log}
}

func (s *storageRepo) Create(ctx context.Context, createStorage models_v1.CreateStorage) (string, error) {
	var id string

	query := `insert into storage (id, product_id, branch_id, total_price, quantity) 
	                          values ($1, $2, $3, $4, $5) returning id`

	if err := s.db.QueryRow(ctx, query,
		uuid.New(),
		createStorage.ProductID,
		createStorage.BranchID,
		createStorage.TotalPrice,
		createStorage.Quantity,
	).Scan(&id); err != nil {
		var pgcon *pgconn.PgError
		if errors.As(err, &pgcon) {
			if pgcon.Code == DuplicateKeyError {
				return "", storage.ErrAlreadyExists
			}
		}
		return "", err
	}

	return id, nil
}

func (s *storageRepo) GetByID(ctx context.Context, id string) (models_v1.Storage, error) {
	query := `select id, product_id, branch_id, total_price, quantity, 
	     created_at, updated_at, deleted_at from storage where id = $1 and deleted_at is null`

	var store models_v1.Storage
	if err := s.db.QueryRow(ctx, query, id).Scan(
		&store.ID,
		&store.ProductID,
		&store.BranchID,
		&store.TotalPrice,
		&store.Quantity,
		&store.CreatedAt,
		&store.UpdatedAt,
		&store.DeletedAt,
	); err != nil {
		return models_v1.Storage{}, err
	}
	return store, nil
}

func (s *storageRepo) GetList(ctx context.Context, store models_v1.StorageRequest) (models_v1.StorageResponse, error) {
	var (
		storages                              []models_v1.Storage
		query, countQuery, pagination, filter string
		count                                 = 0
		offset                                = (store.Page - 1) * store.Limit
	)

	pagination = ` LIMIT $1 OFFSET $2`

	if store.Search != "" {
		filter += fmt.Sprintf(` or product_id = '%s' or branch_id = '%s'`, store.Search, store.Search)
	}

	countQuery = `select count(1) from storages where deleted_at is null` + filter
	if err := s.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		s.log.Error("error is while scanning count", logs.Error(err))
		return models_v1.StorageResponse{}, err
	}

	query = `select id, product_id, branch_id, total_price, quantity, created_at, updated_at, deleted_at
	           from storages where deleted_at is null` + filter + pagination

	rows, err := s.db.Query(ctx, query, store.Limit, offset)
	if err != nil {
		s.log.Error("error is while selecting all from storages", logs.Error(err))
		return models_v1.StorageResponse{}, err
	}

	for rows.Next() {
		var str models_v1.Storage
		if err = rows.Scan(
			&str.ID,
			&str.ProductID,
			&str.BranchID,
			&str.TotalPrice,
			&str.Quantity,
			&str.CreatedAt,
			&str.UpdatedAt,
			&str.DeletedAt,
		); err != nil {
			return models_v1.StorageResponse{}, err
		}

		storages = append(storages, str)
	}

	return models_v1.StorageResponse{
		Storage: storages,
		Count:   count,
	}, nil
}

func (s *storageRepo) Update(ctx context.Context, store models_v1.UpdateStorage) (string, error) {
	var (
		id     string
		params = make(map[string]interface{})
		query  = `update storage set `
		filter = ""
	)

	params["id"] = store.ID

	if store.ProductID != "" {
		params["product_id"] = store.ProductID

		filter += " product_id = @product_id,"
	}

	if store.BranchID != "" {
		params["branch_id"] = store.BranchID

		filter += " branch_id = @branch_id,"
	}

	if store.TotalPrice != 0.0 {
		params["total_price"] = store.TotalPrice

		filter += " total_price = @total_price,"
	}

	if store.Quantity != 0 {
		params["quantity"] = store.Quantity

		filter += " quantity = @quantity,"
	}

	query += filter + ` updated_at = now() where deleted_at is null and id = @id returning id`

	fullQuery, args := helper.ReplaceQueryParams(query, params)

	if err := s.db.QueryRow(ctx, fullQuery, args...).Scan(
		&id,
	); err != nil {
		return "", err
	}

	return id, nil
}

func (s *storageRepo) Delete(ctx context.Context, id string) error {
	q := `update storage set deleted_at = coalesce(deleted_at, now()) where id = $1`
	if _, err := s.db.Exec(ctx, q, id); err != nil {
		return err
	}

	return nil
}
