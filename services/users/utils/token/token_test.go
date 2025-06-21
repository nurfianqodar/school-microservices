package token_test

import (
	"testing"
	"time"

	"github.com/nurfianqodar/school-microservices/services/users/utils/token"
)

func TestCreateToken(t *testing.T) {
	token, err := token.CreateToken(token.TokenTypeAccess, "dummy", time.Hour, []string{"dummy"})
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	t.Log(token)
}

func TestVerifyToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTA2LTIyVDAyOjAxOjQ3LjgyODI2MTAzNSswNzowMCIsImlhdCI6IjIwMjUtMDYtMjJUMDE6MDE6NDcuODI4MjYxMDM1KzA3OjAwIiwibmJmIjoiMjAyNS0wNi0yMlQwMTowMTo0Ny44MjgyNjEwMzUrMDc6MDAiLCJpc3MiOiJhcGkuZXhhbXBsZS5jb20iLCJzdWIiOiJkdW1teSIsImF1ZCI6WyJkdW1teSJdLCJ0eXAiOjB9.4TgVv1uUdqty6iMsHKTVmkjhCIazQKUkEh0UKBph4ZM"
	claims, err := token.VerifyToken(tokenString)
	if err != nil {
		t.Fail()
		t.Log(err)
	}
	t.Log(claims)
}
