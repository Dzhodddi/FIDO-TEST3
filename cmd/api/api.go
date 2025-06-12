package main

import (
	"FIDOtestBackendApp/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type application struct {
	config config
	logger *zap.SugaredLogger
	store  store.Storage
}

type dbConfig struct {
	addr               string
	maxOpenConnections int
	maxIdleConnections int
	maxIdleTime        string
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Minute,
	}
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	//versioned??
	r.Route("/", func(r chi.Router) {
		r.Get("/ping", app.healthCheckHandler)
		r.Route("/quotes", func(r chi.Router) {
			r.Post("/", app.createQuoteHandler)
			r.Get("/", app.getPaginatedQuoteList)
			r.Route("/{quoteID}", func(r chi.Router) {
				r.Get("/", app.getQuoteHandler)
				r.Delete("/", app.deleteQuoteHandler)
				r.Put("/", app.updateQuoteHandler)
			})
		})

	})
	return r
}
