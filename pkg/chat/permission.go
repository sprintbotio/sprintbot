package chat

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
	"github.com/sprintbot.io/sprintbot/pkg/team"
)

type Permissions struct {
	team *team.Service
}

func NewPermissions(ts *team.Service) *Permissions {
	return &Permissions{team: ts}
}

func (p *Permissions) CanUserDoCmd(u *domain.User, cmd Command) (bool, error) {
	logrus.Infof("can user do cmd %s %s ", cmd.ActionType, u.Role)
	if cmd.ActionType == "general" {
		return true, nil
	}
	isInTeam, err := p.team.IsUserInTeam(u.ID, cmd.TeamID)
	if err != nil {
		return false, errors.Wrap(err, "failed to check if user is in team")
	}
	if !isInTeam {
		return false, nil
	}

	if u.Role == "admin" {
		return true, nil
	}

	return u.Role == cmd.ActionType, nil
}
