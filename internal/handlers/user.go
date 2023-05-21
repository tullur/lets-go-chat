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

type RequestBody struct {
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
}

var requestBody RequestBody

func HandleUserList(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(userService.GetList())
	}
}

func HandleUserCreation(userService *service.UserService) http.HandlerFunc {
	type userResponse struct {
		Id   string `json:"id"`
		Name string `json:"userName"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&requestBody)

		user, err := userService.Create(requestBody.UserName, requestBody.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseBody := userResponse{
			Id:   user.ID.String(),
			Name: user.Name,
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(responseBody)
	}
}

func HandleUserLogin(userService *service.UserService) http.HandlerFunc {
	type successResponse struct {
		Url string `json:"url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&requestBody)

		user, err := userService.Login(requestBody.UserName, requestBody.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseBody := successResponse{
			Url: fmt.Sprintf("ws://fancy-chat.io/ws&token=%s", user.ID.String()),
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Expires-After", time.August.String())
		w.Header().Set("X-Rate-Limit", strconv.Itoa(rand.Intn(20)))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}
}
