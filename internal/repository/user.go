package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mestvv/NorthBridgeBackend/internal/db"
	"github.com/mestvv/NorthBridgeBackend/internal/domain"
)

type userRepository struct {
	db *sqlx.DB
}

func newUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	const query = `
				INSERT INTO user (id, login, password, email) 
				VALUES(uuid_to_bin(?), ?, ?, ?);
				`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Login, user.Password, user.Email)

	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok && mysqlError.Number == db.DuplicateEntry {
			return domain.ErrDuplicateEntry
		}
		return fmt.Errorf("db insert user: %v", err)
	}

	return nil
}

func (r *userRepository) GetByCredentials(ctx context.Context, email string, password string) (*uuid.UUID, error) {
	const query = `
				SELECT id FROM user
				WHERE email = ?
				AND password = ?
				`
	var ID uuid.UUID
	if err := r.db.GetContext(ctx, &ID, query, email, password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("select from user failed: %w", err)
	}

	return &ID, nil
}
