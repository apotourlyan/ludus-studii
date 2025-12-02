package main

import (
	"database/sql"
	"os"
	"strconv"
	"time"

	"github.com/apotourlyan/ludus-studii/pkg/dbutil"
	"github.com/apotourlyan/ludus-studii/pkg/envutil"
	"github.com/apotourlyan/ludus-studii/pkg/httputil"
	"github.com/apotourlyan/ludus-studii/pkg/idutil"
	"github.com/apotourlyan/ludus-studii/pkg/passutil"
	"github.com/apotourlyan/ludus-studii/pkg/syncutil"
	"github.com/apotourlyan/ludus-studii/pkg/timeutil"
	"github.com/apotourlyan/ludus-studii/services/user/internal/handler"
	"github.com/apotourlyan/ludus-studii/services/user/internal/repository"
	"github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register"
)

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

func main() {
	variables := getEnvironmentVariables()
	secrets := getSecrets()

	db := dbutil.Initialize(secrets.ConnectionString)
	defer db.Close()

	handlers := createHandlers(db, variables)

	httpserver := httputil.NewServer(
		&httputil.ServerConfig{
			Port:            ":" + strconv.Itoa(variables.Port),
			ShutdownTimeout: variables.ShutdownTimeout,
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			IdleTimeout:     60 * time.Second,
		},
	)

	httpserver.AddEndpoint("POST /api/register", handlers.Register.Execute)

	httpserver.Run()
}

func getSecrets() *secrets {
	// TODO: Get from a dedicated secrets' vault
	return &secrets{
		ConnectionString: os.Getenv("DATABASE_URL"),
	}
}

func getEnvironmentVariables() *envars {
	p := envutil.NewProvider()

	return &envars{
		MachineID:       envutil.MachineID(p).Value(),
		Port:            envutil.Port(p).Value(),
		ShutdownTimeout: envutil.ShutdownTimeout(p).Value(),
	}
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
