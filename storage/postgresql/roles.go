package postgresql

import (
	"app/api/models"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RoleRepo struct {
	db  *pgxpool.Pool
	log logs.LoggerInterface
}

func NewRoleRepo(db *pgxpool.Pool, log logs.LoggerInterface) *RoleRepo {
	return &RoleRepo{
		db:  db,
		log: log,
	}
}

func (s *RoleRepo) CreateRole(ctx context.Context, m models.RoleModel) error {
	q := `insert into roles(id, name, description) values($1, $2, $3)`
	_, err := s.db.Exec(ctx, q, m.ID, m.Name, m.Description)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == DuplicateKeyError {
				return storage.ErrAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (s *RoleRepo) GetRoleByID(ctx context.Context, id string) (*models.RoleModel, error) {
	q := `select * from roles
         where id = $1 and deleted_at is null
		 limit 1;`
	m := models.RoleModel{}

	err := s.db.QueryRow(ctx, q, id).Scan(&m.ID, &m.Name, &m.Description,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *RoleRepo) GetRoleByName(ctx context.Context, name string) (*models.RoleModel, error) {
	q := `select * from roles
         where name = $1 and deleted_at is null
         limit 1;`
	m := models.RoleModel{}

	err := s.db.QueryRow(ctx, q, name).Scan(&m.ID, &m.Name, &m.Description,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *RoleRepo) CreatePermission(ctx context.Context, m models.PermissionModel) error {
	q := `insert into permissions(id, name, description) values($1, $2, $3)`
	_, err := s.db.Exec(ctx, q, m.ID, m.Name, m.Description)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == DuplicateKeyError {
				return storage.ErrAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (s *RoleRepo) GetPermissionByID(ctx context.Context, id string) (*models.PermissionModel, error) {
	q := `select * from permissions where id = $1 limit 1;`
	m := models.PermissionModel{}

	err := s.db.QueryRow(ctx, q, id).Scan(&m.ID, &m.Name, &m.Description,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *RoleRepo) GetPermissionByName(ctx context.Context, name string) (*models.PermissionModel, error) {
	q := `select * from permissions where name = $1 limit 1;`
	var m models.PermissionModel

	if err := s.db.QueryRow(ctx, q, name).Scan(&m.ID, &m.Name, &m.Description,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &m, nil
}

func (s *RoleRepo) Attach(ctx context.Context, rId, pId string) error {
	q := `insert into permission_to_role (role_id, permission_id) values($1,$2);`
	if _, err := s.db.Exec(ctx, q, rId, pId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == DuplicateKeyError {
				return storage.ErrAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (s *RoleRepo) IsRolePermissionAttachExists(ctx context.Context, rId, pId string) (bool, error) {
	q := `select * from permission_to_role
         where role_id = $1 and permission_id = $2 and deleted_at is null`
	var m models.AttachPermission
	if err := s.db.QueryRow(ctx, q, rId, pId).Scan(&m.RoleID, &m.PermissionID,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *RoleRepo) GetRolePermissions(ctx context.Context, role_id string) ([]models.PermissionModel, error) {
	q := `select * from permission_to_role where role_id = $1 and deleted_at is null`
	var res []*models.AttachPermission
	rows, err := s.db.Query(ctx, q, role_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		f := models.AttachPermission{}
		if err := rows.Scan(&f.RoleID, &f.PermissionID,
			&f.CreatedAt,
			&f.UpdatedAt,
			&f.DeletedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, &f)
	}

	q = `select * from permissions where id = $1`

	resp := make([]models.PermissionModel, len(res))
	for i, ap := range res {
		tmp := models.PermissionModel{}
		if err := s.db.QueryRow(ctx, q, ap.PermissionID).Scan(
			&tmp.ID,
			&tmp.Name,
			&tmp.Description,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
			&tmp.DeletedAt,
		); err != nil {
			s.log.Error(
				"permission which occured in attach, not found in permissions",
				logs.String("id", ap.PermissionID),
				logs.Error(err),
			)
			return nil, err
		} else {
			resp[i] = tmp
		}
	}
	return resp, nil
}

func (db *RoleRepo) GetRoles(ctx context.Context) ([]*models.RoleModel, error) {
	q := `select * from roles where deleted_at is null`
	var res []*models.RoleModel

	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	for rows.Next() {
		tmp := models.RoleModel{}
		if err := rows.Scan(
			&tmp.ID,
			&tmp.Name,
			&tmp.Description,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
			&tmp.DeletedAt,
		); err != nil {
			return nil, err
		} else {
			res = append(res, &tmp)
		}
	}
	return res, nil
}

func (db *RoleRepo) GetPermissions(ctx context.Context) ([]*models.PermissionModel, error) {
	q := `select * from permissions where deleted_at is null`
	var res []*models.PermissionModel

	rows, _ := db.db.Query(ctx, q)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	for rows.Next() {
		tmp := models.PermissionModel{}
		if err := rows.Scan(
			&tmp.ID,
			&tmp.Name,
			&tmp.Description,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
			&tmp.DeletedAt,
		); err != nil {
			return nil, err
		} else {
			res = append(res, &tmp)
		}
	}
	return res, nil
}

func (db *RoleRepo) Disattach(ctx context.Context, rId, pId string) error {
	q := `update permission_to_role set deleted_at = now() where role_id = $1 and permission_id = $2`
	if _, err := db.db.Exec(ctx, q, rId, pId); err != nil {
		return err
	}
	return nil
}
