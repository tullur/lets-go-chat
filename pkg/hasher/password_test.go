package hasher

import (
	"strings"
	"testing"
)

func TestHashPasswordFormat(t *testing.T) {
	var (
		password   = "qwerty"
		result, _  = HashPassword(password)
		hashValues = strings.Split(result, "$")
	)

	if len(hashValues) != 6 {
		t.Error("Hash is not in the correct format")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	var (
		password         = "password"
		hash, _          = HashPassword(password)
		correctCompare   = CheckPasswordHash(password, hash)
		incorrectCompare = CheckPasswordHash("dummyPassword", hash)
	)

	if !correctCompare {
		t.Error("Compare is not correct")
	}

	if incorrectCompare {
		t.Error("Compare is not correct")
	}
}
