package web

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/sprintbot.io/sprintbot/pkg/web/middleware"
	"net/http"
)

func BuildRouter() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)
	return r
}

// BuildHTTPHandler puts together our HTTPHandler
func BuildHTTPHandler(r *mux.Router) http.Handler {
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n := negroni.New(recovery)
	n.Use(middleware.Cors{})
	n.UseHandler(r)
	return n
}

func MountSystemHandler(router *mux.Router) {
	syshandler := &SystemHandler{}
	router.HandleFunc("/api/sys/info/alive", syshandler.Alive)
}

func MountHangoutHandler(router *mux.Router, handler *HangoutHandler)  {
	router.HandleFunc("/api/hangout/message", handler.Message)
}