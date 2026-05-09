package errcode

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
)

func TestGenHash(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	t.Errorf("HASH: %s", string(hash))
}
