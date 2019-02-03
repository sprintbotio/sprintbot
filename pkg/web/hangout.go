package web

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/gchat"
)

type HangoutHandler struct {
	chatActionHandler *chat.ActionHandler
	key               string
}

func NewHangoutHandler(chatActionHandler *chat.ActionHandler, key string) *HangoutHandler {
	return &HangoutHandler{chatActionHandler: chatActionHandler, key: key}
}

func (hh *HangoutHandler) Message(rw http.ResponseWriter, req *http.Request) {
	message := gchat.Event{}
	if err := json.NewDecoder(req.Body).Decode(&message); err != nil {
		logrus.Error("failed to decode hangout message ", err)
		http.Error(rw, "could not parse message", http.StatusBadRequest)
		return

	}
	if hh.key != message.Token {
		logrus.Error("incorrect token sent ", message.Token)
		http.Error(rw, "", http.StatusUnauthorized)
		return
	}
	response := hh.chatActionHandler.Handle(&message)

	if _, err := rw.Write([]byte(`{"text":"` + response + `"}`)); err != nil {
		logrus.Error("failed to write response ", err)
		http.Error(rw, "failed to handle chat action", http.StatusInternalServerError)
		return
	}
}
