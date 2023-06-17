package memory

import (
	"reflect"
	"testing"

	"github.com/franela/goblin"
	"github.com/tullur/lets-go-chat/internal/domain/chat/token"
	"github.com/tullur/lets-go-chat/internal/domain/user"
)

var repo = NewMemoryTokenRepositorysitory()

func TestNewMemoryTokenRepositorysitory(t *testing.T) {
	tests := []struct {
		name string
		want *MemoryTokenRepository
	}{
		{
			name: "Correct Initialization",
			want: NewMemoryTokenRepositorysitory(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryTokenRepositorysitory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemoryTokenRepositorysitory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryTokenRepository_Get(t *testing.T) {
	user, err := user.New("TestCase", "Password1234")
	if err != nil {
		t.Fatal(err)
	}

	testToken := token.New(user.Id())
	repo.Add(testToken)

	defer repo.Revoke(testToken.Id())

	type fields struct {
		Tokens map[string]token.Token
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *token.Token
		wantErr bool
	}{
		{
			name:    "Get Token from repository",
			fields:  fields{repo.Tokens},
			args:    args{testToken.Id()},
			want:    testToken,
			wantErr: false,
		},
		{
			name:    "Get unexisted Token from repository",
			fields:  fields{repo.Tokens},
			args:    args{id: "dummy-token-id"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MemoryTokenRepository{
				Tokens: tt.fields.Tokens,
			}
			got, err := repo.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryTokenRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryTokenRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryTokenRepository_Add(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Add()", func() {
		user, err := user.New("TestCase", "Password1234")
		if err != nil {
			g.Fail(err)
		}

		testToken := token.New(user.Id())

		g.Describe("When add new token", func() {
			repo.Add(testToken)

			g.It("Adds new token to the repository", func() {
				g.Assert(len(repo.Tokens)).Equal(1)
			})
		})

		g.Describe("When add same token", func() {
			err := repo.Add(testToken)

			g.It("Returns error", func() {
				g.Assert(err).Equal(ErrTokenAlreadyProvided)
			})
		})
	})
}

func TestMemoryTokenRepository_Revoke(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Revoke()", func() {
		user, err := user.New("TestCase", "Password1234")
		if err != nil {
			g.Fail(err)
		}

		token := token.New(user.Id())
		repo.Add(token)

		g.Describe("When revoke unexisted token", func() {
			err := repo.Revoke("dummy-token-id")

			g.It("Should return error", func() {
				g.Assert(err).Equal(ErrTokeNotFound)
			})
		})

		g.It("Should revoke token from repository", func() {
			g.Assert(len(repo.Tokens)).Equal(2)

			err := repo.Revoke(token.Id())

			g.Assert(len(repo.Tokens)).Equal(1)
			g.Assert(err).Equal(nil)
		})
	})
}
