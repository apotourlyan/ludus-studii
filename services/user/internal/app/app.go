package app

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/apotourlyan/ludus-studii/pkg/dbutil"
	"github.com/apotourlyan/ludus-studii/pkg/envutil"
	"github.com/apotourlyan/ludus-studii/pkg/httputil"
	"github.com/apotourlyan/ludus-studii/pkg/httputil/middleware"
	"github.com/apotourlyan/ludus-studii/pkg/idutil"
	"github.com/apotourlyan/ludus-studii/pkg/passutil"
	"github.com/apotourlyan/ludus-studii/pkg/secretutil"
	"github.com/apotourlyan/ludus-studii/pkg/syncutil"
	"github.com/apotourlyan/ludus-studii/pkg/timeutil"
	"github.com/apotourlyan/ludus-studii/services/user/internal/handler"
	"github.com/apotourlyan/ludus-studii/services/user/internal/repository"
	"github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register"
)

type App interface {
	Run()
	Expose() (http.Handler, *sql.DB)
}

type app struct {
	server *httputil.Server
	db     *sql.DB
}

type envars struct {
	MachineID       int64
	Port            int
	ShutdownTimeout time.Duration
}

type secrets struct {
	ConnectionString string
}

type handlers struct {
	Register handler.Post
}

func New(ep envutil.Provider, sp secretutil.Provider) App {
	variables := &envars{
		MachineID:       envutil.MachineID(ep).Value(),
		Port:            envutil.Port(ep).Value(),
		ShutdownTimeout: envutil.ShutdownTimeout(ep).Value(),
	}

	secrets := &secrets{
		ConnectionString: secretutil.ConnectionString(sp).Value(),
	}

	db := dbutil.Initialize(secrets.ConnectionString)

	handlers := createHandlers(db, variables)

	server := httputil.NewServer(
		&httputil.ServerConfig{
			Port:            ":" + strconv.Itoa(variables.Port),
			ShutdownTimeout: variables.ShutdownTimeout,
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			IdleTimeout:     60 * time.Second,
		},
	)

	server.AddMiddleware(middleware.CorrelationID)

	server.AddEndpoint("POST /api/register", handlers.Register.Execute)

	return &app{server, db}
}

func createHandlers(db *sql.DB, envars *envars) *handlers {
	timeProvider := timeutil.NewProvider()
	counter := syncutil.NewCounter()
	idgen := idutil.NewGenerator(timeProvider, counter, envars.MachineID)
	hasher := passutil.NewHasher()

	userRepository := repository.NewUserRepository(db)
	registerService := register.NewService(userRepository, idgen, hasher)
	registerHandler := handler.NewRegisterHandler(registerService)

	return &handlers{
		Register: registerHandler,
	}
}

func (a *app) Expose() (http.Handler, *sql.DB) {
	return a.server.Handler(), a.db
}

func (a *app) Run() {
	defer a.db.Close()
	a.server.Run()
}
