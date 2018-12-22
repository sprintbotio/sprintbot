package hangout

import (
	"github.com/sirupsen/logrus"
)

func (ah *ActionHandler )handleRegister(m *Event) (string,error) {
	logrus.Info("register ", m.Space, m.User.Name, m.User.DisplayName, m.User.Email)
	// check if this user is in a team as this could be the first direct message
	id, err := ah.userService.RegisterAdmin(m.User.DisplayName,m.User.Name,m.Space.Name)
	if err != nil{
		logrus.Error("failed to register admin ", err)
		return "I was unable to register you", err
	}
	if _, err := ah.teamService.RegisterTeam(m.Space.DisplayName, m.Space.Name, id); err != nil{
		return "I was unable to register the team", err
	}

	return `Thank you for registering SprintBot @`+m.User.DisplayName+`. 
You have been added as the admin for this space. Your team is named `+m.Space.DisplayName+`
You can add more admins using the `+"```sprintbot admin make-user-admin <user>```"+`
Admins can register teams and members to teams.
`+"```sprintbot admin help```"+` will give you more info`, nil
}
