package gchat

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sprintbot.io/sprintbot/pkg/domain"

	"github.com/sprintbot.io/sprintbot/pkg/user"

	"github.com/pkg/errors"
	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/standup"
)

type StandUpUseCase struct {
	standUpService *standup.Service
	userService    *user.Service
}

func NewStandUpUseCase(standUpService *standup.Service, us *user.Service) *StandUpUseCase {
	return &StandUpUseCase{standUpService: standUpService, userService: us}
}

func (ss *StandUpUseCase) ScheduleStandUp(cmd chat.Command, event *Event) (string, error) {
	if len(cmd.Args) < 4 && cmd.NoEmptyArgs() {
		return `It looks like you are trying to schedule a stand up?
please resend with the following format ` +
			"```schedule standup at HH:SS <TimeZone (E.G Europe/Dublin)>.```" + `
`, nil
	}
	hour, err := strconv.ParseInt(cmd.Args[1], 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse hour for meeting")
	}
	min, err := strconv.ParseInt(cmd.Args[2], 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse min for meeting")
	}
	if err := ss.standUpService.SaveTime(cmd.TeamID, hour, min, cmd.Args[3]); err != nil {
		return "", err
	}
	return "Stand Up Scheduled", nil
}

func (ss *StandUpUseCase) StandUpLog(cmd chat.Command, event *Event) (string, error) {
	standUp, err := ss.standUpService.LoadStandUp(cmd.TeamID, time.Now())
	if err != nil {
		if domain.IsNotFoundErr(err) {
			return "There is no log for a stand up today. Perhaps it has not run yet?", nil
		}
		return "", err
	}

	res := "Present: " + strings.Join(standUp.Present, " , ") + " \n"
	res += "Absent: " + strings.Join(standUp.Absent, " , ") + " \n"

	for _, l := range standUp.Log {
		res += fmt.Sprintf("@" + l.UserName + " logged: \n ```" + l.Message + "``` \n")
	}
	return res, nil
}

func (ss *StandUpUseCase) RemoveLastStandUp(cmd chat.Command, event *Event) (string, error) {
	if err := ss.standUpService.RemoveMostRecentStandUp(cmd.TeamID); err != nil {
		return "", err
	}
	return "stand up removed", nil
}

func (ss *StandUpUseCase) LogStandUpStatus(cmd chat.Command, event *Event) (string, error) {
	// NEED TO FIGURE OUT WHICH STAND UP THIS IS FOR
	// if  a standup exists and is complete then this is for the next day.
	// if a standup exists but is not yet started then someone else has logged a status and we need to add to it
	// If no standup present then it is for the current day in the users timezone
	if event.Space.Type != "DM" {
		return "please use a direct message for this.", nil
	} // needs to be done as direct message
	successFmt := "Your status was logged. It will be shown to the team during the stand up for the date  %d-%02d-%02d "
	log := &domain.StandUpLog{
		Message:  strings.Join(cmd.Args, " "),
		UserName: cmd.Requester,
		UserID:   cmd.RequesterID,
	}
	u, err := ss.userService.GetUser(cmd.RequesterID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get user when loggin stand up status")
	}
	if u.Team != cmd.TeamID {
		return "", errors.New("team members can only log statuses for team they are a part of")
	}
	l, err := time.LoadLocation(u.Timezone)
	if err != nil {
		return "", errors.Wrap(err, "failed to load timezone for user")
	}
	t := time.Now().In(l)
	s, err := ss.standUpService.LoadStandUp(cmd.TeamID, t)
	if err != nil && !domain.IsNotFoundErr(errors.Cause(err)) {
		return "", errors.Wrap(err, "unexpected error loading stand up")
	}
	if s == nil {
		// no standup
		s, err := ss.standUpService.CreateStandUp(t, cmd.TeamID)
		if err != nil {
			return "", err
		}
		if err := ss.standUpService.AddStandUpLog(s.ID, log); err != nil {
			return "", err
		}
		return fmt.Sprintf(successFmt, t.Year(), t.Month(), t.Day()), nil

	}
	if s.StartTime == 0 && s.EndTime == 0 {
		// not run yet but another user has logged a status
		if err := ss.standUpService.AddStandUpLog(s.ID, log); err != nil {
			return "", err
		}
		return fmt.Sprintf(successFmt, t.Year(), t.Month(), t.Day()), nil
	} else if s.StartTime != 0 && s.EndTime == 0 {
		// standup is running promt the user to join the stand up
		return "It looks like your team's stand up is currently in progress. Join the main room and add your status there. Or wait until the stand up is finished", nil

	}
	// stand up is complete for this day so we will log for the next day
	tomorrow := t.Add(time.Hour + 24)
	ts, err := ss.standUpService.CreateStandUp(tomorrow, cmd.TeamID)
	if err != nil {

	}
	if err := ss.standUpService.AddStandUpLog(ts.ID, log); err != nil {
		return "", err
	}
	return fmt.Sprintf(successFmt, t.Year(), t.Month(), t.Day()), nil

}

func (ss *StandUpUseCase) Register() {
	actionHandlers[chat.CommandStandUpLog] = ss.StandUpLog
	actionHandlers[chat.CommandScheduleStandUp] = ss.ScheduleStandUp
	actionHandlers[chat.CommandStandUpStatus] = ss.LogStandUpStatus
	actionHandlers[chat.CommandRemoveLatestStandUp] = ss.RemoveLastStandUp
}
