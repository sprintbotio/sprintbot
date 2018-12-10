package hangout

import (
	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/team"
	"regexp"
	"strings"
)

type ActionHandler struct {

  teamService *team.Service

}

func NewActionHandler(ts *team.Service)*ActionHandler  {
	return &ActionHandler{teamService: ts}
}

func (ah *ActionHandler)Handle(m chat.Message) (string,error) {
	message := m.(*Event)
	logrus.Info("got new event ", message.Type, message.Message, message.User)
	if message.Type == "ADDED_TO_SPACE"{
		return ah.handleRegister(message)
	}else if message.Type == "MESSAGE"{
		cmd := ah.parseCommand(message.Message.ArgumentText)
		switch cmd.name {
		case "admin":
			return ah.handleAdmin(cmd, message)

		default:
			return "",chat.NewUknownCommand(cmd.name)
		}
	}
	logrus.Info("message arguments ", message.Message.ArgumentText)
	return "", chat.NewUknownCommand(message.Type)
}

func (ah *ActionHandler)parseCommand(argumentText string)command  {
	r, _ := regexp.Compile(`/\s{2,}/g`)
	valid := strings.TrimSpace(r.ReplaceAllString(argumentText," "))
	cmdParts := strings.Split(valid, " ")
	return command{name:cmdParts[0], args:cmdParts[1:]}
}

func (ah *ActionHandler)Platform()string  {
	return "hangout"
}

type command struct{
	name string
	args []string
}