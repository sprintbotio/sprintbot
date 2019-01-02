package hangout

import (
	"fmt"
	"strings"
	"time"
)

func (ah *ActionHandler) handleGeneral(cmd command, event *Event) (string, error) {
	switch cmd.name {
	case cmdStandUpLog:
		return ah.generalStandUpLog(cmd, event)
	}
	msg := ah.generalHelp()
	return msg, nil
}

func (ah *ActionHandler) generalHelp() string {
	return "general help"
}

func (ah *ActionHandler) generalStandUpLog(cmd command, event *Event) (string, error) {
	teamID := cmd.team.ID
	standUp, err := ah.standUpService.LoadStandUp(teamID, time.Now())
	if err != nil {
		return "", err
	}

	res := "Present: " + strings.Join(standUp.Present, " , ") + " \n"
	res += "Absent: " + strings.Join(standUp.Absent, " , ") + " \n"

	for _, l := range standUp.Log {
		res += fmt.Sprintf("@" + l.UserName + " logged: \n ```" + l.Message + "``` \n")
	}
	return res, nil
}
