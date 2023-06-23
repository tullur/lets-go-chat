package user

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestNewUser(t *testing.T) {
	g := Goblin(t)

	g.Describe("New User Creation", func() {
		user, err := New("Test User", "PasswordQwerty")
		if err != nil {
			g.Fail(err)
		}

		g.It("Should create user with valid UUID", func() {
			g.Assert(len(user.Id())).Equal(36)
		})

		g.It("Should create user with proper name", func() {
			g.Assert(user.Name()).Equal("Test User")
		})
	})
}

func TestUserPasswordVerification(t *testing.T) {
	g := Goblin(t)

	g.Describe("User Password Verification", func() {
		user, err := New("Test User", "PasswordQwerty")
		if err != nil {
			g.Fail(err)
		}

		g.Describe("When Password is correct", func() {
			result, err := user.VerifyPassword("PasswordQwerty")
			if err != nil {
				g.Fail(err)
			}

			g.It("Should be true", func() {
				g.Assert(result).Equal(true)
			})

			g.It("Should have empty error", func() {
				g.Assert(err).Equal(nil)
			})
		})

		g.Describe("When Password is incorrect", func() {
			result, err := user.VerifyPassword("NotCorrectPassword")

			g.It("Should be false", func() {
				g.Assert(result).Equal(false)
			})

			g.It("Should throw InvalidPassword error", func() {
				g.Assert(err).Equal(ErrInvalidPassword)
			})
		})
	})
}

func TestNewUserErrors(t *testing.T) {
	type testCase struct {
		test          string
		name          string
		password      string
		expectedError error
	}

	testCases := []testCase{
		{
			test:          "Name Length Validation",
			name:          "U",
			password:      "qwerty12345",
			expectedError: ErrShortUserName,
		},
		{
			test:          "Password Length Validation",
			name:          "UserTest",
			password:      "q",
			expectedError: ErrShortPasswordLength,
		},
		{
			test:          "Empty User Validation",
			name:          "",
			password:      "Password1234567",
			expectedError: ErrEmptyValues,
		},
		{
			test:          "Empty Password Validation",
			name:          "UserTest",
			password:      "",
			expectedError: ErrEmptyValues,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := New(tc.name, tc.password)
			if err != tc.expectedError {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}
