package postgresql

import (
	models_v1 "app/api/models/v1"
	"app/pkg/helper"
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

func (i *incomeRepo) Create(ctx context.Context, income models_v1.CreateIncome) (models_v1.Income, error) {
	var createdIncome models_v1.Income
	query := `insert into incomes(id, branch_id, total_price, comment, courier_id)
                         values($1, $2, $3, $4, $5)
                returning id, branch_id, total_price, comment, courier_id`
	if err := i.db.QueryRow(ctx, query,
		uuid.New(),
		income.BranchID,
		income.TotalPrice,
		income.Comment,
		income.CourierID,
	).Scan(
		&createdIncome.ID,
		&createdIncome.BranchID,
		&createdIncome.TotalPrice,
		&createdIncome.Comment,
		&createdIncome.CourierID,
	); err != nil {
		i.log.Error("error is while creating income", logs.Error(err))
		return models_v1.Income{}, err
	}

	return createdIncome, nil
}

func (i *incomeRepo) GetByID(ctx context.Context, id string) (models_v1.Income, error) {
	var income models_v1.Income

	query := `select id, branch_id, total_price, comment, courier_id from incomes where deleted_at is null`

	if err := i.db.QueryRow(ctx, query, id).Scan(
		&income.ID,
		&income.BranchID,
		&income.TotalPrice,
		&income.Comment,
		&income.CourierID,
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
		filter += fmt.Sprintf(` and (branch_id = '%s' or courier_id = '%s')`, request.Search, request.Search)
	}

	countQuery = `select count(1) from incomes where deleted_at is null ` + filter
	if err := i.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		i.log.Error("error is while scanning count", logs.Error(err))
		return models_v1.IncomeResponse{}, err
	}

	query = `select id, branch_id, total_price, comment, courier_id from incomes where deleted_at is null` + filter + pagination
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
			&income.CourierID,
		); err != nil {
			i.log.Error("error is while scanning all from incomes", logs.Error(err))
			return models_v1.IncomeResponse{}, err
		}

		response.Incomes = append(response.Incomes, income)
	}

	response.Count = count

	return response, nil
}

func (i *incomeRepo) Update(ctx context.Context, income models_v1.UpdateIncome) (models_v1.Income, error) {
	var (
		updatedIncome models_v1.Income
		params        = make(map[string]interface{})
		query         = `update incomes set `
		filter        = ""
	)

	params["id"] = income.ID

	if income.BranchID != "" {
		params["branch_id"] = income.BranchID

		filter += " branch_id = @branch_id,"
	}

	if income.TotalPrice != 0.0 {
		params["total_price"] = income.TotalPrice

		filter += " total_price = @total_price,"
	}

	if income.Comment != "" {
		params["comment"] = income.Comment

		filter += " comment = @comment,"
	}

	if income.CourierID != "" {
		params["courier_id"] = income.CourierID

		filter += " courier_id = @courier_id,"
	}

	query += filter + ` updated_at = now() where deleted_at is null and id = @id returning id, branch_id, total_price, comment, courier_id, created_at, updated_at, deleted_at`

	fullQuery, args := helper.ReplaceQueryParams(query, params)

	if err := i.db.QueryRow(ctx, fullQuery, args...).Scan(
		&updatedIncome.ID,
		&updatedIncome.BranchID,
		&updatedIncome.TotalPrice,
		&updatedIncome.Comment,
		&updatedIncome.CourierID,
		&updatedIncome.CreatedAt,
		&updatedIncome.UpdatedAt,
		&updatedIncome.DeletedAt,
	); err != nil {
		return models_v1.Income{}, err
	}

	return updatedIncome, nil
}

func (i *incomeRepo) Delete(ctx context.Context, id string) error {
	q := `update incomes set deleted_at = coalesce(deleted_at, now()) where id = $1`
	if _, err := i.db.Exec(ctx, q, id); err != nil {
		return err
	}

	return nil
}
