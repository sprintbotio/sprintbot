package middleware

import (
	"net/http"
)

type AuthToken struct {

}

func NewAuthToken() *AuthToken {
	return &AuthToken{}
}


func (at *AuthToken) TokenHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		next( writer, request)
	}
}
