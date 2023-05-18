package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/tullur/lets-go-chat/internal/domain/user"
	"github.com/tullur/lets-go-chat/pkg/hasher"
)

func HandleUserList(userRepo user.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(userRepo.List())
	}
}

func HandleUserCreation(userRepo user.Repository) http.HandlerFunc {
	type userResponse struct {
		Id   string `json:"id"`
		Name string `json:"userName"`
	}

	type validatationError struct {
		Err string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Expires-After", time.August.String())
		w.Header().Set("X-Rate-Limit", strconv.Itoa(rand.Intn(20)))

		request := map[string]string{}

		json.NewDecoder(r.Body).Decode(&request)

		if len(request["userName"]) == 0 || len(request["password"]) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validatationError{Err: "Empty username or password"})

			return
		}

		hashedPassword, err := hasher.HashPassword(request["password"])
		if err != nil {
			log.Println(err)
		}

		createdUser := userRepo.Create(
			user.User{
				ID:       uuid.New(),
				Name:     request["userName"],
				Password: hashedPassword,
			},
		)

		responseBody := userResponse{
			Id:   createdUser.ID.String(),
			Name: createdUser.Name,
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(responseBody)
	}
}

func HandleUserLogin(userRepo user.Repository) http.HandlerFunc {
	type successResponse struct {
		Url string `json:"url"`
	}

	type errorResponse struct {
		Err string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request := map[string]string{}

		json.NewDecoder(r.Body).Decode(&request)

		currentUser, err := userRepo.FindByName(request["userName"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Err: err.Error()})

			return
		}

		checker := hasher.CheckPasswordHash(request["password"], currentUser.Password)
		if !checker {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Err: "Invalid username/password"})

			return
		}

		responseBody := successResponse{
			Url: "ws://fancy-chat.io/ws&token=one-time-token",
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}
}
