package hangout

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/team"
)

type ActionHandler struct {

  teamService *team.Service

}

func NewActionHandler(ts *team.Service)*ActionHandler  {
	return &ActionHandler{teamService: ts}
}

func (ah *ActionHandler)Handle(m chat.Message) (string,error) {
	message := m.(*Event)
	logrus.Info("got new event ", message.Type, message.Message)
	if message.Type == "ADDED_TO_SPACE"{

		return ah.handleRegister(message)
	}
	return "", fmt.Errorf("unknown message type " + message.Type)
}

func (ah *ActionHandler)Platform()string  {
	return "hangout"
}