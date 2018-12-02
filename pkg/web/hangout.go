package web

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/hangout"
	"net/http"
)

type HangoutHandler struct {
	chatActionHandler *chat.ActionHandler

}

func NewHangoutHandler(chatActionHandler *chat.ActionHandler)*HangoutHandler{
	return &HangoutHandler{chatActionHandler:chatActionHandler}
}

func (hh *HangoutHandler)Message(rw http.ResponseWriter, req *http.Request)  {
	message := hangout.Event{}
	if err := json.NewDecoder(req.Body).Decode(&message); err != nil{
		logrus.Error("failed to decode hangout message ", err)
		http.Error(rw, "could not parse message", http.StatusBadRequest)
		return

	}
	response, err := hh.chatActionHandler.Handle(&message)
	if err != nil{
		logrus.Error("failed to handle chat action ", err)
		http.Error(rw,"failed to handle chat action", http.StatusInternalServerError)
		return
	}
	rw.Write([]byte(`{"text":"`+response+`"}`))
}
