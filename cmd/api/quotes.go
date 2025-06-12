package main

import (
	"FIDOtestBackendApp/internal/store"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var (
	ValidationError = errors.New("validation error")
)

type CreateQuotePayload struct {
	Author string `json:"author" validate:"required,max=200"`
	Text   string `json:"text" validate:"required,max=1000"`
}

type UpdateQuotePayload struct {
	Author string `json:"author" validate:"required,max=200"`
	Text   string `json:"text" validate:"required,max=1000"`
}

func (app *application) createQuoteHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateQuotePayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	quote := &store.Quote{
		Author: payload.Author,
		Text:   payload.Text,
	}
	ctx := r.Context()
	app.logger.Info("Create Quote Handler", quote)
	if err := app.store.Quotes.CreateQuote(ctx, quote); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, quote); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) updateQuoteHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdateQuotePayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	quote, err := app.getQuoteByID(r)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestError(w, r, err)
		case errors.Is(err, ValidationError):
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	updatedQuote := &store.Quote{
		ID:     quote.ID,
		Author: payload.Author,
		Text:   payload.Text,
	}
	ctx := r.Context()
	app.logger.Info(updatedQuote.ID, updatedQuote.Author, updatedQuote.Text)
	err = app.store.Quotes.UpdateQuote(ctx, updatedQuote)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, updatedQuote); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getQuoteHandler(w http.ResponseWriter, r *http.Request) {

	quote, err := app.getQuoteByID(r)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestError(w, r, err)
		case errors.Is(err, ValidationError):
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, quote); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deleteQuoteHandler(w http.ResponseWriter, r *http.Request) {
	quote, err := app.getQuoteByID(r)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestError(w, r, err)
		case errors.Is(err, ValidationError):
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	ctx := r.Context()
	err = app.store.Quotes.DeleteQuote(ctx, quote.ID)
	app.logger.Info("Error", err)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func (app *application) getQuoteByID(r *http.Request) (*store.Quote, error) {
	quoteID := chi.URLParam(r, "quoteID")
	id, err := strconv.ParseInt(quoteID, 10, 64)
	if err != nil {
		return nil, ValidationError
	}
	quote, err := app.store.Quotes.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return nil, store.ErrNotFound
		default:
			return nil, err
		}
	}
	return quote, nil

}

func (app *application) getPaginatedQuoteList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	author := r.URL.Query().Get("author")
	filterDefault := store.PaginatedQuery{
		Limit:  10,
		Offset: 0,
	}
	filterQuery, err := filterDefault.Parse(r)
	if err != nil {
		app.badRequestError(w, r, ValidationError)
		return
	}

	if err := Validate.Struct(filterQuery); err != nil {
		app.badRequestError(w, r, ValidationError)
		return
	}
	quotes, err := app.store.Quotes.GetByAuthor(ctx, author, filterQuery)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, quotes); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
