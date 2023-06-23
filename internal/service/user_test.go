package service

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/golang/mock/gomock"
	"github.com/tullur/lets-go-chat/internal/domain/user"
	"github.com/tullur/lets-go-chat/internal/domain/user/mocks"
)

var (
	ctrl *gomock.Controller
	g    *goblin.G
)

func TestUserService(t *testing.T) {
	g = goblin.Goblin(t)
	ctrl = gomock.NewController(t)

	defer ctrl.Finish()

	g.Describe("User Service", func() {
		mockRepo := mocks.NewMockRepository(ctrl)

		us, err := NewUserService(WithRepository(mockRepo))
		if err != nil {
			g.Fail(err)
		}

		mockRepo.EXPECT().Create(gomock.Any()).Times(1).Return(nil)

		testUser, err := us.CreateUser("Test", "Qwerty123454")

		g.Describe("When create valid user", func() {
			g.It("Returns user", func() {
				g.Assert(testUser.Name()).Equal("Test")
			})

			g.It("Return nil error", func() {
				g.Assert(err).Equal(nil)
			})
		})

		g.Describe("When login correct user", func() {
			mockRepo.EXPECT().FindByName(gomock.Any()).Times(1).Return(testUser, nil)
			loginUser, err := us.LoginUser(testUser.Name(), "Qwerty123454")

			g.It("Returns correct user", func() {
				g.Assert(loginUser.Name()).Equal(testUser.Name())
			})

			g.It("Should not return errors", func() {
				g.Assert(err).Equal(nil)
			})
		})

		g.Describe("GetList()", func() {
			mockRepo.EXPECT().List().Return([]user.User{*testUser})

			g.It("Returns correct users length", func() {
				g.Assert(len(us.GetList())).Equal(1)
			})
		})
	})
}
