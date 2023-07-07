package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tullur/lets-go-chat/internal/service"
)

type (
	RequestParams struct {
		UserName string `json:"userName,omitempty" example:"Gopher"`
		Password string `json:"password,omitempty" example:"Qwerty1234"`
	}

	createUserResponse struct {
		Id   string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
		Name string `json:"userName" example:"Gopher"`
	}

	loginUserResponse struct {
		Url string `json:"url" example:"ws://localhost:8080/v1/chat?token=b8606575-4593-424f-9f70-11b5dce79d54"`
	}
)

var requestParams RequestParams

// GetUsers lists all existing users
//
//	@Summary      User List
//	@Description  get users
//	@Tags         users
//	@Produce      json
//	@Success      200
//	@Router       /users [get]
func GetUsers(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(userService.GetList())
	}
}

// CreateUser
//
// @Summary      Creat user
// @Description  post user
// @Tags         users
// @Param	     user	body		RequestParams	true	"Login user"
// @Accept       json
// @Produce      json
// @Success		 201	{object}	createUserResponse
// @Failure      400
// @Failure      500
// @Router       /users [post]
func CreateUser(userService *service.UserService) http.HandlerFunc {
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

// LoginUser
//
// @Summary      Login user
// @Description  Provide login token to the chat service
// @Tags         users
// @Param	     user	body		RequestParams	true	"Add user"
// @Accept       json
// @Produce      json
// @Success		 201	{object}	loginUserResponse
// @Failure      400
// @Failure      500
// @Router       /login [post]
func LoginUser(userService *service.UserService, chatService *service.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&requestParams)

		user, err := userService.LoginUser(requestParams.UserName, requestParams.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := chatService.GenerateAccessToken(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseBody := loginUserResponse{
			Url: fmt.Sprintf("ws://%s/v1/chat?token=%s", r.Host, token.Id()),
		}

		w.Header().Set("X-Expires-After", token.ExpiresAfter())
		w.Header().Set("X-Rate-Limit", strconv.Itoa(100))
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}
}
