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
	standUpRepo domain.StandUpRepo
}

var standupThread = map[string]string{}

func NewStandUpRunner(chat *Service, ts *team.Service, sr domain.StandUpRepo) *StandUpRunner {
	return &StandUpRunner{chat: chat, teamService: ts, standUpRepo: sr}
}

func (sr *StandUpRunner) InProgress(teamID string) bool {
	_, ok := standupThread[teamID]
	return ok
}

func (sr *StandUpRunner) Run(teamID, tz string, msgChan chan domain.StandUpMsg) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		logrus.Errorf("failed to load location cannot run stand up ", err)
	}
	now := time.Now().In(loc)
	var (
		userDone                       = make(chan struct{})
		userTimer                      *time.Timer
		standUpID                      = fmt.Sprintf("%s-%d-%02d-%02d", teamID, now.Year(), now.Month(), now.Day())
		currentUserID, currentUserName string
		present                        bool
	)
	defer close(userDone)
	msgBuilder := NewMessageBuilder()
	// create a stand up rep

	var startMsg = "Stand up is now starting. Each team member will be asked for their status update. If a team member doesn't respond within 2 mins they will be considered absent." +
		"Due to the limitations of hangouts, *Please ensure to mention @sprintbot in your replies*\n"
	threadID := standupThread[teamID]
	msg := msgBuilder.Text(startMsg).Thread(threadID).Build()
	if _, err := sr.chat.SendMessageToTeam(teamID, msg); err != nil {
		logrus.Error("failed to start standup ", err)
		return
	}

	// set to midnight to help querying
	standUp := domain.StandUp{
		ID:        standUpID,
		StartTime: now.Unix(),
		TeamID:    teamID,
	}
	if err := sr.standUpRepo.SaveUpdate(&standUp); err != nil {
		logrus.Errorf("failed to save stand up", err)
		msgBuilder.Text("I am unable to log this stand up due to an unexpected error. I will continue with the stand up")
		if _, err := sr.chat.SendMessageToTeam(teamID, msg); err != nil {
			logrus.Error("failed to start standup ", err)
		}
	}

	t, err := sr.teamService.PopulateTeam(teamID)
	if err != nil {
		logrus.Errorf("failed to get a populated team", err)
		return
	}
	msgFormat := "<%s> please give your update. \n Remember to *mention @sprintbot in your reply* \n Admins can choose to skip this user using @sprinbot skip"

	go func() {
		// read any status messages
		for m := range msgChan {
			standUp, err := sr.standUpRepo.Get(standUpID)
			if err != nil {
				logrus.Errorf("failed to get stand up ", err)
			}

			if strings.TrimSpace(strings.ToLower(m.Msg)) == "skip" {
				logrus.Info("skipping ", currentUserName)
				userTimer.Reset(time.Second)
				continue
			}
			if m.UserID != currentUserID {
				msgBuilder.Text("sorry <" + m.UserID + "> only <" + currentUserID + "> can add their own status during the stand up")
				if _, err := sr.chat.SendMessageToTeam(teamID, msgBuilder.Build()); err != nil {
					logrus.Error("failed to send message ", err)
				}
				continue
			}
			present = true
			standUp.Present = append(standUp.Present, currentUserName)
			// add a stand up log reset the timer for the user (in case they want to add more detail)
			if strings.TrimSpace(strings.ToLower(m.Msg)) == "complete" {
				logrus.Info("got a complete message")
				msgBuilder.Text("thanks for your status <" + m.UserID + ">")
				if _, err := sr.chat.SendMessageToTeam(teamID, msgBuilder.Build()); err != nil {
					logrus.Errorf("failed to send message to team ", err)
				}
				logrus.Info("sending user done")
				// TODO think about allowing time for ppl to respond with questions.
				userTimer.Reset(time.Second)
				logrus.Info("reset timer ")
				continue
			}
			// actual standup msg so log it

			standUp.Log = append(standUp.Log, &domain.StandUpLog{
				UserID:   m.UserID,
				UserName: m.UserName,
				Message:  m.Msg,
			})
			if err := sr.standUpRepo.SaveUpdate(standUp); err != nil {
				logrus.Errorf("failed to get stand up ", err)
				msgBuilder.Text("I am unable to log this stand up due to an unexpected error. I will continue with the stand up")
				if _, err := sr.chat.SendMessageToTeam(teamID, msgBuilder.Build()); err != nil {
					logrus.Error("failed to save stand up log ", err)
				}
			}
			msgBuilder.Text("Logged status. Are you finished? Use: *@sprintbot complete* to end your status. I will naturally move on in 1 minute")
			if _, err := sr.chat.SendMessageToTeam(teamID, msgBuilder.Build()); err != nil {
				logrus.Errorf("failed to send message to team ", err)
				msgBuilder.Text("I am unable to log this stand up due to an unexpected error. I will continue with the stand up")
				if _, err := sr.chat.SendMessageToTeam(teamID, msgBuilder.Build()); err != nil {
					logrus.Error("failed to save stand up log ", err)
				}
			}
			logrus.Info("resetting the timer")
			userTimer.Reset(time.Minute)
		}
		fmt.Println("msg chan closed for stand up ", teamID)
	}()

	for _, m := range t.Users {
		present = false
		currentUserID = m.ID
		currentUserName = m.Name
		logrus.Info("moving on to user ", m.Name)
		standUp, err := sr.standUpRepo.Get(standUpID)
		if err != nil {
			logrus.Errorf("failed to get stand up ", err)
		}
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
		logrus.Info("present ", present, currentUserName)
		if !present {
			standUp.Absent = append(standUp.Absent, currentUserName)
			if err := sr.standUpRepo.SaveUpdate(standUp); err != nil {
				logrus.Errorf("failed to save stand up ", err)
			}
		}
	}
	//TODO allow for other comments
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
