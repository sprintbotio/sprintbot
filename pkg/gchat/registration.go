package gchat

import (
	"fmt"

	"github.com/sprintbot.io/sprintbot/pkg/chat"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/team"
	"github.com/sprintbot.io/sprintbot/pkg/user"
)

type Registration struct {
	userService *user.Service
	teamService *team.Service
}

const (
	registerResponse = `Thank you for registering SprintBot <%s>. 
You have been added as the admin for this space. Your team is named %s.
To find out more of what is available use the following command:
*@sprintbot admin help*`
)

func NewRegisterationUseCase(userService *user.Service, teamService *team.Service) *Registration {
	return &Registration{
		userService: userService,
		teamService: teamService,
	}
}

//TODO add metrics
func (r *Registration) HandleRegistration(cmd chat.Command, event *Event) (string, error) {
	logrus.Debug("register handler")
	id, err := r.userService.RegisterAdmin(event.User.DisplayName, event.User.Name, event.Space.Name)
	if err != nil {
		return "<" + event.User.Name + "> I was unable to register you or your team. Try removing me from the room and re adding me.", err
	}
	if _, err := r.teamService.RegisterTeam(event.Space.DisplayName, event.Space.Name, id); err != nil {
		return "I was unable to register the team", err
	}

	return fmt.Sprintf(registerResponse, event.User.Name, event.Space.DisplayName), nil
}

func (r *Registration) HandleUnRegister(cmd chat.Command, event *Event) (string, error) {
	// clean up the users and team and stand ups
	t, err := r.teamService.PopulateTeam(cmd.TeamID)
	if err != nil {
		return "", errors.Wrap(err, "failed to populate team to remove")
	}
	for _, u := range t.Users {
		if err := r.userService.DeleteUser(u.ID); err != nil {
			return "", err
		}
	}
	if err := r.teamService.RemoveTeam(cmd.TeamID); err != nil {
		return "", errors.Wrap(err, "failed to remove the team with id "+cmd.TeamID+" after being unregistered")
	}
	return "", nil

}

func (r *Registration) Register() {
	actionHandlers[chat.CommandRegister] = r.HandleRegistration
	actionHandlers[chat.CommandUnRegister] = r.HandleUnRegister
}
