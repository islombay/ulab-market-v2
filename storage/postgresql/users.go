package postgresql

import (
	"app/api/models"
	"app/storage"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (db *UserRepo) CreateClient(ctx context.Context, m models.Client) error {
	q := `insert into clients (
                     id,
                     name,
                     phone_number,
                     email,
                     created_at,
                     updated_at, 
                     deleted_at
            ) values ($1,$2,$3,$4,$5,$6,$7)`
	if _, err := db.db.Exec(ctx, q,
		m.ID,
		m.Name,
		m.PhoneNumber,
		m.Email,
		m.CreatedAt,
		m.UpdatedAt,
		m.DeletedAt,
	); err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.Code == DuplicateKeyError {
				return storage.ErrAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (db *UserRepo) CreateStaff(ctx context.Context, m models.Staff) error {
	usr, _ := db.GetStaffByLogin(ctx, models.GetStringValue(m.Email))
	usr2, _ := db.GetStaffByLogin(ctx, models.GetStringValue(m.PhoneNumber))
	if usr != nil || usr2 != nil {
		return storage.ErrAlreadyExists
	}
	q := `insert into staff (
                     id,
                     name,
                     phone_number,
                     email,
                     password,
                     role_id,
                     created_at,
                     updated_at, 
                     deleted_at
            ) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	if _, err := db.db.Exec(ctx, q,
		m.ID,
		m.Name,
		m.PhoneNumber,
		m.Email,
		m.Password,
		m.RoleID,
		m.CreatedAt,
		m.UpdatedAt,
		m.DeletedAt,
	); err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.Code == DuplicateKeyError {
				return storage.ErrAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (db *UserRepo) GetClientByEmail(ctx context.Context, e string) (*models.Client, error) {
	q := `select 
			id, name, surname, phone_number,
			email, created_at, updated_at,
			deleted_at
		from clients where email = $1 limit 1;`
	var m models.Client

	if err := db.db.QueryRow(ctx, q, e).Scan(
		&m.ID, &m.Name, &m.Surname,
		&m.PhoneNumber, &m.Email,
		&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *UserRepo) GetClientByPhone(ctx context.Context, p string) (*models.Client, error) {
	q := `select 
				id, name, surname, phone_number,
				email, created_at, updated_at,
				deleted_at
			from clients where phone_number = $1 limit 1;`
	var m models.Client

	if err := db.db.QueryRow(ctx, q, p).Scan(
		&m.ID, &m.Name, &m.Surname,
		&m.PhoneNumber, &m.Email,
		&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *UserRepo) GetClientByLogin(ctx context.Context, l string) (*models.Client, error) {
	q := `select
				id, name, surname, phone_number,
				email, created_at, updated_at,
				deleted_at
			from clients where phone_number = $1 or email = $1 limit 1;`
	var m models.Client

	if err := db.db.QueryRow(ctx, q, l).Scan(
		&m.ID, &m.Name, &m.Surname,
		&m.PhoneNumber, &m.Email,
		&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *UserRepo) GetStaffByID(ctx context.Context, id string) (*models.Staff, error) {
	q := `select * from staff where id = $1 and deleted_at is null;`
	var m models.Staff

	if err := db.db.QueryRow(ctx, q, id).Scan(
		&m.ID,
		&m.Name, &m.PhoneNumber, &m.Email, &m.Password,
		&m.RoleID, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *UserRepo) GetClientByID(ctx context.Context, id string) (*models.Client, error) {
	q := `select 
			id, name, surname, phone_number, email
			created_at, updated_at, deleted_at
		from clients where id = $1 and deleted_at is null;`
	var m models.Client

	if err := db.db.QueryRow(ctx, q, id).Scan(
		&m.ID,
		&m.Name, &m.Surname, &m.PhoneNumber, &m.Email,
		&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *UserRepo) ChangeStaffPassword(ctx context.Context, id, pwd string) error {
	q := `update staff set password = $1 where id = $2`
	r, err := db.db.Exec(ctx, q, pwd, id)
	if err != nil {
		return err
	}
	if r.RowsAffected() != 1 {
		return storage.ErrNotAffected
	}
	return nil
}

func (db *UserRepo) GetStaffByLogin(ctx context.Context, l string) (*models.Staff, error) {
	q := `select 
			id, name, phone_number, email,
			password, role_id, created_at,
			updated_at, deleted_at
		from staff where (phone_number = $1 or email = $1) and deleted_at is null limit 1;`
	var m models.Staff

	if err := db.db.QueryRow(ctx, q, l).Scan(
		&m.ID,
		&m.Name, &m.PhoneNumber, &m.Email, &m.Password,
		&m.RoleID, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (db *UserRepo) GetStaffByRole(ctx context.Context, roleID string) ([]models.Staff, error) {
	q := `select
				id, name, phone_number, email,
				password, role_id, created_at,
				updated_at, deleted_at
			from staff where role_id = $1 and deleted_at is null`
	m := []models.Staff{}

	rows, _ := db.db.Query(ctx, q, roleID)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	for rows.Next() {
		mt := models.Staff{
			ID:          "",
			Name:        "",
			PhoneNumber: nil,
			Email:       nil,
			Password:    "",
			RoleID:      "",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			DeletedAt:   nil,
		}
		err := rows.Scan(&mt.ID,
			&mt.Name, &mt.PhoneNumber, &mt.Email, &mt.Password,
			&mt.RoleID, &mt.CreatedAt, &mt.UpdatedAt, &mt.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		m = append(m, mt)
	}
	return m, nil
}

func (db *UserRepo) DeleteStaff(ctx context.Context, id string) error {
	now := time.Now()
	q := `update staff set deleted_at=$1 where id = $2`
	_, err := db.GetStaffByID(ctx, id)
	if err != nil {
		return err
	}
	upt, err := db.db.Exec(ctx, q,
		now,
		id,
	)
	if err != nil {
		return err
	}
	if upt.RowsAffected() != 1 {
		return storage.ErrNotAffected
	}

	return nil
}

func (db *UserRepo) ChangeStaff(ctx context.Context, m models.Staff) error {
	updateFields := make(map[string]interface{})
	if m.Name != "" {
		updateFields["name"] = m.Name
	}
	if m.Email != nil {
		updateFields["email"] = m.Email
	}
	if m.PhoneNumber != nil {
		updateFields["phone_number"] = m.PhoneNumber
	}
	if m.RoleID != "" {
		updateFields["role_id"] = m.RoleID
	}
	if m.Password != "" {
		updateFields["password"] = m.Password
	}
	updateFields["updated_at"] = time.Now()

	if len(updateFields) == 0 {
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
	q := fmt.Sprintf("update staff set %s where id = $%d",
		strings.Join(setParts, ", "), iv)
	args = append(args, m.ID)

	if _, err := db.db.Exec(ctx, q, args...); err != nil {
		return err
	}
	return nil
}

func (db *UserRepo) GetClientList(ctx context.Context) ([]models.Client, error) {
	q := `select
			id, name, surname, phone_number,
			email, created_at, updated_at
		from clients
		where deleted_at is null
		order by created_at desc`

	var res []models.Client

	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var tmp models.Client

		if err := rows.Scan(
			&tmp.ID, &tmp.Name, &tmp.Surname,
			&tmp.PhoneNumber, &tmp.Email,
			&tmp.CreatedAt, &tmp.UpdatedAt,
		); err != nil {
			return nil, err
		}

		res = append(res, tmp)
	}

	return res, nil
}
