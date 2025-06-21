package hasher_test

import (
	"log"
	"testing"

	"github.com/nurfianqodar/school-microservices/utils/hasher"
)

func TestHash(t *testing.T) {
	password := "secretpassword"
	validPassword := "secretpassword"
	invalidPassword := "invalidpassword"
	hash, err := hasher.GenerateFromPassword(password, hasher.DefaultConfig)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log(hash)

	if err := hasher.CompareHashWithPassword(hash, validPassword); err != nil {
		log.Println(err)
		t.Fail()
	}

	if err := hasher.CompareHashWithPassword(hash, invalidPassword); err == nil {
		t.Fail()
	}
}
