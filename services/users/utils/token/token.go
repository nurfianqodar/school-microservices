package token

import (
	"log"
	"os"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nurfianqodar/school-microservices/utils/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenType int8

const (
	TokenTypeAccess TokenType = iota
	TokenTypeRefresh
)

var _ jwt.Claims = (*Claims)(nil)

type Claims struct {
	Exp time.Time `json:"exp"`
	Iat time.Time `json:"iat"`
	Nbf time.Time `json:"nbf"`
	Iss string    `json:"iss"`
	Sub string    `json:"sub"`
	Aud []string  `json:"aud"`
	Typ TokenType `json:"typ"`
}

func (c *Claims) GetAudience() (jwt.ClaimStrings, error) {
	return c.Aud, nil
}

func (c *Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(c.Exp), nil
}

func (c *Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(c.Iat), nil
}

func (c *Claims) GetIssuer() (string, error) {
	return c.Iss, nil
}

func (c *Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(c.Nbf), nil
}

func (c *Claims) GetSubject() (string, error) {
	return c.Sub, nil
}

func (c *Claims) GetTokenType() TokenType {
	return c.Typ
}

func CreateToken(typ TokenType, sub string, expAfter time.Duration, aud []string) (string, error) {
	now := time.Now()
	exp := now.Add(expAfter)
	c := &Claims{
		Iss: "api.example.com",
		Iat: now,
		Nbf: now,
		Exp: exp,
		Sub: sub,
		Aud: aud,
		Typ: typ,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	appSecret, ok := os.LookupEnv("SECRET")
	if !ok {
		log.Println("error: unable to get SECRET environment variable")
		return "", errs.ErrInternalServer
	}

	tokenString, err := token.SignedString([]byte(appSecret))
	if err != nil {
		log.Printf("error: failed to sign token. %s\n", err.Error())
		return "", errs.ErrInternalServer
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*Claims, error) {
	var (
		c  *Claims = new(Claims)
		ok bool
	)

	token, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (any, error) {
		appSecret, ok := os.LookupEnv("SECRET")
		if !ok {
			log.Fatalln("error: unable to get SECRET environment variable")
			// Exit when app secret not set
		}
		return []byte(appSecret), nil
	})
	if err != nil {
		log.Printf("error: failed to parse with claims. %s\n", err.Error())
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	c, ok = token.Claims.(*Claims)
	if ok {
		return c, nil
	}

	log.Printf("error: incompatible token claims type. found %s\n", reflect.TypeOf(c).Name())
	return nil, status.Error(codes.Unauthenticated, "invalid token")
}
