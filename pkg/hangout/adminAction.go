package hangout

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	cmdAddUsersToTeam = "add-users-to-team"
	cmdViewTeam = "view-team"
	cmdRemoveUserFromTeam = "remove-user-from-team"
)


func (ah *ActionHandler)handleAdmin(c command, m *Event)(string,error)  {
	// verify user is an admin before doing the command
	switch c.name {
	case "help":
		return ah.adminHelp(), nil
	case cmdAddUsersToTeam:
		return ah.adminAddUsersToTeam(c, m)
	case cmdViewTeam:
		return ah.adminViewTeam(m.Space.Name)
	case cmdRemoveUserFromTeam:
		return ah.adminRemoveUserFromTeam(c,m)
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

func (ah *ActionHandler)adminAddUsersToTeam(cmd command, m *Event)(string, error) {
	//parse out user names
	var added []string
	for _, a := range m.Message.Annotations{
		if a != nil && a.UserMention.User.DisplayName != "sprintbot"{
			//a.UserMention.User.
			logrus.Info("adding ", a.UserMention.User.DisplayName, " to the team")
			is, err := ah.teamService.IsUserInTeam(a.UserMention.User.DisplayName, m.Space.Name)
			if err != nil{
				return "", err
			}
			if is{
				continue
			}
			if err := ah.teamService.AddUserToTeam(a.UserMention.User.DisplayName,m.Space.Name); err != nil{
				return "", err
			}
			added = append(added, a.UserMention.User.DisplayName)
		}
	}

	return "added "+ strings.Join(added," : ") + " to the team " + m.Space.DisplayName, nil
}

func (ah *ActionHandler)adminRemoveUserFromTeam(cmd command, m *Event)(string,error)  {
	for _, a := range  m.Message.Annotations{
		if a != nil && a.UserMention.User.DisplayName != "sprintbot"{
			logrus.Info("removing user from team ", a.UserMention.User.DisplayName)
			if err := ah.teamService.RemoveUserFromTeam(a.UserMention.User.DisplayName, m.Space.Name); err != nil{
				return "",err
			}
		}
	}
	return `users removed`, nil
}

func (ah *ActionHandler)adminViewTeam(team string)(string, error)  {
	t, err := ah.teamService.PopulateTeam(team)
	if err != nil{
		return "", errors.Wrap(err,"failed to get the team")
	}
	return `| team name : `+t.Name+` |
| Members   : `+strings.Join(t.Members,"\n")+``, nil
}
