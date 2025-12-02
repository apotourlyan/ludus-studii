package repository

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/apotourlyan/ludus-studii/pkg/errorutil"
	"github.com/apotourlyan/ludus-studii/services/user/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type postgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

//go:embed scripts/user_insert.sql
var insertScript string

//go:embed scripts/user_exists_by_email.sql
var existsByEmailScript string

func (r *postgresUserRepository) Create(
	ctx context.Context, user *domain.User,
) error {
	_, err := r.db.ExecContext(
		ctx,
		insertScript,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
	)

	return errorutil.DatabaseError(err)
}

func (r *postgresUserRepository) ExistsByEmail(
	ctx context.Context, email string,
) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, existsByEmailScript, email).Scan(&exists)
	return exists, errorutil.DatabaseError(err)
}
