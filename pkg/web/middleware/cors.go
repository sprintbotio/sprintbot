package middleware

import (
	"net/http"
	"os"
)

type Cors struct{}

func (c Cors) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	origin := "https://storypoint.me"
	if os.Getenv("APP_ENVIRONMENT") == "development" {
		origin = "http://localhost:3000"
	}
	headers := w.Header()
	if headers.Get("Access-Control-Allow-Origin") == "" {
		headers.Add("Access-Control-Allow-Origin", origin)
	}
	if headers.Get("Access-Control-Allow-Credentials") == "" {
		headers.Add("Access-Control-Allow-Credentials", "true")
	}
	if headers.Get("Access-Control-Allow-Methods") == "" {
		headers.Add("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
	}
	if headers.Get("Access-Control-Allow-Headers") == "" {
		headers.Add("Access-Control-Allow-Headers", "content-type")
	}
	if req.Method != "OPTIONS" {
		next(w, req)
		return
	}
}
