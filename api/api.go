package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Id uuid.UUID

type User struct {
	Id        Id     `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type Application struct {
	Data map[Id]User
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func NewHandler(db *Application) http.Handler {
	r := chi.NewMux()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Post("/api/users", Insert(db))
	r.Get("/api/users", FindAll(db))
	r.Get("/api/users/{id}", FindById(db))
	r.Put("/api/users/{id}", Update(db))
	r.Delete("/api/users/{id}", Delete(db))
	return r
}

func Insert(db *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bodyUser User

		if err := json.NewDecoder(r.Body).Decode(&bodyUser); err != nil {
			sendJSON(w, Response{Error: "invalid Body"}, http.StatusUnprocessableEntity)
			return
		}

		if db == nil || db.Data == nil {
			sendJSON(w, Response{Error: "Internal Server Error"}, http.StatusInternalServerError)
			return
		}

		if len(strings.TrimSpace(bodyUser.FirstName)) < 2 || len(strings.TrimSpace(bodyUser.FirstName)) > 20 {
			sendJSON(w, Response{Error: "O nome precisa ter no máximo 20 caracteres e no mínimo 2."}, http.StatusBadRequest)
			return
		}

		if len(strings.TrimSpace(bodyUser.LastName)) < 2 || len(strings.TrimSpace(bodyUser.LastName)) > 20 {
			sendJSON(w, Response{Error: "O sobrenome precisa ter no máximo 20 caracteres e no mínimo 2."}, http.StatusBadRequest)
			return
		}

		if len(strings.TrimSpace(bodyUser.Biography)) < 20 || len(strings.TrimSpace(bodyUser.Biography)) > 450 {
			sendJSON(w, Response{Error: "A biografia precisa ter no mínimo 20 caracteres e no máximo 450."}, http.StatusBadRequest)
			return
		}

		userID := uuid.New()
		bodyUser.Id = Id(userID)
		db.Data[Id(userID)] = bodyUser
		sendJSON(w, Response{Data: bodyUser}, http.StatusCreated)
	}

}

func FindAll(db *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if db == nil || db.Data == nil {
			sendJSON(w, Response{Error: "Internal Server Error"}, http.StatusInternalServerError)
			return
		}

		users := make([]User, 0, len(db.Data))

		for _, value := range db.Data {
			users = append(users, value)
		}
		sendJSON(w, Response{Data: users}, http.StatusOK)
	}
}

func FindById(db *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		idParsed, err := uuid.Parse(idStr)

		if err != nil {
			sendJSON(w, Response{Error: "invalid user id"}, http.StatusBadRequest)
			return
		}

		if db == nil || db.Data == nil {
			sendJSON(w, Response{Error: "internal server error"}, http.StatusInternalServerError)
			return
		}

		id := Id(idParsed)
		user, ok := db.Data[id]
		if !ok {
			sendJSON(w, Response{Error: "user not found."}, http.StatusNotFound)
			return
		}
		sendJSON(w, Response{Data: user}, http.StatusOK)

	}

}

func Update(db *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idSrt := chi.URLParam(r, "id")
		idParsed, err := uuid.Parse(idSrt)

		if err != nil {
			sendJSON(w, Response{Error: "invalid user id"}, http.StatusBadRequest)
			return
		}

		if db == nil || db.Data == nil {
			sendJSON(w, Response{Error: "internal server error"}, http.StatusInternalServerError)
			return
		}

		id := Id(idParsed)
		_, ok := db.Data[id]
		if !ok {
			sendJSON(w, Response{Error: "user not found"}, http.StatusNotFound)
			return
		}

		var bodyUser User

		if err := json.NewDecoder(r.Body).Decode(&bodyUser); err != nil {
			sendJSON(w, Response{Error: "invalid Body"}, http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(bodyUser.FirstName) == "" || strings.TrimSpace(bodyUser.LastName) == "" || strings.TrimSpace(bodyUser.Biography) == "" {
			sendJSON(w, Response{Error: "missing required fields"}, http.StatusBadRequest)
			return
		}

		bodyUser.Id = id
		db.Data[id] = bodyUser

		sendJSON(w, Response{Data: bodyUser}, http.StatusOK)

	}
}

func Delete(db *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		idParsed, err := uuid.Parse(idStr)

		if err != nil {
			sendJSON(w, Response{Error: "invalid user id"}, http.StatusBadRequest)
			return
		}

		if db == nil || db.Data == nil {
			sendJSON(w, Response{Error: "internal server error"}, http.StatusInternalServerError)
			return
		}

		id := Id(idParsed)
		user, ok := db.Data[id]
		if !ok {
			sendJSON(w, Response{Error: "user not found"}, http.StatusNotFound)
			return
		}

		delete(db.Data, id)
		sendJSON(w, Response{Data: user}, http.StatusOK)

	}
}

func sendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("error ao fazer marshal de json", "error", err)
		sendJSON(w, Response{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("error ao enviar a resposta:", "error", err)
		return
	}
}
