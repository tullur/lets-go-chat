package token

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/google/uuid"
)

var userUUID = uuid.New()

func TestToken(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("New Token", func() {
		token := New(userUUID.String())

		g.It("Should create token", func() {
			g.Assert(token.User()).Equal(userUUID.String())
			g.Assert(len(token.Id())).Equal(36)
			g.Assert(token.ExpiresAfter()).Equal(token.expiresAfter.Local().String())
		})
	})
}
