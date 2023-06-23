package service

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestUserService(t *testing.T) {
	g := Goblin(t)

	g.Describe("User Service", func() {
		us, err := NewUserService(WithInMemoryRepository())
		if err != nil {
			g.Fail(err)
		}

		user, err := us.CreateUser("Test", "Password1234")

		g.Describe("When create valid user", func() {
			g.It("Returns user", func() {
				g.Assert(user.Name()).Equal("Test")
			})

			g.It("Return nil error", func() {
				g.Assert(err).Equal(nil)
			})
		})

		g.Describe("When login correct user", func() {
			loginUser, err := us.LoginUser(user.Name(), "Password1234")

			g.It("Returns correct user", func() {
				g.Assert(loginUser.Name()).Equal(user.Name())
			})

			g.It("Should not return errors", func() {
				g.Assert(err).Equal(nil)
			})
		})

		g.Describe("GetList()", func() {
			g.It("Returns correct users length", func() {
				g.Assert(len(us.GetList())).Equal(1)
			})
		})
	})
}
