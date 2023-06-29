package security

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

type Security struct {
	users          map[string]string
	expirationTime int
	jwtKey         []byte
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func New(users map[string]string, exTime int, jwtKey []byte) *Security {
	return &Security{users: users, expirationTime: exTime, jwtKey: jwtKey}
}

func (sec *Security) Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expectedPassword, ok := sec.users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	duration := time.Now().Add(time.Duration(sec.expirationTime) * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(duration),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, errToken := token.SignedString(sec.jwtKey)

	if errToken != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, tokenString)
}

func (sec *Security) VerifyToken(tokenStr string) (string, bool) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return sec.jwtKey, nil
	})
	if err != nil {
		log.Error().Msg(err.Error())
		return err.Error(), false
	}
	return "", token.Valid
}
