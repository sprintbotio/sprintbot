package hangout

import (
	"fmt"
	"strings"
	"time"

	"github.com/sprintbot.io/sprintbot/pkg/domain"

	"github.com/sprintbot.io/sprintbot/pkg/team"

	"github.com/sirupsen/logrus"
)

type StandUpRunner struct {
	chat        *Service
	teamService *team.Service
}

var standupThread = map[string]string{}
var userDone = make(chan struct{})
var userTimer *time.Timer

func NewStandUpRunnder(chat *Service, ts *team.Service) *StandUpRunner {
	return &StandUpRunner{chat: chat, teamService: ts}
}

func (sr *StandUpRunner) InProgress(teamID string) bool {
	_, ok := standupThread[teamID]
	return ok
}

func (sr *StandUpRunner) Run(teamID string, msgChan chan domain.StandUpMsg) {
	msgBuilder := NewMessageBuilder()
	var startMsg = "Stand up is now starting. Each team member will be asked for their status update. If a team member doesn't respond within 2 mins they will be considered absent." +
		"Due to the limitations of hangouts, *Please ensure to mention @sprintbot in your replies*\n"
	threadID := standupThread[teamID]
	msg := msgBuilder.Text(startMsg).Thread(threadID).Build()
	if _, err := sr.chat.SendMessageToTeam(teamID, msg); err != nil {
		logrus.Error("failed to start standup ", err)
		return
	}
	t, err := sr.teamService.PopulateTeam(teamID)
	if err != nil {
		logrus.Errorf("failed to get a populated team", err)
		return
	}
	msgFormat := "<%s> please give your update. Remember to *mention @sprintbot in your reply*"

	go func() {
		// read any status messages
		for m := range msgChan {
			fmt.Println("got message for stand up ", teamID, m)
			// add a stand up log reset the timer for the user (in case they want to add more detail)
			if strings.TrimSpace(strings.ToLower(m.Msg)) == "complete" {
				logrus.Info("got a complete message")
				msgBuilder.Text("thanks for your status <" + m.UserID + ">")
				if _, err := sr.chat.SendMessageToTeam(teamID, msgBuilder.Build()); err != nil {
					logrus.Errorf("failed to send message to team ", err)
				}
				logrus.Info("sending user done")
				userTimer.Reset(time.Second)
				logrus.Info("reset timer ")
				continue
			}
			msgBuilder.Text("Logged status. Are you finished? Use: *@sprintbot complete* to end your status. I will naturally move on in 1 minute")
			if _, err := sr.chat.SendMessageToTeam(teamID, msgBuilder.Build()); err != nil {
				logrus.Errorf("failed to send message to team ", err)
			}
			logrus.Info("resetting the timer")
			userTimer.Reset(time.Minute)
		}
		fmt.Println("msg chan closed for stand up ", teamID)
	}()

	for _, m := range t.Users {
		logrus.Info("moving on to user ", m.Name)
		userTimer = time.AfterFunc(time.Minute*2, func() {
			userDone <- struct{}{}
		})
		msg := msgBuilder.Text(fmt.Sprintf(msgFormat, m.ID)).Mention(m.ID, m.Name, 0).Build()
		// ask for an update.
		if _, err := sr.chat.SendMessageToTeam(teamID, msg); err != nil {
			logrus.Errorf("failed to send hangout message", err)
		}

		<-userDone
		logrus.Info("user done ", m.Name, "stopping timer")
		userTimer.Stop()
	}
	logrus.Info("standup complete")
	msgBuilder.Text("The standup is now complete. To see what was recorded during this standup use *@sprintbot standup log* \n")
	if _, err := sr.chat.SendMessageToTeam(teamID, msgBuilder.Build()); err != nil {
		logrus.Errorf("failed to send message to team ", err)
	}
	delete(standupThread, teamID)
}

func (sr *StandUpRunner) Announce(teamID string, minutesBefore time.Duration) {
	msg := NewMessageBuilder().Text(fmt.Sprintf("stand up will start in %d minutes", minutesBefore/time.Minute)).Build()
	id, err := sr.chat.SendMessageToTeam(teamID, msg)
	if err != nil {
		logrus.Error("failed to announce standup ", err)
	}
	standupThread[teamID] = id
}

func (sr *StandUpRunner) LogStandUpMessage(cmd command, event *Event) {

}
