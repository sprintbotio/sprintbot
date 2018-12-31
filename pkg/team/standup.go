package team

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
)

type StandUpService struct {
	ts            domain.TeamRepo
	sr            domain.ScheduleRepo
	sl            domain.ScheduleRepo
	standUpRunner StandUpRunner
}

func NewStandUpService(ts domain.TeamRepo, sr domain.ScheduleRepo, srun StandUpRunner) *StandUpService {
	return &StandUpService{ts: ts, sr: sr, standUpRunner: srun}
}

type RepeatSchedule uint8

const (
	RepeatEveryDay = RepeatSchedule(iota)
	RepeatEveryWeekDay
)

func (ss *StandUpService) SaveTime(teamID string, hour, minute int64, timeZone string) error {
	fmt.Println(hour, minute)

	_, err := ss.ts.GetTeam(teamID)
	if err != nil {
		return err
	}
	schedule := domain.StandupSchedule{
		TeamID:   teamID,
		Hour:     hour,
		Min:      minute,
		TimeZone: timeZone,
	}
	return ss.sr.SaveUpdate(teamID, schedule)
}

var standUpMsgs = map[string]chan domain.StandUpMsg{}

func (ss *StandUpService) Schedule(ctx context.Context) {
	//TODO NEEDS TO RECOVER IF A STANDUP WAS IN PROGRESS BUT THE SERVER FAILED AND IT IS WITHIN A CERTAIN WINDOW (IE 5 MINS)
	tick := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-tick.C:
			logrus.Info("checking stand up schedule")
			schedules, err := ss.sr.List()
			if err != nil {
				logrus.Error("failed to check standups ", err)
				return
			}
			for _, s := range schedules {
				should, err := ss.shouldRunStandUp(s)
				if err != nil {
					logrus.Errorf("failed to check if should run standup")
					return
				}
				if should {
					logrus.Info("stand up will run ")
					ss.standUpRunner.Announce(s.TeamID, time.Minute*2)
					go func(teamID string) {
						c := time.After(time.Minute * 2)
						<-c
						standUpChan := make(chan domain.StandUpMsg)
						standUpMsgs[teamID] = standUpChan
						ss.standUpRunner.Run(teamID, standUpChan)
						defer close(standUpChan)
					}(s.TeamID)
				}
			}
		case <-ctx.Done():
			tick.Stop()
		}
	}
}

func (ss *StandUpService) IsStandUpInProgress(teamID string) bool {
	return ss.standUpRunner.InProgress(teamID)
}

func (ss *StandUpService) LogStandUpMessage(teamID, userID, msg string) error {
	standUpChan := standUpMsgs[teamID]
	if nil == standUpChan {
		return errors.New("failed to log stand up message\n")
	}
	// this can block if nothing is reading from it
	standUpChan <- domain.StandUpMsg{Msg: msg, UserID: userID}
	return nil
}

func (ss *StandUpService) shouldRunStandUp(s *domain.StandupSchedule) (bool, error) {
	logrus.Info("should run standup? ", s.TimeZone, s.Hour, s.Min)
	l, err := time.LoadLocation(s.TimeZone)
	if err != nil {
		return false, errors.Wrap(err, "failed to check time for standup ")
	}
	t := time.Now().In(l)
	// we only care about to the minute time
	minuteTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, l)
	standup := time.Date(t.Year(), t.Month(), t.Day(), int(s.Hour), int(s.Min), 0.0, 0, l)
	logrus.Info("checking time for standup ", t.Unix(), "standup time ", standup.Unix(), t, standup, minuteTime.Unix() == standup.Unix())
	return minuteTime.Unix() == standup.Unix(), nil
}

type ChatInterface interface {
	SendMessageToTeam(teamID string, msg string) error
}

type StandUpRunner interface {
	Announce(teamID string, minutesBefore time.Duration)
	Run(teamID string, msgChan chan domain.StandUpMsg)
	InProgress(teamID string) bool
}
