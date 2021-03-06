package standup

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
)

type Service struct {
	ts            domain.TeamRepo
	sr            domain.ScheduleRepo
	sl            domain.StandUpRepo
	standUpRunner Runner
	checkInterval time.Duration
}

func NewStandUpService(ts domain.TeamRepo, sr domain.ScheduleRepo, sRepo domain.StandUpRepo, checkInterval time.Duration) *Service {
	return &Service{ts: ts, sr: sr, sl: sRepo, checkInterval: checkInterval}
}

var DefaultCheckInterval = time.Second * 60

func (ss *Service) SaveTime(teamID string, hour, minute int64, timeZone string) error {

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

func (ss *Service) CreateStandUp(time time.Time, teamID string) (*domain.StandUp, error) {
	s := &domain.StandUp{}
	s.ID = ss.sl.GenerateID(teamID, time)
	s.TeamID = teamID
	s.Log = []*domain.StandUpLog{}
	s.Absent = []string{}
	s.Present = []string{}
	if err := ss.sl.SaveUpdate(s); err != nil {
		return nil, errors.Wrap(err, "failed to create stand up")
	}
	return s, nil
}

func (ss *Service) CreateNextStandUp(t time.Time, teamID string) {
	// has stand up ran based on timezone time t
	// if not is there one scheduled to happen today
	// if one is scheduled to happen today create a standup for today
	// otherwise create a standup for the next day
}

var standUpMsgs = map[string]chan domain.StandUpMsg{}

func (ss *Service) Schedule(ctx context.Context, runner Runner) {
	ss.standUpRunner = runner
	//TODO NEEDS TO RECOVER IF A STANDUP WAS IN PROGRESS BUT THE SERVER FAILED AND IT IS WITHIN A CERTAIN WINDOW (IE 5 MINS)
	tick := time.NewTicker(ss.checkInterval)
L:
	for {
		select {
		case <-tick.C:
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
					runner.Announce(s.TeamID, time.Minute*2)

					go func(teamID string) {
						c := time.After(time.Minute * 2)
						<-c
						standUpChan := make(chan domain.StandUpMsg)
						standUpMsgs[teamID] = standUpChan
						runner.Run(teamID, s.TimeZone, standUpChan)
						defer close(standUpChan)
					}(s.TeamID)
				}
			}
		case <-ctx.Done():
			tick.Stop()
			break L
		}
	}
}

func (ss *Service) IsStandUpInProgress(teamID string) bool {
	return ss.standUpRunner.InProgress(teamID)
}

func (ss *Service) HandleStandUpMessage(teamID, userID, userName, msg string) error {
	standUpChan := standUpMsgs[teamID]
	if nil == standUpChan {
		return errors.New("failed to log stand up message\n")
	}
	// this can block if nothing is reading from it
	standUpChan <- domain.StandUpMsg{Msg: msg, UserID: userID, UserName: userName}
	return nil
}

func (ss *Service) shouldRunStandUp(s *domain.StandupSchedule) (bool, error) {
	logrus.Debug("should run standup? ", s.TimeZone, s.Hour, s.Min)
	l, err := time.LoadLocation(s.TimeZone)
	if err != nil {
		return false, errors.Wrap(err, "failed to check time for standup ")
	}
	t := time.Now().In(l)
	// we only care about to the minute time
	minuteTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, l)
	standup := time.Date(t.Year(), t.Month(), t.Day(), int(s.Hour), int(s.Min), 0.0, 0, l)
	logrus.Debug("checking time for standup ", t.Unix(), "standup time ", standup.Unix(), t, standup, minuteTime.Unix() == standup.Unix())
	if s.PausedUntil == -1 || s.PausedUntil > standup.Unix() {
		logrus.Debug("standup is paused ")
		return false, nil
	}
	return minuteTime.Unix() == standup.Unix(), nil
}

func (ss *Service) AddStandUpLog(sid string, log *domain.StandUpLog) error {
	su, err := ss.sl.Get(sid)
	if err != nil {
		return err
	}
	// has this user already added  a status if so add it to the original
	for i, l := range su.Log {
		if l.UserID == log.UserID {
			su.Log[i].Message += "\n" + log.Message
			return ss.sl.SaveUpdate(su)
		}
	}
	su.Log = append(su.Log, log)
	return ss.sl.SaveUpdate(su)
}

func (ss *Service) RemoveAllStandUpsForTeam(tid string) error {
	return ss.sl.DeleteAllForTeam(tid)
}

func (ss *Service) LoadStandUp(teamID string, time time.Time) (*domain.StandUp, error) {
	standUpID := ss.sl.GenerateID(teamID, time)
	standUp, err := ss.sl.Get(standUpID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load stand up logs")
	}
	return standUp, nil
}

func (ss *Service) RemoveMostRecentStandUp(teamID string) error {
	standUps, err := ss.sl.List(teamID)
	if err != nil {
		return errors.Wrap(err, "failed to remove most recent stand up")
	}
	if len(standUps) == 0 {
		return errors.New("no stand ups logged")
	}

	latestStandUp := standUps[0]

	for _, s := range standUps {
		if s.StartTime > latestStandUp.StartTime {
			latestStandUp = s
		}
	}
	return ss.sl.Delete(latestStandUp.ID)
}

type ChatInterface interface {
	SendMessageToTeam(teamID string, msg string) error
}

//go:generate moq -out standUpRunner_mock.go . Runner
type Runner interface {
	Announce(teamID string, minutesBefore time.Duration)
	Run(teamID, tz string, msgChan chan domain.StandUpMsg)
	InProgress(teamID string) bool
}
