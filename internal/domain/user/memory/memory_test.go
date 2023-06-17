package memory

import (
	"reflect"
	"testing"

	"github.com/franela/goblin"
	"github.com/tullur/lets-go-chat/internal/domain/user"
)

var repo = NewInMemoryRepository()

func TestNewInMemoryRepository(t *testing.T) {
	tests := []struct {
		name string
		want *InMemoryRepo
	}{
		{
			name: "Correct Initialization",
			want: NewInMemoryRepository(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInMemoryRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInMemoryRepository() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestInMemoryRepositoryOperations(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("InMemory Repository Operations", func() {
		testUser, err := user.New("Goblin", "Qwerty12345678")
		if err != nil {
			g.Fail(err)
		}

		g.It("Should create user", func() {
			g.Assert(len(repo.Users)).Equal(0)

			repo.Create(testUser)

			g.Assert(len(repo.Users)).Equal(1)
		})

		g.Describe("FindByName()", func() {
			g.It("Should find user by name", func() {
				expect, err := repo.FindByName("Goblin")
				if err != nil {
					g.Fail(err)
				}

				g.Assert(expect).Equal(testUser)
			})

			g.It("Returns error when user not found", func() {
				_, err := repo.FindByName("DummyUser")

				g.Assert(err).Equal(user.ErrUserNotFound)
			})
		})

		g.It("Returns list of users", func() {
			g.Assert(len(repo.List())).Equal(1)
		})
	})
}
