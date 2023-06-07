package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/tullur/lets-go-chat/internal/domain/chat"
	"github.com/tullur/lets-go-chat/internal/service"
)

type RequestParams struct {
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
}

type createUserResponse struct {
	Id   string `json:"id"`
	Name string `json:"userName"`
}

type loginUserResponse struct {
	Url string `json:"url"`
}

var requestParams RequestParams

func HandleUserList(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(userService.GetList())
	}
}

func HandleUserCreation(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&requestParams)

		user, err := userService.CreateUser(requestParams.UserName, requestParams.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseBody := createUserResponse{
			Id:   user.Id(),
			Name: user.Name(),
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(responseBody)
	}
}

func HandleUserLogin(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&requestParams)

		user, err := userService.LoginUser(requestParams.UserName, requestParams.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token := chat.NewToken(user.Id())
		responseBody := loginUserResponse{
			Url: fmt.Sprintf("ws://%s/chat&token=%s", r.Host, token.Id()),
		}

		w.Header().Set("X-Expires-After", time.August.String())
		w.Header().Set("X-Rate-Limit", strconv.Itoa(rand.Intn(20)))
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}
}
