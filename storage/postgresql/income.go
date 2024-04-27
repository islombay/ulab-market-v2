package postgresql

import (
	models_v1 "app/api/models/v1"
	"app/pkg/logs"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type incomeRepo struct {
	db  *pgxpool.Pool
	log logs.LoggerInterface
}

func NewIncomeRepo(db *pgxpool.Pool, log logs.LoggerInterface) *incomeRepo {
	return &incomeRepo{
		db:  db,
		log: log,
	}
}

// income

func (i *incomeRepo) Create(ctx context.Context, income models_v1.CreateIncome) (models_v1.Income, error) {
	var createdIncome models_v1.Income
	query := `insert into incomes(id, branch_id, total_price, comment)
                         values($1, $2, $3, $4)
                returning id, branch_id, comment, total_price, created_at, updated_at, deleted_at`

	tx, err := i.db.Begin(ctx)
	if err != nil {
		return models_v1.Income{}, err
	}
	if err := tx.QueryRow(ctx, query,
		uuid.NewString(),
		income.BranchID,
		income.TotalPrice,
		income.Comment,
	).Scan(
		&createdIncome.ID,
		&createdIncome.BranchID,
		&createdIncome.Comment,
		&createdIncome.TotalPrice,
		&createdIncome.CreatedAt,
		&createdIncome.UpdatedAt,
		&createdIncome.DeletedAt,
	); err != nil {
		fmt.Println(i.log)
		i.log.Error("error is while creating income", logs.Error(err))
		return models_v1.Income{}, err
	}

	createdIncome.Products = make([]models_v1.IncomeProduct, len(income.Products))

	// create income products
	for index, e := range income.Products {
		q := `insert into income_products (id, income_id, product_id, quantity, product_price, total_price)
                     values($1, $2, $3, $4, $5, $6)
    		returning id, income_id, product_id, quantity, product_price, total_price, created_at, updated_at, deleted_at`

		e.TotalPrice = e.ProductPrice * float32(e.Quantity)

		var tmp models_v1.IncomeProduct

		if err := tx.QueryRow(ctx, q, uuid.NewString(), createdIncome.ID, e.ProductID, e.Quantity, e.ProductPrice, e.TotalPrice).Scan(
			&tmp.ID, &tmp.IncomeID, &tmp.ProductID, &tmp.Quantity, &tmp.ProductPrice, &tmp.TotalPrice,
			&tmp.CreatedAt, &tmp.UpdatedAt, &tmp.DeletedAt,
		); err != nil {
			i.log.Error("error is while scanning income product", logs.Error(err))
			defer tx.Rollback(ctx)
			return models_v1.Income{}, err
		}
		createdIncome.Products[index] = tmp
	}

	if err := tx.Commit(ctx); err != nil {
		i.log.Error("could not commit transaction", logs.Error(err))
		defer tx.Rollback(ctx)
		return models_v1.Income{}, err
	}

	return createdIncome, nil
}

func (i *incomeRepo) GetByID(ctx context.Context, id string) (models_v1.Income, error) {
	var income models_v1.Income

	query := `select id, branch_id, total_price, comment, created_at, updated_at, deleted_at from incomes where deleted_at is null`

	if err := i.db.QueryRow(ctx, query, id).Scan(
		&income.ID,
		&income.BranchID,
		&income.TotalPrice,
		&income.Comment,
		&income.CreatedAt,
		&income.UpdatedAt,
		&income.DeletedAt,
	); err != nil {
		i.log.Error("error is while getting income by id", logs.Error(err))
		return models_v1.Income{}, err
	}

	return income, nil
}

func (i *incomeRepo) GetList(ctx context.Context, request models_v1.IncomeRequest) (models_v1.IncomeResponse, error) {
	var (
		query, countQuery, pagination, filter string
		count                                 = 0
		offset                                = (request.Page - 1) * request.Limit
		response                              models_v1.IncomeResponse
	)

	pagination = ` LIMIT $1 OFFSET $2`

	if request.Search != "" {
		filter += fmt.Sprintf(` and (branch_id = '%s')`, request.Search, request.Search)
	}

	countQuery = `select count(1) from incomes where deleted_at is null ` + filter
	if err := i.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		i.log.Error("error is while scanning count", logs.Error(err))
		return models_v1.IncomeResponse{}, err
	}

	query = `select id, branch_id, total_price, comment, created_at, updated_at, deleted_at from incomes where deleted_at is null` + filter + pagination
	rows, err := i.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		i.log.Error("error is while selecting all from incomes", logs.Error(err))
		return models_v1.IncomeResponse{}, err
	}

	for rows.Next() {
		income := models_v1.Income{}
		if err = rows.Scan(
			&income.ID,
			&income.BranchID,
			&income.TotalPrice,
			&income.Comment,
			&income.CreatedAt,
			&income.UpdatedAt,
			&income.DeletedAt,
		); err != nil {
			i.log.Error("error is while scanning all from incomes", logs.Error(err))
			return models_v1.IncomeResponse{}, err
		}

		response.Incomes = append(response.Incomes, income)
	}

	response.Count = count

	return response, nil
}

