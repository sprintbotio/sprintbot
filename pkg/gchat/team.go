package gchat

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/team"
)

type Team struct {
	teamService *team.Service
}

func NewTeamUseCase(ts *team.Service) *Team {
	return &Team{teamService: ts}
}

const (
	addedUserResp = "Added %s to team: %s \n" +
		"update this users timzone using ```set timezone @user Europe/Dublin ``` \n" +
		"For a full list of timezones visit https://timezonedb.com/time-zones"
)

func (t *Team) AddUserToTeam(cmd chat.Command, event *Event) (string, error) {
	var added []string
	for _, a := range event.Message.Annotations {
		if a != nil && (a.UserMention.User.DisplayName != "sprintbot" && a.UserMention.User.DisplayName != "") {
			userToAdd := a.UserMention.User.DisplayName
			userID := a.UserMention.User.Name
			is, err := t.teamService.IsUserInTeam(userToAdd, cmd.TeamID)
			if err != nil {
				return "", err
			}
			if is {
				continue
			}
			if err := t.teamService.AddUserToTeam(userToAdd, userID, cmd.TeamID, "member"); err != nil {
				return "", err
			}
			added = append(added, userToAdd)
		}
	}

	return fmt.Sprintf(addedUserResp, strings.Join(added, ","), event.Space.DisplayName), nil
}

func (t *Team) RemoveUserFromTeam(cmd chat.Command, event *Event) (string, error) {
	var removed []string
	for _, a := range event.Message.Annotations {
		if a != nil && a.UserMention.User.DisplayName != "sprintbot" {
			logrus.Info("removing user from team ", a.UserMention.User.DisplayName)
			if err := t.teamService.RemoveUserFromTeam(a.UserMention.User.DisplayName, cmd.TeamID); err != nil {
				return "", err
			}
			removed = append(removed, a.UserMention.User.DisplayName)
		}
	}
	res := `Removed ` + strings.Join(removed, ",") + ` from team ` + event.Space.DisplayName + ``
	if len(removed) == 0 {
		res = "no users removed from the team"
	}
	return res, nil
}

func (t *Team) ViewTeam(cmd chat.Command, event *Event) (string, error) {
	team, err := t.teamService.PopulateTeam(cmd.TeamID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get the team")
	}
	var teamView = `|` + team.Name + `|\n`
	for _, u := range team.Users {
		teamView += u.Name + ` ` + u.Timezone + `\n`
	}
	return teamView, nil
}

func (t *Team) Register() {
	actionHandlers[chat.CommandAddUserToTeam] = t.AddUserToTeam
	actionHandlers[chat.CommandRemoveUserFromTeam] = t.RemoveUserFromTeam
	actionHandlers[chat.CommandViewTeam] = t.ViewTeam
}
