package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franela/goblin"
	"github.com/tullur/lets-go-chat/internal/domain/user"
	"github.com/tullur/lets-go-chat/internal/service"
)

type mockUserService struct{}

func (s *mockUserService) GetList() []user.User {
	return []user.User{
		{},
		{},
		{},
	}
}

func TestUserHandlers(t *testing.T) {
	g := goblin.Goblin(t)

	// Mock it
	userService, err := service.NewUserService(service.WithInMemoryRepository())
	if err != nil {
		g.Fail(err)
	}

	g.Describe("GetList()", func() {
		req, err := http.NewRequest("GET", "/v1/user", nil)
		if err != nil {
			g.Fail(err)
		}

		res := httptest.NewRecorder()

		handler := GetUsers(userService)
		handler.ServeHTTP(res, req)

		g.It("Returns success status code", func() {
			g.Assert(res.Code).Equal(http.StatusOK)
		})
	})

	g.Describe("Create()", func() {
		payload := struct {
			UserName string `json:"userName"`
			Password string `json:"password"`
		}{
			UserName: "TestUser",
			Password: "QwertyTest1234",
		}

		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(payload)

		req, err := http.NewRequest(http.MethodPost, "/v1/user", &buf)
		if err != nil {
			g.Fail(err)
		}

		res := httptest.NewRecorder()

		handler := CreateUser(userService)
		handler.ServeHTTP(res, req)

		g.It("Returns created status code", func() {
			g.Assert(res.Code).Equal(http.StatusCreated)
		})

		g.Describe("Internal Server Error", func() {
			payload := struct {
				UserName string `json:"userName"`
				Password string `json:"password"`
			}{
				UserName: "",
				Password: "",
			}

			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(payload)

			req, err := http.NewRequest(http.MethodPost, "/v1/user", &buf)
			if err != nil {
				g.Fail(err)
			}

			res := httptest.NewRecorder()

			handler := CreateUser(userService)
			handler.ServeHTTP(res, req)

			g.It("Returns error status code", func() {
				g.Assert(res.Code).Equal(http.StatusInternalServerError)
			})
		})
	})

	g.Describe("Login()", func() {
		// Mock it
		tokenService, err := service.NewChatService(service.WithInMemoryTokenRepository())
		if err != nil {
			g.Fail(err)
		}

		payload := struct {
			UserName string `json:"userName"`
			Password string `json:"password"`
		}{
			UserName: "TestUser",
			Password: "QwertyTest1234",
		}

		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(payload)

		req, err := http.NewRequest(http.MethodPost, "/v1/user/login", &buf)
		if err != nil {
			g.Fail(err)
		}

		res := httptest.NewRecorder()

		handler := LoginUser(userService, tokenService)
		handler.ServeHTTP(res, req)

		g.It("Returns success status code", func() {
			g.Assert(res.Code).Equal(http.StatusOK)
		})

		g.Describe("Internal Server Error", func() {
			payload := struct {
				UserName string `json:"userName"`
				Password string `json:"password"`
			}{
				UserName: "",
				Password: "",
			}

			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(payload)

			req, err := http.NewRequest(http.MethodPost, "/v1/user/login", &buf)
			if err != nil {
				g.Fail(err)
			}

			res := httptest.NewRecorder()

			handler := LoginUser(userService, tokenService)
			handler.ServeHTTP(res, req)

			g.It("Returns error status code", func() {
				g.Assert(res.Code).Equal(http.StatusInternalServerError)
			})
		})
	})
}
