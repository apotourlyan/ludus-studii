package integration

import (
	"database/sql"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/envutil"
	"github.com/apotourlyan/ludus-studii/pkg/envutil/envvar"
	"github.com/apotourlyan/ludus-studii/pkg/secretutil"
	"github.com/apotourlyan/ludus-studii/pkg/secretutil/secretvar"
	"github.com/apotourlyan/ludus-studii/pkg/syncutil"
	"github.com/apotourlyan/ludus-studii/pkg/testutil/txutil"
	"github.com/apotourlyan/ludus-studii/services/user/internal/app"
)

var (
	server     *httptest.Server
	db         *sql.DB
	txRegistry = syncutil.NewConcurrentMap[string, *sql.Tx]()
)

func getPath(path string) string {
	return server.URL + path
}

func TestMain(m *testing.M) {
	log.Println("Setting up test environment...")

	// Set environment variables
	os.Setenv(secretvar.DbConnection, "postgres://ludus:ludus@db:5432/ludus_studii_user_db_test")
	os.Setenv(envvar.MachineId, "1")
	os.Setenv(envvar.Port, "8080")
	os.Setenv(envvar.ShutdownTimeout, "5s")

	// Setup
	ep := envutil.NewProvider()
	sp := secretutil.NewProvider()
	app := app.New(ep, sp)

	handler, tdb := app.Expose()
	db = tdb
	if db != nil {
		log.Println("Database initialized successfully")
	}

	// Provide test transactions to the HTTP server
	txutil.SetTxRegistry(txRegistry)
	handler = txutil.TxMiddleware(handler)

	tsrv := httptest.NewServer(handler)
	server = tsrv
	log.Printf("Test server started at: %s\n", server.URL)

	log.Println("Running tests...")
	// Run tests
	code := m.Run()

	// Cleanup
	log.Println("Cleaning up...")
	tsrv.Close()
	tdb.Close()

	os.Exit(code)
}
