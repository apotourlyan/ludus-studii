package txutil

import (
	"context"
	"database/sql"
	"net/http"
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/httputil/header"
	"github.com/apotourlyan/ludus-studii/pkg/panicutil"
	"github.com/apotourlyan/ludus-studii/pkg/syncutil"
	"github.com/apotourlyan/ludus-studii/pkg/typeutil"
)

var (
	txRegistry syncutil.ConcurrentMap[string, *sql.Tx]
)

func SetTxRegistry(value syncutil.ConcurrentMap[string, *sql.Tx]) {
	txRegistry = value
}

func TxMiddleware(next http.Handler) http.Handler {
	panicutil.RequireNotNil(txRegistry, "transactions registry")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txID := r.Header.Get(header.CorrelationID)
		if txID != "" {
			if tx, ok := txRegistry.Get(txID); ok {
				ctx := context.WithValue(r.Context(), typeutil.KeyTestTx, tx)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func TxTest(t *testing.T, db *sql.DB, txID string, tfn func()) {
	t.Helper()
	panicutil.RequireNotNil(txRegistry, "transactions registry")

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin test transaction: %v", err)
	}
	defer tx.Rollback()

	// Register transaction
	txRegistry.Set(txID, tx)
	defer txRegistry.Remove(txID)

	// Run test
	tfn()
}

func TxQueryValue[T any](
	t *testing.T, txID string, query string, args ...any,
) T {
	t.Helper()
	panicutil.RequireNotNil(txRegistry, "transaction")
	panicutil.RequireNotEmptyOrWhitespace(txID, "transaction id")

	tx, ok := txRegistry.Get(txID)
	if !ok {
		t.Fatalf("Failed to get transaction with ID: %v", txID)
	}

	var result T
	err := tx.QueryRowContext(t.Context(), query, args...).Scan(&result)

	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	return result
}
