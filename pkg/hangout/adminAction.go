package hangout

import (
	"github.com/pkg/errors"
	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"strings"
)

func (ah *ActionHandler)handleAdmin(c command, m *Event)(string,error)  {
	if len(c.args) == 0{
		return ah.adminHelp(), nil
	}
	subCMD := c.args[0]
	switch subCMD {
	case "help":
		return ah.adminHelp(), nil
	case "add-users-to-team":
		if len(c.args) < 2{
			return "", &chat.MissingArgs{}
		}
		return ah.adminAddUsersToTeam(c.args[1], m.Space.Name)
	case "view-team":
		return ah.adminViewTeam(m.Space.Name)
		
	}

	return ah.adminHelp(), nil
}

func (ah *ActionHandler)adminHelp()string{
	return `available commands are:
`+"```@sprinbot admin add-user-to-team <user1>,<user2> ... ```"+`
`+"```@sprintbot admin make-user-admin <user>```"+`
`+"```@sprintbot admin view-team```"+`
`
}

func (ah *ActionHandler)adminAddUsersToTeam(user string, team string)(string, error) {

	return "", nil
}

func (ah *ActionHandler)adminViewTeam(team string)(string, error)  {
	t, err := ah.teamService.PopulateTeam(team)
	if err != nil{
		return "", errors.Wrap(err,"failed to get the team")
	}
	return `| team name : `+t.Name+` |
| Members   : `+strings.Join(t.Members,"\n")+``, nil
}
