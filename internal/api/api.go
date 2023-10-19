package api

import (
	"refactoring/internal/storage"

	"github.com/go-chi/chi/v5"
)

type API struct {
	db storage.StorageInterface
	r  *chi.Mux
}

func New(db storage.StorageInterface) *API {
	a := API{db: db, r: chi.NewRouter()}
	a.endpoints()
	return &a
}

func (api *API) Router() *chi.Mux {
	return api.r
}
