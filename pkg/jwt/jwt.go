package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"persons/internal/domain/models"
	"time"
)

const sold = "ija0rjvvnaoivn"

func NewToken(person models.Person, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = person.Id
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(sold))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
