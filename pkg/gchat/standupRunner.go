package gchat

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
		standUpID                      = sr.standUpRepo.GenerateID(teamID, now)
		currentUserID, currentUserName string
		present                        bool
		standUpErr                     error
		standup                        *domain.StandUp
	)
	defer close(userDone)
	defer delete(standupThread, teamID)

	t, err := sr.teamService.PopulateTeam(teamID)
	if err != nil {
		logrus.Errorf("failed to get a populated team", err)
		return
	}
	msgBuilder := NewMessageBuilder()
	// create a stand up rep
	teamMembers := sr.buildTeamMemberList(t)
	var startMsg = teamMembers + "The stand up is now starting. Each team member will be asked for their status update. If a team member doesn't respond within 2 minutes they will be considered absent." +
		"*Please ensure to mention @sprintbot in your replies*\n If you wish to have a comment logged after a team member's status is given, ensure to also mention @sprintbot in the comment"
	threadID := standupThread[teamID]
	msg := msgBuilder.Text(startMsg).Thread(threadID).Build()
	if _, err := sr.chat.SendMessageToTeam(teamID, msg); err != nil {
		logrus.Error("failed to start standup ", err)
		return
	}

	// look for existing stand up as users may have already logged statuses
	standup, standUpErr = sr.standUpRepo.Get(standUpID)
	if standUpErr != nil && !domain.IsNotFoundErr(standUpErr) {
		logrus.Error("failed to find an existing stand up", err)
		return
	}
	if standup == nil {
		// set to midnight to help querying
		standup = &domain.StandUp{
			ID:        standUpID,
			StartTime: now.Unix(),
			TeamID:    teamID,
		}
	}
	if err := sr.standUpRepo.SaveUpdate(standup); err != nil {
		logrus.Errorf("failed to save stand up", err)
		msgBuilder.Text("I am unable to log this stand up due to an unexpected error. I will continue with the stand up")
		if _, err := sr.chat.SendMessageToTeam(teamID, msg); err != nil {
			logrus.Error("failed to start standup ", err)
		}
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
				// This could be a comment based on what the user said. We should record this
				if len(standUp.Log) > 0 {
					// get last log and add a comment to it
					standUp.Log[len(standUp.Log)-1].Comments = append(standUp.Log[len(standUp.Log)-1].Comments, m.Msg)
					if err := sr.standUpRepo.SaveUpdate(standUp); err != nil {
						logrus.Error("failed to save the stand up comment ", err)
					}

				} else {
					msgBuilder.Text("sorry <" + m.UserID + "> only <" + currentUserID + "> has not yet added a status. Please hold any comments until the a status has been added")
					if _, err := sr.chat.SendMessageToTeam(teamID, msgBuilder.Build()); err != nil {
						logrus.Error("failed to send message ", err)
					}
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
		statusLogged := false
		logrus.Info("moving on to user ", m.Name)
		standUp, err := sr.standUpRepo.Get(standUpID)
		if err != nil {
			logrus.Errorf("failed to get stand up ", err)
		}
		userTimer = time.AfterFunc(time.Minute*2, func() {
			userDone <- struct{}{}
		})
		existingStatusMsg := "<" + currentUserID + "> logged the following status before the stand up started: \n %s"
		// check if user has already logged a status. Show it if so and then move on
		for _, status := range standUp.Log {
			logrus.Info("status found ", status.UserName, status.UserID)
			if status.UserID == currentUserID {
				msg := msgBuilder.Text(fmt.Sprintf(existingStatusMsg, status.Message)).Build()
				if _, err := sr.chat.SendMessageToTeam(teamID, msg); err != nil {
					logrus.Errorf("failed to send message to room ", err)
				}
				statusLogged = true
				userTimer.Reset(time.Second * 2)
				break
			}
		}
		if !statusLogged {
			msg := msgBuilder.Text(fmt.Sprintf(msgFormat, m.ID)).Mention(m.ID, m.Name, 0).Build()
			// ask for an update.
			if _, err := sr.chat.SendMessageToTeam(teamID, msg); err != nil {
				logrus.Errorf("failed to send hangout message", err)
			}
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
	msgBuilder.Text("any final comments you want logged before we wrap up the standup? If nothing is sent I will end the standup in 1 minute")
	if _, err := sr.chat.SendMessageToTeam(teamID, msgBuilder.Build()); err != nil {
		logrus.Errorf("failed to send message to team ", err)
	}
	// wait a final minute for any final comments
	<-time.Tick(1 * time.Minute)
	logrus.Info("standup complete")
	msgBuilder.Text("The standup is now complete. To see what was recorded during this standup use *@sprintbot standup log* \n")
	if _, err := sr.chat.SendMessageToTeam(teamID, msgBuilder.Build()); err != nil {
		logrus.Errorf("failed to send message to team ", err)
	}
	standup, err = sr.standUpRepo.Get(standup.ID)
	if err != nil {
		logrus.Errorf("failed to send message to team ", err)
	}
	standup.EndTime = time.Now().In(loc).Unix()
	if err := sr.standUpRepo.SaveUpdate(standup); err != nil {
		logrus.Errorf("failed to save stand up ", err)
	}

}

func (sr *StandUpRunner) Announce(teamID string, minutesBefore time.Duration) {
	t, err := sr.teamService.PopulateTeam(teamID)
	if err != nil {
		logrus.Error("failed to announce standup ", err)
	}
	teamMembers := sr.buildTeamMemberList(t)
	msg := NewMessageBuilder().Text(fmt.Sprintf("%s stand up will start in %d minutes", teamMembers, minutesBefore/time.Minute)).Build()
	id, err := sr.chat.SendMessageToTeam(teamID, msg)
	if err != nil {
		logrus.Error("failed to announce standup ", err)
	}
	standupThread[teamID] = id
}

func (sr *StandUpRunner) buildTeamMemberList(team *domain.Team) string {
	teamMembers := ""
	for _, u := range team.Users {
		teamMembers += "<" + u.ID + "> "
	}
	return teamMembers
}