func (i *incomeRepo) Delete(ctx context.Context, id string) error {
	q := `update incomes set deleted_at = coalesce(deleted_at, now()) where id = $1`
	if _, err := i.db.Exec(ctx, q, id); err != nil {
		return err
	}

	return nil
}

// income_product

func (i *incomeRepo) CreateIncomeProduct(ctx context.Context, product models_v1.CreateIncomeProduct) (models_v1.IncomeProduct, error) {
	var incomeProduct models_v1.IncomeProduct

	query := `insert into income_products (id, income_id, product_id, quantity, product_price, total_price)
                     values($1, $2, $3, $4, $5, $6)
    returning id, income_id, product_id, quantity, product_price, total_price, created_at, updated_at, deleted_at`

	if err := i.db.QueryRow(ctx, query,
		uuid.New(),
		product.IncomeID,
		product.ProductID,
		product.Quantity,
		product.ProductPrice,
		product.TotalPrice,
	).Scan(
		&incomeProduct.ID,
		&incomeProduct.ProductID,
		&incomeProduct.Quantity,
		&incomeProduct.ProductPrice,
		&incomeProduct.TotalPrice,
		&incomeProduct.CreatedAt,
		&incomeProduct.UpdatedAt,
		&incomeProduct.DeletedAt,
	); err != nil {
		i.log.Error("error is while scanning income product", logs.Error(err))
		return models_v1.IncomeProduct{}, err
	}

	return incomeProduct, nil
}

func (i *incomeRepo) GetByIncomeProductID(ctx context.Context, id string) (models_v1.IncomeProduct, error) {
	var incomeProduct models_v1.IncomeProduct

	query := `select id, income_id, product_id, quantity, product_price, total_price, created_at, updated_at, deleted_at from income_products where deleted_at is null`

	if err := i.db.QueryRow(ctx, query, id).Scan(
		&incomeProduct.ID,
		&incomeProduct.ProductID,
		&incomeProduct.Quantity,
		&incomeProduct.ProductPrice,
		&incomeProduct.TotalPrice,
		&incomeProduct.CreatedAt,
		&incomeProduct.UpdatedAt,
		&incomeProduct.DeletedAt,
	); err != nil {
		i.log.Error("error is while getting income by id", logs.Error(err))
		return models_v1.IncomeProduct{}, err
	}

	return incomeProduct, nil
}

func (i *incomeRepo) GetIncomeProductsList(ctx context.Context, request models_v1.IncomeProductRequest) (models_v1.IncomeProductResponse, error) {
	var (
		query, countQuery, pagination, filter string
		count                                 = 0
		offset                                = (request.Page - 1) * request.Limit
		response                              models_v1.IncomeProductResponse
	)

	pagination = ` LIMIT $1 OFFSET $2`

	if request.Search != "" {
		filter += fmt.Sprintf(` and (income_id = '%s' or product_id = '%s')`, request.Search, request.Search)
	}

	countQuery = `select count(1) from income_products where deleted_at is null ` + filter
	if err := i.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		i.log.Error("error is while scanning count", logs.Error(err))
		return models_v1.IncomeProductResponse{}, err
	}

	query = `select id, income_id, product_id, quantity, product_price, total_price, created_at, updated_at, deleted_at from incomes where deleted_at is null` + filter + pagination
	rows, err := i.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		i.log.Error("error is while selecting all from incomes", logs.Error(err))
		return models_v1.IncomeProductResponse{}, err
	}

	for rows.Next() {
		incomeProduct := models_v1.IncomeProduct{}
		if err = rows.Scan(
			&incomeProduct.ID,
			&incomeProduct.ProductID,
			&incomeProduct.Quantity,
			&incomeProduct.ProductPrice,
			&incomeProduct.TotalPrice,
			&incomeProduct.CreatedAt,
			&incomeProduct.UpdatedAt,
			&incomeProduct.DeletedAt,
		); err != nil {
			i.log.Error("error is while scanning all from incomes", logs.Error(err))
			return models_v1.IncomeProductResponse{}, err
		}

		response.IncomeProducts = append(response.IncomeProducts, incomeProduct)
	}

	response.Count = count

	return response, nil
}
