package dbutil

import (
	"context"
	"database/sql"

	"github.com/apotourlyan/ludus-studii/pkg/typeutil"
)

// Executor interface - both *sql.DB and *sql.Tx implement this
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Query runs function with executor
// - In test environment: uses transaction from context (if present)
// - Otherwise: uses DB connection pool
func (r *Repository) Query(ctx context.Context, fn func(Executor) error) error {
	// Check for test transaction in context
	if tx, ok := ctx.Value(typeutil.KeyTestTx).(*sql.Tx); ok {
		// Use test transaction
		return fn(tx)
	}

	// Use DB directly
	return fn(r.db)
}

// Command runs function in transaction
// - In test environment: uses transaction from context (doesn't commit)
// - In production: creates new transaction and commits
func (r *Repository) Command(ctx context.Context, fn func(Executor) error) error {
	// Check for test transaction in context
	if tx, ok := ctx.Value(typeutil.KeyTestTx).(*sql.Tx); ok {
		// Use test transaction, DON'T commit
		return fn(tx)
	}

	// Create production transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute function
	if err := fn(tx); err != nil {
		return err
	}

	// Commit in production
	return tx.Commit()
}
