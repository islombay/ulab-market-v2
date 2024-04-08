package start

import (
	"app/api/models"
	"app/config"
	auth_lib "app/pkg/auth"
	"app/pkg/logs"
	"app/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

func Init(cfg *config.DBConfig, log logs.LoggerInterface, isDown bool, rolesDB storage.RoleInterface, userDB storage.UserInterface) error {
	if err := migration(cfg, log, isDown); err != nil {
		return err
	}
	if err := roles_setup(rolesDB, log); err != nil {
		return err
	}
	if err := defaultUser(userDB, log); err != nil {
		return err
	}
	return nil
}

func migration(cfg *config.DBConfig, log logs.LoggerInterface, isDown bool) error {
	var dbURL string
	if os.Getenv("ENV") == config.LocalMode {
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			os.Getenv("DB_USER"), os.Getenv("DB_PWD"),
			cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	} else if os.Getenv("ENV") == config.ProdMode {
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			os.Getenv("DB_USER"), os.Getenv("DB_PWD"),
			cfg.Host, cfg.Port, cfg.DBName)
	}

	migrationsPath := fmt.Sprintf("file://%s", cfg.MigrationsPath)

	log.Debug("initializing migrations")
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Error("could not create new migration", logs.Error(err))
		return err
	}

	if isDown {
		log.Debug("migrating down")
		err = m.Down()
		if err != nil {
			log.Error("could not migrate down", logs.Error(err))
			return err
		}
		log.Info("migrated down")
	}

	log.Debug("migrating up")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error("could not migrate up", logs.Error(err))
		return err
	}
	log.Info("migrated up")

	return nil
}

func roles_setup(db storage.RoleInterface, log logs.LoggerInterface) error {
	for _, e := range *auth_lib.GetPermissionsList() {
		err := db.CreatePermission(context.Background(), *e)
		if err != nil {
			if !errors.Is(err, storage.ErrAlreadyExists) {
				log.Error("could not create permission in start", logs.Error(err))
				return err
			}
		}
		r, _ := db.GetPermissionByName(context.Background(), e.Name)
		e.ID = r.ID
	}
	for _, e := range *auth_lib.GetRolesList() {
		err := db.CreateRole(context.Background(), *e)
		if err != nil {
			if !errors.Is(err, storage.ErrAlreadyExists) {
				log.Error("could not create role in start", logs.Error(err))
				return err
			}
		}
		r, err := db.GetRoleByName(context.Background(), e.Name)
		e.ID = r.ID
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Error("could not find role by name after creation", logs.String("name", e.Name))
				return err
			}
		}
		for _, p := range e.Permissions {
			pNew, err := db.GetPermissionByName(context.Background(), p.Name)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					log.Error("could not find permission by name for attaching", logs.String("name", p.Name))
					return err
				}
			}
			//fmt.Println(e)
			//fmt.Println(pNew)
			//fmt.Println()
			ok, err := db.IsRolePermissionAttachExists(context.Background(), e.ID, pNew.ID)
			if err != nil {
				log.Error("could not get whether role and permission already attached", logs.Error(err))
				return err
			}
			if !ok {
				if err = db.Attach(context.Background(), e.ID, pNew.ID); err != nil {
					if !errors.Is(err, storage.ErrAlreadyExists) {
						log.Error("could not attach permission to role", logs.Error(err))
						return err
					}
				}
			}
		}
	}

	return nil
}

func defaultUser(db storage.UserInterface, log logs.LoggerInterface) error {
	superPWD := os.Getenv("ROOT_PWD")
	superEmail := os.Getenv("ROOT_EMAIL")

	pwd, err := auth_lib.GetHashPassword(superPWD)
	if err != nil {
		log.Error("could not generate hash password", logs.Error(err))
		return err
	}

	if err := db.CreateStaff(context.Background(), models.Staff{
		ID:          uuid.New().String(),
		Name:        "Super",
		PhoneNumber: sql.NullString{Valid: false},
		Email:       sql.NullString{Valid: true, String: superEmail},
		Password:    pwd,
		RoleID:      auth_lib.RoleSuper.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   sql.NullTime{Valid: false},
	}); err != nil {
		if !errors.Is(err, storage.ErrAlreadyExists) {
			log.Error("error to create staff in db", logs.Error(err))
			return err
		}
	}

	return nil
}
