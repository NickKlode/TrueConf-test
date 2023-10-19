package api

import (
	"encoding/json"
	"log"
	"net/http"
	"refactoring/internal/api/response"
	"refactoring/internal/storage"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (api *API) endpoints() {
	api.r.Use(middleware.RequestID)
	api.r.Use(middleware.RealIP)
	api.r.Use(middleware.Logger)
	api.r.Use(middleware.Recoverer)

	api.r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	api.r.Get("/api/v1/users/", api.getAllUsers)
	api.r.Post("/api/v1/users/", api.createUser)

	api.r.Get("/api/v1/users/{id}", api.getUserByID)
	api.r.Patch("/api/v1/users/{id}", api.updateUser)
	api.r.Delete("/api/v1/users/{id}", api.deleteUser)

}

func (api *API) createUser(w http.ResponseWriter, r *http.Request) {

	var userInput storage.UserInput

	err := render.DecodeJSON(r.Body, &userInput)
	if err != nil {
		log.Printf("cannot decode request body. %v", err)
		render.JSON(w, r, response.Error("enter correct data"))
		return
	}
	if userInput.DisplayName == "" || userInput.Email == "" {
		log.Printf("incorrect data. %v", err)
		render.JSON(w, r, response.Error("enter correct data"))
		return
	}

	id, err := api.db.CreateUser(userInput)
	if err != nil {
		log.Printf("error while creating user. %v", err)
		render.JSON(w, r, response.Error("cannot create user"))
		return
	}
	render.JSON(w, r, response.OK(map[string]string{"id": id}))
}
func (api *API) deleteUser(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	err := api.db.DeleteUser(id)
	if err != nil {
		log.Printf("error while deleting user. %v", err)
		render.JSON(w, r, response.Error("cannot delete user"))
		return
	}
	render.JSON(w, r, response.OK(map[string]string{"Message": "user deleted"}))
}
func (api *API) getAllUsers(w http.ResponseWriter, r *http.Request) {

	usrs, err := api.db.GetAllUsers()
	if err != nil {
		if err != nil {
			log.Printf("error while getting all users. %v", err)
			render.JSON(w, r, response.Error("cannot get all users"))
			return
		}
	}
	render.JSON(w, r, response.OK(usrs))
}
func (api *API) updateUser(w http.ResponseWriter, r *http.Request) {

	var userInput storage.UserInput
	id := chi.URLParam(r, "id")
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		log.Printf("cannot decode request body. %v", err)
		render.JSON(w, r, response.Error("enter correct data"))
		return
	}
	if userInput.DisplayName == "" && userInput.Email == "" {
		log.Printf("incorrect data. %v", err)
		render.JSON(w, r, response.Error("enter correct data"))
		return
	}
	err = api.db.UpdateUser(id, userInput)
	if err != nil {
		log.Printf("error while updating user. %v", err)
		render.JSON(w, r, response.Error("cannot update user"))
		return
	}
	render.JSON(w, r, response.OK(map[string]string{"Message": "user updated"}))
}
func (api *API) getUserByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	usr, err := api.db.GetUserByID(id)
	if err != nil {
		log.Printf("error while getting user by id. %v", err)
		render.JSON(w, r, response.Error("cannot get user by id"))
		return
	}
	render.JSON(w, r, response.OK(usr))
}
