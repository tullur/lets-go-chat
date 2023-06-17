package ws

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/franela/goblin"
	"github.com/gorilla/websocket"
	"github.com/tullur/lets-go-chat/internal/service"
)

func TestChat(t *testing.T) {
	g := goblin.Goblin(t)

	// Mock it
	tokenService, err := service.NewTokenService(service.WithInMemoryTokenRepository())
	if err != nil {
		g.Fail(err)
	}

	userService, err := service.NewUserService(service.WithInMemoryRepository())
	if err != nil {
		g.Fail(err)
	}

	user, err := userService.CreateUser("TestUser", "Password1234")
	if err != nil {
		g.Fail(err)
	}

	currentUser, err := userService.LoginUser(user.Name(), "Password1234")
	if err != nil {
		g.Fail(err)
	}

	g.Describe("[GET] /v1/chat/users", func() {
		req, err := http.NewRequest(http.MethodGet, "/v1/chat/users", nil)
		if err != nil {
			g.Fail(err)
		}

		res := httptest.NewRecorder()

		handler := GetChatUsers()
		handler.ServeHTTP(res, req)

		g.It("Returns success status code", func() {
			g.Assert(res.Code).Equal(http.StatusOK)
		})
	})

	g.Describe("[GET] /v1/chat", func() {
		g.Describe("Unathorized error", func() {
			req, err := http.NewRequest(http.MethodGet, "/v1/chat", nil)
			if err != nil {
				g.Fail(err)
			}

			res := httptest.NewRecorder()

			handler := ChatConnection(tokenService)
			handler.ServeHTTP(res, req)

			g.It("Returns unauthorized error", func() {
				g.Assert(res.Code).Equal(http.StatusUnauthorized)
			})
		})

		g.Describe("When invalid token", func() {
			req, err := http.NewRequest(http.MethodGet, "/v1/chat?token=fake", nil)
			if err != nil {
				g.Fail(err)
			}

			res := httptest.NewRecorder()

			handler := ChatConnection(tokenService)
			handler.ServeHTTP(res, req)

			g.It("Returns unauthorized error", func() {
				g.Assert(res.Code).Equal(http.StatusUnauthorized)
			})
		})

		g.Describe("Chat Connection", func() {
			loginToken, err := tokenService.GenerateAccessToken(currentUser)
			if err != nil {
				g.Fail(err)
			}

			url := fmt.Sprintf("/v1/chat?token=%s", loginToken.Id())

			go BroadcastMessages()

			s := httptest.NewServer(http.HandlerFunc(ChatConnection(tokenService)))
			defer s.Close()

			u := "ws" + strings.TrimPrefix(s.URL, "http")

			log.Println(u + url)

			ws, _, err := websocket.DefaultDialer.Dial(u+url, nil)
			if err != nil {
				t.Fatalf("%v", err)
			}

			defer ws.Close()

			// Send message to server, read response and check to see if it's what we expect.
			for i := 0; i < 10; i++ {
				if err := ws.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
					t.Fatalf("%v", err)
				}
				_, p, err := ws.ReadMessage()
				if err != nil {
					t.Fatalf("%v", err)
				}
				log.Println(p)
				// if string(p) != "hello" {
				// 	t.Fatalf("bad message")
				// }
			}

		})
	})
}
