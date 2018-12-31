package hangout

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/team"
	"github.com/sprintbot.io/sprintbot/pkg/user"
)

type ActionHandler struct {
	teamService    *team.Service
	userService    *user.Service
	standUpService *team.StandUpService
}

const (
	unexpectedERRRESPONSE = `sorry I was unable to %s. Please try again later`
)

func NewActionHandler(ts *team.Service, us *user.Service, ss *team.StandUpService) *ActionHandler {
	return &ActionHandler{teamService: ts, userService: us, standUpService: ss}
}

func (ah *ActionHandler) Handle(m chat.Message) string {
	message, ok := m.(*Event)
	if !ok {
		logrus.Error("the chat message passed to the google chat handler was not an Event it was ", reflect.TypeOf(m))
		return fmt.Sprintf(unexpectedERRRESPONSE, "complete your requested action")
	}
	// look for user and set team
	t, err := ah.teamService.GetTeamForUser(message.User.Name)
	if err != nil {
		logrus.Errorf("failed when trying to find team user a member of ", err, reflect.TypeOf(err))
		//return fmt.Sprintf(unexpectedERRRESPONSE, "complete the action")
	}
	logrus.Info("got new event ", message.Type, "space: ", message.Space, message.User, message.Message.Annotations, message.Message.Thread.Name)
	if message.Type == "ADDED_TO_SPACE" && t == nil {
		resp, err := ah.handleRegister(message)
		if err != nil {
			logrus.Errorf("failed to register new bot for gChat %+v", err)
			return fmt.Sprintf(unexpectedERRRESPONSE, "register your bot")
		}
		return resp
	}
	if message.Type == "REMOVED_FROM_SPACE" && t != nil {
		if err := ah.teamService.RemoveTeam(t.ID); err != nil {
			logrus.Errorf("failed to remove team %s %+v", t.ID, err)
			return ""
		}
		return ""
	}
	if message.Type == "MESSAGE" {
		cleanCmd := ah.cleanText(message.Message.ArgumentText)
		cmd, err := ah.parseCommand(cleanCmd)

		if err != nil {
			logrus.Info("cmd err ", err)
			if _, ok := err.(*chat.UnkownCommand); ok && ah.standUpService.IsStandUpInProgress(t.ID) {
				logrus.Info("stand up is in progress so logging standup message ", message.Message.Text)
				// this is a but of hack but when a standup is in progress anything could be sent to the bot so any unknown command errors are treated as standup messages
				if err := ah.standUpService.LogStandUpMessage(t.ID, message.User.Name, message.Message.ArgumentText); err != nil {
					logrus.Errorf("failed to log stand up message ", err)
					return "Unable to log the stand up status"
				}
				return ""
			}

			logrus.Errorf("error parsing command in google chat message %+v", err)
			return fmt.Sprintf(unexpectedERRRESPONSE, "complete your requested action")

		}
		if t != nil {
			cmd.team = t
		}
		switch cmd.actionType {
		case "admin":
			resp, err := ah.handleAdmin(cmd, message)
			if err != nil {
				logrus.Errorf("failed to handle admin action %s %+v", cmd.name, err)
				return fmt.Sprintf(unexpectedERRRESPONSE, "complete your "+cmd.name+" action")
			}
			return resp
		}
	}
	return fmt.Sprintf("sorry I do not understand. Try @sprintbot help")
}

func (ah *ActionHandler) cleanText(argumentText string) string {
	r, _ := regexp.Compile(`/\s{2,}/g`)
	clean := strings.TrimSpace(argumentText)
	clean = r.ReplaceAllString(clean, " ")
	return clean
}

func (ah *ActionHandler) parseCommand(argumentText string) (command, error) {
	cmd := command{}
	if chat.AddToTeamRegexp.MatchString(argumentText) {
		cmd.name = cmdAddUsersToTeam
		cmd.actionType = "admin"
	}
	if chat.ViewTeamRegexp.MatchString(argumentText) {
		cmd.name = cmdViewTeam
		cmd.actionType = "admin"
	}
	if chat.RemoveFromTeamRegexp.MatchString(argumentText) {
		cmd.name = cmdRemoveUserFromTeam
		cmd.actionType = "admin"
	}
	if chat.MakeUsersAdmins.MatchString(argumentText) {
		cmd.name = cmdMakeUserAdmin
		cmd.actionType = "admin"
	}
	if chat.SetUserTimeZone.MatchString(argumentText) {
		cmd.name = cmdSetUserTimeZone
		cmd.actionType = "admin"
		m := chat.SetUserTimeZone.FindStringSubmatch(argumentText)
		cmd.args = []string{strings.TrimSpace(m[1])}
	}
	if chat.ScheduleStandUp.MatchString(argumentText) {
		cmd.name = cmdScheduleStandUp
		cmd.actionType = "admin"
		m := chat.ScheduleStandUp.FindStringSubmatch(argumentText)
		fmt.Println("substr match ", m[1:], m, len(m))

		cmd.args = m[1:]
	}
	if cmd.name == "" {
		return cmd, chat.NewUknownCommand(argumentText)
	}
	return cmd, nil
}

func (ah *ActionHandler) Platform() string {
	return "hangout"
}
