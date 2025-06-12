package main

import (
	"FIDOtestBackendApp/internal/db"
	"FIDOtestBackendApp/internal/store"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: ":8080",
		env:  "development",
		db: dbConfig{
			addr:               "postgresql://admin:yourpassword@localhost:5432/postgres?sslmode=disable",
			maxOpenConnections: 10,
			maxIdleConnections: 10,
			maxIdleTime:        "15m",
		},
	}

	// Logger init
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database init
	database, err := db.New(cfg.db.addr,
		cfg.db.maxOpenConnections,
		cfg.db.maxIdleConnections,
		cfg.db.maxIdleTime,
	)
	defer database.Close()
	if err != nil {
		logger.Fatal(err)
	}

	// Storage init
	storage := store.NewStorage(database)
	app := &application{
		config: cfg,
		logger: logger,
		store:  storage,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
