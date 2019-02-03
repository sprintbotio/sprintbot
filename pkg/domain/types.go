package domain

import (
	"time"
)

type User struct {
	Admin    bool
	Org      string
	Name     string
	ID       string
	Team     string
	Role     string // admin | member | general
	Timezone string
}

func (u User) IsAdmin() bool {
	return u.Admin
}

type Team struct {
	Name    string
	Owner   string
	ID      string
	Members []string
	Users   []*User
}

type StandupSchedule struct {
	Hour        int64
	Min         int64
	TimeZone    string
	TeamID      string
	PausedUntil int64
}

type StandUp struct {
	ID        string
	TeamID    string
	Log       []*StandUpLog
	StartTime int64
	EndTime   int64
	Absent    []string
	Present   []string
	DateTime  string
}

func NewStandUP(teamID string) *StandUp {
	return &StandUp{
		TeamID:  teamID,
		Log:     []*StandUpLog{},
		Absent:  []string{},
		Present: []string{},
	}
}

type StandUpMsg struct {
	UserID   string
	UserName string
	Msg      string
}

type StandUpLog struct {
	UserName string
	UserID   string
	Message  string
	Comments []*StandUpLog
}

//go:generate moq -out mockStandUpRepo.go . StandUpRepo
type StandUpRepo interface {
	SaveUpdate(s *StandUp) error
	Get(sid string) (*StandUp, error)
	FindByTeam(tid string) (*StandUp, error)
	GenerateID(teamID string, t time.Time) string
	Delete(id string) error
	List(teamID string) ([]*StandUp, error)
	DeleteAllForTeam(id string) error
}

//go:generate moq -out mockUserRepo_mock.go . UserRepo
type UserRepo interface {
	AddUser(u *User) (string, error)
	GetUser(id string) (*User, error)
	Update(u *User) error
	Delete(id string) error
}

//go:generate moq -out mockTeamRepo_mock.go . TeamRepo
type TeamRepo interface {
	AddTeam(team Team) (string, error)
	GetTeam(id string) (*Team, error)
	Update(t *Team) error
	Delete(id string) error
}

//go:generate moq -out mockScheduleRepo_mock.go . ScheduleRepo
type ScheduleRepo interface {
	SaveUpdate(teamID string, schedule StandupSchedule) error
	List() ([]*StandupSchedule, error)
}
