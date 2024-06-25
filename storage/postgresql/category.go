package postgresql

import (
	"app/api/models"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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
	q := `insert into category (
		id, name_uz, name_ru, parent_id, created_at, icon_id
		) values ($1, $2, $3, $4, $5, $6)`
	ra, err := db.db.Exec(ctx, q, m.ID, m.NameUz, m.NameRu, m.ParentID, m.CreatedAt, m.IconID)
	if err != nil {
		var pgcon *pgconn.PgError
		if errors.As(err, &pgcon) {
			if pgcon.Code == DuplicateKeyError {
				db.log.Error("category already exists", logs.Error(err))
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
	q := `select
				id, name_uz, name_ru, image, icon_id, parent_id,
       			created_at, updated_at, deleted_at
		from category where id = $1`
	var m models.Category
	if err := db.db.QueryRow(ctx, q, id).Scan(
		&m.ID,
		&m.NameUz,
		&m.NameRu,
		&m.Image,
		&m.IconID,
		&m.ParentID,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *CategoryRepo) GetAll(ctx context.Context, pagination models.Pagination, onlySub bool) ([]*models.Category, int, error) {
	var whereClause strings.Builder

	whereClause.WriteString("c.deleted_at is null")
	if onlySub {
		whereClause.WriteString(" and c.parent_id is not null")
	} else {
		whereClause.WriteString(" and c.parent_id is null")
	}

	if pagination.Search.Query != "" {
		pagination.Query = `'%` + pagination.Query + `%'`
		whereClause.WriteString(`
			and (c.name_uz ilike ` + pagination.Query + ` or
			c.name_ru ilike ` + pagination.Query + `)
		`)
	}

	q := fmt.Sprintf(`
		SELECT
    		c.id, c.name_uz, c.name_ru,
    		c.image, icon.url, c.parent_id,
    		c.created_at, c.updated_at
    	FROM category AS c
		LEFT JOIN icons_list AS icon ON icon.id = c.icon_id
		WHERE %s
		ORDER BY c.created_at DESC
		LIMIT %d OFFSET %d`, whereClause.String(), pagination.Limit, pagination.Offset)

	m := []*models.Category{}
	rows, err := db.db.Query(ctx, q)
	if err != nil {
		return nil, 0, err
	}

	var count int

	if err := db.db.QueryRow(ctx, fmt.Sprintf("SELECT count(*) FROM category AS c WHERE %s", whereClause.String())).Scan(&count); err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var tmp models.Category
		if err := rows.Scan(
			&tmp.ID,
			&tmp.NameUz,
			&tmp.NameRu,
			&tmp.Image,
			&tmp.IconID,
			&tmp.ParentID,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
		); err != nil {
			db.log.Error("could not scan category", logs.Error(err))
		}
		m = append(m, &tmp)
	}
	return m, count, nil
}

// func (db *CategoryRepo) GetAll(ctx context.Context, pagination models.Pagination, onlySub bool) ([]*models.Category, int, error) {
// 	whereClause := "c.deleted_at is null"
// 	if onlySub {
// 		whereClause += ` and c.parent_id is not null`
// 	} else {
// 		whereClause += ` and c.parent_id is null`
// 	}

// 	if pagination.Search.Query != "" {
// 		pagination.Query = `'%` + pagination.Query + `%'`
// 		whereClause += `
// 			and (c.name_uz ilike ` + pagination.Query + ` or
// 			c.name_ru ilike ` + pagination.Query + `)
// 		`
// 	}

// 	q := fmt.Sprintf(`select
//     		c.id, c.name_uz, c.name_ru,
//     		c.image, icon.url, c.parent_id,
//     		c.created_at, c.updated_at
//     	from category as c
// 		left join icons_list as icon on icon.id = c.icon_id
// 		where %s
// 		order by c.created_at desc
// 		limit %d offset %d`, whereClause, pagination.Limit, pagination.Offset)

// 	fmt.Println(q)
// 	m := []*models.Category{}
// 	rows, _ := db.db.Query(ctx, q)
// 	if rows.Err() != nil {
// 		return nil, 0, rows.Err()
// 	}

// 	var count int

// 	if err := db.db.QueryRow(ctx, fmt.Sprintf("select count(*) from category where %s", whereClause)).Scan(&count); err != nil {
// 		return nil, 0, err
// 	}

// 	for rows.Next() {
// 		var tmp models.Category
// 		if err := rows.Scan(
// 			&tmp.ID,
// 			&tmp.NameUz,
// 			&tmp.NameRu,
// 			&tmp.Image,
// 			&tmp.IconID,
// 			&tmp.ParentID,
// 			&tmp.CreatedAt,
// 			&tmp.UpdatedAt,
// 		); err != nil {
// 			db.log.Error("could not scan category", logs.Error(err))
// 		}
// 		m = append(m, &tmp)
// 	}
// 	return m, count, nil
// }

func (db *CategoryRepo) ChangeImage(ctx context.Context, cid, imageUrl, iconURL *string) error {
	updateFields := make(map[string]interface{})
	if imageUrl != nil {
		updateFields["image"] = imageUrl
	}

	if iconURL != nil {
		updateFields["icon_id"] = iconURL
	}

	updateFields["updated_at"] = time.Now()

	if len(updateFields) == 1 {
		return storage.ErrNoUpdate
	}

	setParts := []string{}
	args := []interface{}{}
	iv := 1
	for k, v := range updateFields {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", k, iv))
		args = append(args, v)
		iv++
	}
	q := fmt.Sprintf("update category set %s where id = $%d",
		strings.Join(setParts, ", "), iv)
	args = append(args, cid)

	ra, err := db.db.Exec(ctx, q, args...)
	if err != nil {
		return err
	}

	if ra.RowsAffected() != 1 {
		return storage.ErrNotAffected
	}
	return nil
}

func (db *CategoryRepo) ChangeCategory(ctx context.Context, m models.Category) error {
	updateFields := make(map[string]interface{})
	if m.NameRu != "" {
		updateFields["name_ru"] = m.NameRu
	}
	if m.NameUz != "" {
		updateFields["name_uz"] = m.NameUz
	}
	if m.IconID != nil {
		updateFields["icon_id"] = m.IconID
	}

	updateFields["updated_at"] = time.Now()

	if len(updateFields) == 1 {
		return storage.ErrNoUpdate
	}

	setParts := []string{}
	args := []interface{}{}
	iv := 1
	for k, v := range updateFields {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", k, iv))
		args = append(args, v)
		iv++
	}

	q := fmt.Sprintf("update category set %s where id = $%d",
		strings.Join(setParts, ", "), iv)

	args = append(args, m.ID)

	ra, err := db.db.Exec(ctx, q, args...)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.Code == DuplicateKeyError {
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

func (db *CategoryRepo) GetSubcategories(ctx context.Context, id string) ([]*models.Category, error) {
	q := `select
			c.id, c.name_uz, c.name_ru,
       		c.image, icon.url, c.parent_id, c.created_at,
       		c.updated_at, c.deleted_at
       	from category as c
		left join icons_list as icon on icon.id = c.icon_id
		where c.parent_id = $1 and c.deleted_at is null
	   	order by c.created_at desc`
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
			&tmp.IconID,
			&tmp.ParentID,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
			&tmp.DeletedAt,
		); err != nil {
			db.log.Error("could not subcategory", logs.Error(err), logs.String("cid", id))
		} else {
			m = append(m, &tmp)
		}
	}
	return m, nil
}

func (db *CategoryRepo) DeleteCategory(ctx context.Context, id string) error {
	q := `update category set deleted_at = coalesce(deleted_at, now()) where id = $1`

	_, err := db.db.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *CategoryRepo) GetByName(ctx context.Context, name string) (*models.Category, error) {
	q := `select 
    		id, name_uz, name_ru,
    		image, icon_id, parent_id,
    		created_at, updated_at, deleted_at
		from category where (name_ru = $1 or name_uz = $1) and deleted_at is null`

	var res models.Category
	if err := db.db.QueryRow(ctx, q, name).Scan(
		&res.ID, &res.NameUz, &res.NameRu,
		&res.Image, &res.IconID, &res.ParentID,
		&res.CreatedAt, &res.UpdatedAt, &res.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &res, nil
}

func (db *CategoryRepo) GetBrands(ctx context.Context, id string) ([]models.Brand, error) {
	q := `with all_category_ids as(
				select id from category
				where (id = $1
					or parent_id = $1)
				and deleted_at is null
			)
			select distinct p.brand_id, b.name from all_category_ids as c
				join products as p on p.category_id = c.id
				join brands as b on p.brand_id = b.id;`

	rows, _ := db.db.Query(ctx, q, id)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	res := []models.Brand{}
	for rows.Next() {
		var tmp models.Brand

		if err := rows.Scan(
			&tmp.ID, &tmp.Name,
		); err != nil {
			db.log.Error("could not get category brands", logs.Error(err),
				logs.String("category_id", id))
		} else {
			res = append(res, tmp)
		}
	}

	return res, nil
}
