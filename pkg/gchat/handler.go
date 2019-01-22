package gchat

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/sprintbot.io/sprintbot/pkg/domain"

	"github.com/sprintbot.io/sprintbot/pkg/user"

	"github.com/sprintbot.io/sprintbot/pkg/standup"

	"github.com/sprintbot.io/sprintbot/pkg/team"

	"github.com/sirupsen/logrus"

	"github.com/sprintbot.io/sprintbot/pkg/chat"
)

type ActionHandler struct {
	teamService    *team.Service
	standUpService *standup.Service
	userService    *user.Service
}

const (
	unexpectedERRRESPONSE = `sorry I was unable to %s. Please try again later`
)

type actionHandler func(chat.Command, *Event) (string, error)

// this is added to by each actions Register method
var actionHandlers = map[string]actionHandler{}

func NewActionHandler(teamService *team.Service, standupService *standup.Service, userService *user.Service) *ActionHandler {
	return &ActionHandler{teamService: teamService, standUpService: standupService, userService: userService}
}

func (ah *ActionHandler) Handle(m chat.Message) string {

	var commonErrResp = fmt.Sprintf(unexpectedERRRESPONSE, "figure out the action you wanted to take. ")
	//errors should be logged at this level
	event, ok := m.(*Event)
	if !ok {
		// really bad
		logrus.Error("Unexpected type passed. Should have been an *Event but got ", reflect.TypeOf(m))
		return commonErrResp
	}
	logrus.Info("handling message ", event.Message.ArgumentText, event.Space.Type, event.User, event.User.Name)
	clean := ah.cleanText(event.Message.ArgumentText)
	u, err := ah.userService.ResolveUser(event.User.Name, event.User.DisplayName)
	if err != nil {
		logrus.Error("failed to resolve user ", err)
		return commonErrResp
	}
	t, err := ah.teamService.ResolveTeamForUser(u)
	if err != nil {
		logrus.Error("failed to resolve team ", err)
		return commonErrResp
	}
	cmd, err := ah.parseCommand(t.ID, event, clean)
	//TODO need to have a think about how to handle free form messages
	if err != nil {
		if chat.IsUnkownCommandErr(err) && ah.standUpService.IsStandUpInProgress(cmd.TeamID) {
			if err := ah.standUpService.HandleStandUpMessage(cmd.TeamID, cmd.RequesterID, cmd.Requester, event.Message.ArgumentText); err != nil {
				logrus.Error("failed to handle stand up message", err)
				return "failed to log your stand up message"
			}
			return ""
		}
		logrus.Error("failed to parse a valid command", err)
		return commonErrResp
	}
	logrus.Info("cmd parsed ", cmd.Name, cmd.Args, cmd.Requester)
	if !chat.CanUserDoCmd(u, cmd) {
		logrus.Error("user ", u.Name, " cannot perform action ", cmd.Name)
		return "sorry but you cannot do that."
	}
	if h, ok := actionHandlers[cmd.Name]; ok {
		resp, err := h(cmd, event)
		if err != nil {
			logrus.Error("failed to execute command ", cmd.ActionType, err)
			switch err.(type) {
			case *domain.NotDirectMessageErr:
				return err.Error()
			default:
				return commonErrResp
			}

		}
		return resp
	}
	logrus.Info("no handler found for cmd")
	return commonErrResp
}

func (ah *ActionHandler) cleanText(argumentText string) string {
	r, _ := regexp.Compile(`/\s{2,}/g`)
	clean := strings.TrimSpace(argumentText)
	clean = r.ReplaceAllString(clean, " ")
	return clean
}

func (ah *ActionHandler) parseCommand(teamID string, event *Event, argumentText string) (chat.Command, error) {
	cmd := chat.Command{
		Requester:   event.User.DisplayName,
		RequesterID: event.User.Name,
		TeamID:      teamID,
	}
	if event.Type == "ADDED_TO_SPACE" {
		cmd.Name = chat.CommandRegister
		cmd.ActionType = "general"
		return cmd, nil

	}
	if event.Type == "REMOVED_FROM_SPACE" {
		cmd.Name = chat.CommandUnRegister
		cmd.ActionType = "general"
		return cmd, nil
	}
	logrus.Println("parse command text ", argumentText)
	for k, v := range chat.Commands {
		if k.MatchString(argumentText) {
			logrus.Println("found command ", k.String())
			v.TeamID = teamID
			v.Requester = event.User.DisplayName
			v.RequesterID = event.User.Name
			sm := k.FindStringSubmatch(argumentText)
			logrus.Info("sub matches ", sm, len(sm))
			v.Args = sm[1:]
			namedArgs := map[string]string{}
			for i, name := range k.SubexpNames() {
				if name != "" {
					namedArgs[name] = sm[i]
				}
			}
			v.MappedArgs = namedArgs
			return v, nil
		}
	}
	return cmd, chat.NewUknownCommand(argumentText)
}

func (ah *ActionHandler) Platform() string {
	return "hangout"
}
