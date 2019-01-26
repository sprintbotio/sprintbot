package chat

import (
	"regexp"

	"github.com/sirupsen/logrus"

	"github.com/sprintbot.io/sprintbot/pkg/domain"
)

type ActionHandler struct {
	handlers map[string]Handler
}

func (ah *ActionHandler) RegisterHandler(h Handler) {
	ah.handlers[h.Platform()] = h
}

func NewActionHandler() *ActionHandler {
	return &ActionHandler{handlers: map[string]Handler{}}
}

func (ah *ActionHandler) Handle(m Message) string {
	switch m.Platform() {
	case "hangout":
		return ah.handlers["hangout"].Handle(m)

	}
	return ""
}

type Handler interface {
	Handle(m Message) string
	Platform() string
}

type Message interface {
	Platform() string
}

type Command struct {
	ActionType  string
	TeamID      string
	Requester   string
	RequesterID string
	Name        string
	Args        []string
	MappedArgs  map[string]string
	Room        string
}

func (c Command) NoEmptyArgs() bool {
	for _, a := range c.Args {
		if a == "" {
			return false
		}
	}
	return true
}

const (
	CommandRegister            = "register"
	CommandUnRegister          = "unregister"
	CommandAddUserToTeam       = "add user to team"
	CommandViewTeam            = "view team"
	CommandRemoveUserFromTeam  = "remove user from team"
	CommandMakeUserAdmin       = "add admin"
	CommandSetUserTZ           = "set user tz"
	CommandScheduleStandUp     = "schedule standUP"
	CommandStandUpLog          = "standUp log"
	CommandStandUpStatus       = "log status"
	CommandAdminHelp           = "admin help"
	CommandGeneralHelp         = "general help"
	CommandRemoveLatestStandUp = "remove stand up"
)

var (
	AddToTeamRegexp           = regexp.MustCompile(`^add.*to (the\s)?team\s?(in\stimezone\s)?(set\stimezone\s)?(?P<tz>\w+\/\w+\/?\w+)?`)
	ViewTeamRegexp            = regexp.MustCompile(`^(view|show)?(\s\w+)?(\s\w+)?\steam.*`)
	RemoveFromTeamRegexp      = regexp.MustCompile(`^remove.*from team.*`)
	MakeUsersAdminsRegexp     = regexp.MustCompile(`^make.*admin(s)?`)
	SetUserTimeZoneRegexp     = regexp.MustCompile(`^set (?P<self>my)?\s?timezone.*\s(?P<zone>\w+\/\w+\/?\w+)`)
	ScheduleStandUpRegexp     = regexp.MustCompile(`^schedule standup\s?(for|at)?\s?(\d\d)?:?(\d\d)?\s?(\w+\/\w+)?`)
	StandUpLogRegexp          = regexp.MustCompile(`^(view)?\s?standup log$`)
	AdminHelpRegexp           = regexp.MustCompile(`^admin help$`)
	GeneralHelpRegexp         = regexp.MustCompile(`^help$`)
	LogStandUpStatusRegexp    = regexp.MustCompile("^log standup status (.*)")
	RemoveLatestStandUpRegexp = regexp.MustCompile(`^remove latest standup$`)
	Commands                  = map[*regexp.Regexp]Command{
		GeneralHelpRegexp:         {ActionType: "general", Name: CommandGeneralHelp},
		AdminHelpRegexp:           {ActionType: "admin", Name: CommandAdminHelp},
		AddToTeamRegexp:           {ActionType: "admin", Name: CommandAddUserToTeam},
		ViewTeamRegexp:            {ActionType: "admin", Name: CommandViewTeam},
		RemoveFromTeamRegexp:      {ActionType: "admin", Name: CommandRemoveUserFromTeam},
		MakeUsersAdminsRegexp:     {ActionType: "admin", Name: CommandMakeUserAdmin},
		SetUserTimeZoneRegexp:     {ActionType: "admin", Name: CommandSetUserTZ},
		ScheduleStandUpRegexp:     {ActionType: "admin", Name: CommandScheduleStandUp},
		StandUpLogRegexp:          {ActionType: "general", Name: CommandStandUpLog},
		LogStandUpStatusRegexp:    {ActionType: "member", Name: CommandStandUpStatus},
		RemoveLatestStandUpRegexp: {ActionType: "admin", Name: CommandRemoveLatestStandUp},
	}
)

func CanUserDoCmd(u *domain.User, cmd Command) bool {
	logrus.Infof("can user do cmd %s %s ", cmd.ActionType, u.Role)
	if cmd.ActionType == "general" {
		return true
	}
	if u.Role == "admin" {
		return true
	}

	return u.Role == cmd.ActionType
}
