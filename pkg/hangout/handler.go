package hangout

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/team"
	"reflect"
	"regexp"
	"strings"
)

type ActionHandler struct {
  teamService *team.Service
}


const(
	unexpectedERRRESPONSE = `sorry I was unable to %s. Please try again later`
)

var (
	addToTeamRegexp = regexp.MustCompile(`^add.*to (the\s)?team.*`)
	viewTeamRegexp = regexp.MustCompile(`^(view|show)?(\s\w+)?(\s\w+)?\steam.*`)
	removeFromTeamRegexp = regexp.MustCompile(`^remove.*from team.*`)
	makeUsersAdmins = regexp.MustCompile(`^make.*admin(s)?`)
)

func NewActionHandler(ts *team.Service)*ActionHandler  {
	return &ActionHandler{teamService: ts}
}

func (ah *ActionHandler)Handle(m chat.Message) string {
	message, ok := m.(*Event)
	if ! ok{
		logrus.Error("the chat message passed to the google chat handler was not an Event it was ", reflect.TypeOf(m))
		return fmt.Sprintf(unexpectedERRRESPONSE, "complete your requested action")
	}
	// look for user and set team
	t,err := ah.teamService.GetTeamForUser(message.User.Name)
	if err != nil{
		logrus.Errorf("failed when trying to find team user a member of ",err, reflect.TypeOf(err))
		//return fmt.Sprintf(unexpectedERRRESPONSE, "complete the action")
	}
	logrus.Info("got new event ", message.Type, message.User, message.Message.Annotations)
	if message.Type == "ADDED_TO_SPACE" && t == nil{
		resp, err :=  ah.handleRegister(message)
		if err != nil{
			logrus.Errorf("failed to register new bot for gChat %+v",err)
			return fmt.Sprintf(unexpectedERRRESPONSE, "register your bot")
		}
		return resp
	}
	if message.Type == "REMOVED_FROM_SPACE" && t != nil{
		if err := ah.teamService.RemoveTeam(t.ID); err != nil{
			logrus.Errorf("failed to remove team %s %+v",t.ID, err)
			return ""
		}
		return ""
	}
	if message.Type == "MESSAGE"{
		cleanCmd := ah.cleanText(message.Message.ArgumentText)
		cmd, err := ah.parseCommand(cleanCmd)
		if err != nil{
			logrus.Errorf("error parsing command in google chat message %+v", err)
			return fmt.Sprintf(unexpectedERRRESPONSE, "complete your requested action")
		}
		if t != nil{
			cmd.team =t
		}
		switch cmd.actionType {
		case "admin":
			resp, err := ah.handleAdmin(cmd, message)
			if err != nil{
				logrus.Errorf("failed to handle admin action %s %+v",cmd.name, err)
				return fmt.Sprintf(unexpectedERRRESPONSE, "complete your "+cmd.name+" action")
			}
			return resp
		}
	}
	return fmt.Sprintf("sorry I do not understand. Try @sprintbot help")
}

func (ah *ActionHandler)cleanText(argumentText string)string{
	r, _ := regexp.Compile(`/\s{2,}/g`)
	clean := strings.TrimSpace(argumentText)
	clean = r.ReplaceAllString(clean," ")
	return clean
}

func (ah *ActionHandler)parseCommand(argumentText string)(command,error)  {
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
	if makeUsersAdmins.MatchString(argumentText){
		cmd.name = cmdMakeUserAdmin
		cmd.actionType = "admin"
	}
	if cmd.name == ""{
		return cmd, chat.NewUknownCommand(argumentText)
	}
	return cmd, nil
}


func (ah *ActionHandler)Platform()string  {
	return "hangout"
}

