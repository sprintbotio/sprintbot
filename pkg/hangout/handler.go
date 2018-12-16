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

var (
	addToTeamRegexp = regexp.MustCompile(`^add.*to team.*`)
	viewTeamRegexp = regexp.MustCompile(`^view\steam.*`)
	removeFromTeamRegexp = regexp.MustCompile(`^remove.*from team.*`)
)

func NewActionHandler(ts *team.Service)*ActionHandler  {
	return &ActionHandler{teamService: ts}
}

func (ah *ActionHandler)Handle(m chat.Message) (string,error) {
	message := m.(*Event)
	logrus.Info("got new event ", message.Type, message.User, message.Message.Annotations)
	for _, a := range message.Message.Annotations{
		if a != nil{
			logrus.Info("mentions ", a.UserMention.User.DisplayName)
		}
	}
	if message.Type == "ADDED_TO_SPACE"{
		return ah.handleRegister(message)
	}else if message.Type == "MESSAGE"{
		logrus.Info("message arguments ", message.Message.ArgumentText)
		cleanCmd := ah.cleanText(message.Message.ArgumentText)
		cmd := ah.parseCommand(cleanCmd)

		switch cmd.actionType {
		case "admin":
			return ah.handleAdmin(cmd, message)
		default:
			return "",chat.NewUknownCommand(cmd.name)
		}

	}
	return "", chat.NewUknownCommand(message.Type)
}

func (ah *ActionHandler)cleanText(argumentText string)string{
	r, _ := regexp.Compile(`/\s{2,}/g`)
	clean := strings.TrimSpace(argumentText)
	clean = r.ReplaceAllString(clean," ")
	return clean
}

func (ah *ActionHandler)parseCommand(argumentText string)command  {
	cmd := command{argsText:argumentText}
	if addToTeamRegexp.MatchString(argumentText){
			cmd.name = cmdAddUsersToTeam
			cmd.actionType = "admin"
	}
	if viewTeamRegexp.MatchString(argumentText){
		cmd.name = cmdViewTeam
		cmd.actionType = "admin"
	}
	if removeFromTeamRegexp.MatchString(argumentText){
		cmd.name = cmdRemoveUserFromTeam
		cmd.actionType = "admin"
	}
	return cmd
}


func (ah *ActionHandler)Platform()string  {
	return "hangout"
}

