package web

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type SystemHandler struct {

}

func (sh *SystemHandler)Alive(rw http.ResponseWriter, req *http.Request)  {
	rw.Header().Add("content-type","application/json")
	if _, err := rw.Write([]byte(`{"http":"ok"}`)); err != nil{
		logrus.Error("failed to write response", err)
		http.Error(rw, "failed to write response", http.StatusInternalServerError)
		return
	}
}
