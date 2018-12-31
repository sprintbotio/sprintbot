package chat

import "regexp"

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

var (
	AddToTeamRegexp      = regexp.MustCompile(`^add.*to (the\s)?team.*`)
	ViewTeamRegexp       = regexp.MustCompile(`^(view|show)?(\s\w+)?(\s\w+)?\steam.*`)
	RemoveFromTeamRegexp = regexp.MustCompile(`^remove.*from team.*`)
	MakeUsersAdmins      = regexp.MustCompile(`^make.*admin(s)?`)
	SetUserTimeZone      = regexp.MustCompile(`^set timezone.*\s(?P<zone>\w+\/\w+\/?\w+)`)
	ScheduleStandUp      = regexp.MustCompile(`^schedule standup\s?(for|at)?\s?(\d\d)?:?(\d\d)?\s?(\w+\/\w+)?`)
)
