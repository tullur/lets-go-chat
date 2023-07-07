package service

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/tullur/lets-go-chat/internal/domain/chat/token/memory"
)

func TestChatService(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Chat Service Operations", func() {
		cs, err := NewChatService(WithInMemoryTokenRepository())
		if err != nil {
			g.Fail(err)
		}

		us, err := NewUserService(WithInMemoryRepository())
		if err != nil {
			g.Fail(err)
		}

		user, err := us.CreateUser("Test", "Password1234")
		if err != nil {
			g.Fail(err)
		}

		chatToken, err := cs.GenerateAccessToken(user)
		if err != nil {
			g.Fail(err)
		}

		g.Describe("GenerateAccessToken()", func() {
			g.It("Returns correct user association", func() {
				g.Assert(chatToken.User()).Equal(user.Id())
			})
		})

		g.Describe("GetToken()", func() {
			result, err := cs.GetToken(chatToken.Id())
			if err != nil {
				g.Fail(err)
			}

			g.It("Returns correct token", func() {
				g.Assert(result).Equal(chatToken)
			})

			_, err = cs.GetToken("test-no-existing-id")
			g.It("Returns error", func() {
				g.Assert(err).Equal(memory.ErrTokenNotExists)
			})
		})

		g.Describe("RevokeToken()", func() {
			err := cs.RevokeToken(chatToken.Id())

			g.It("Should revoke token", func() {
				g.Assert(err).Equal(nil)
			})
		})
	})
}
