package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/tullur/lets-go-chat/internal/service"
)

var request map[string]string

type ErrorResponse struct {
	Err string `json:"error"`
}

func HandleUserList(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(userService.GetList())
	}
}

func HandleUserCreation(userService service.UserService) http.HandlerFunc {
	type userResponse struct {
		Id   string `json:"id"`
		Name string `json:"userName"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&request)

		user, err := userService.CreateUser(request["userName"], request["password"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Err: err.Error()})

			return
		}

		responseBody := userResponse{
			Id:   user.ID.String(),
			Name: user.Name,
		}

		w.WriteHeader(http.StatusCreated)

		w.Header().Set("X-Expires-After", time.August.String())
		w.Header().Set("X-Rate-Limit", strconv.Itoa(rand.Intn(20)))
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(responseBody)
	}
}

func HandleUserLogin(userService service.UserService) http.HandlerFunc {
	type successResponse struct {
		Url string `json:"url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&request)

		user, err := userService.LoginUser(request["userName"], request["password"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Err: err.Error()})

			return
		}

		responseBody := successResponse{
			Url: fmt.Sprintf("ws://fancy-chat.io/ws&token=%s", user.ID.String()),
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}
}
