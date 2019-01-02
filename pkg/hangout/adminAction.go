package hangout

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	cmdAddUsersToTeam     = "add-users-to-team"
	cmdViewTeam           = "view-team"
	cmdRemoveUserFromTeam = "remove-user-from-team"
	cmdMakeUserAdmin      = "make-user-admin"
	cmdUnMakeUserAdmin    = "remove-admin"
	cmdSetUserTimeZone    = "set-user-timezone"
	cmdScheduleStandUp    = "schedule-standup"
	cmdStandUpLog         = "standup-log"
)

func userUserIDAndTeamFromEvent(m *Event) (string, string, string) {
	return m.User.DisplayName, m.User.Name, m.Space.Name
}

func (ah *ActionHandler) handleAdmin(c command, m *Event) (string, error) {
	// verify user is an admin before doing the command
	_, uid, _ := userUserIDAndTeamFromEvent(m)
	is, err := ah.teamService.IsUserAdminForTeam(uid, c.team.ID)
	if err != nil {
		return "", err
	}
	if !is {
		return `sorry you cannot do that`, nil
	}
	switch c.name {
	case "help":
		return ah.adminHelp(), nil
	case cmdAddUsersToTeam:
		return ah.adminAddUsersToTeam(c, m)
	case cmdViewTeam:
		return ah.adminViewTeam(c.team.ID)
	case cmdRemoveUserFromTeam:
		return ah.adminRemoveUserFromTeam(c, m)
	case cmdSetUserTimeZone:
		return ah.adminSetUserTimezone(c, m)
	case cmdScheduleStandUp:
		return ah.adminScheduleStandUp(c, m)
	}

	return ah.adminHelp(), nil
}

func (ah *ActionHandler) adminHelp() string {
	return `available- commands are:
` + "```@sprinbot add <user1>,<user2> to team ... ```" + `
` + "```@sprintbot make <user1>,<user2> admins```" + `
` + "```@sprintbot view team```" + `
`
}

func (ah *ActionHandler) adminSetUserTimezone(cmd command, m *Event) (string, error) {
	//ah.userService.UpdateTZ()
	updated := []string{}
	for _, a := range m.Message.Annotations {
		if a != nil && (a.UserMention.User.DisplayName != "sprintbot" && a.UserMention.User.DisplayName != "") {
			userToAdd := a.UserMention.User.DisplayName
			userID := a.UserMention.User.Name
			if err := ah.userService.UpdateTZ(userID, cmd.args[0]); err != nil {
				return "", err
			}
			updated = append(updated, userToAdd)
		}
	}
	return `updated timezone to ` + cmd.args[0] + ` for ` + strings.Join(updated, " , "), nil
}

func (ah *ActionHandler) adminAddUsersToTeam(cmd command, m *Event) (string, error) {
	//parse out user names
	_, _, teamName := userUserIDAndTeamFromEvent(m)
	var added []string
	for _, a := range m.Message.Annotations {
		if a != nil && (a.UserMention.User.DisplayName != "sprintbot" && a.UserMention.User.DisplayName != "") {
			userToAdd := a.UserMention.User.DisplayName
			userID := a.UserMention.User.Name
			logrus.Info("adding ", userToAdd, " to the team")
			is, err := ah.teamService.IsUserInTeam(userToAdd, teamName)
			if err != nil {
				return "", err
			}
			if is {
				continue
			}
			if err := ah.teamService.AddUserToTeam(userToAdd, userID, m.Space.Name); err != nil {
				return "", err
			}
			added = append(added, userToAdd)
		}
	}

	return "Added " + strings.Join(added, " : ") + " to team: " + m.Space.DisplayName + " \n" +
		"update this users timzone using ```set timezone @user Europe/Dublin ``` \n" +
		"For a full list of timezones visit https://timezonedb.com/time-zones", nil
}

func (ah *ActionHandler) adminRemoveUserFromTeam(cmd command, m *Event) (string, error) {
	var removed []string
	for _, a := range m.Message.Annotations {
		if a != nil && a.UserMention.User.DisplayName != "sprintbot" {
			logrus.Info("removing user from team ", a.UserMention.User.DisplayName)
			if err := ah.teamService.RemoveUserFromTeam(a.UserMention.User.DisplayName, m.Space.Name); err != nil {
				return "", err
			}
			removed = append(removed, a.UserMention.User.DisplayName)
		}
	}
	res := `Removed ` + strings.Join(removed, ",") + ` from team ` + m.Space.DisplayName + ``
	if len(removed) == 0 {
		res = "no users removed from the team"
	}
	return res, nil
}

func (ah *ActionHandler) adminMakeUserAdmin(cmd command, m *Event) (string, error) {
	return "", nil
}

func (ah *ActionHandler) adminViewTeam(team string) (string, error) {
	logrus.Info("view team ", team)
	t, err := ah.teamService.PopulateTeam(team)
	if err != nil {
		return "", errors.Wrap(err, "failed to get the team")
	}
	var teamView = `|` + t.Name + `|\n`
	for _, u := range t.Users {
		teamView += u.Name + ` ` + u.Timezone.Name + `\n`
	}
	return teamView, nil
}

func (ah *ActionHandler) adminScheduleStandUp(cmd command, m *Event) (string, error) {
	logrus.Info("cmd args ", strings.Join(cmd.args, ","), len(cmd.args))
	if len(cmd.args) < 4 && cmd.NoEmptyArgs() {
		return `It looks like you are trying to schedule a stand up?
please resend with the following format ` +
			"```schedule standup at HH:SS <TimeZone (E.G Europe/Dublin)>.```" + `
`, nil
	}
	hour, err := strconv.ParseInt(cmd.args[1], 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse hour for meeting")
	}
	min, err := strconv.ParseInt(cmd.args[2], 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse min for meeting")
	}
	if err := ah.standUpService.SaveTime(cmd.team.ID, hour, min, cmd.args[3]); err != nil {
		return "", err
	}
	return "Stand Up Scheduled", nil
}
