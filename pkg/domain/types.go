package domain

type User struct {
	Admin    bool
	Org      string
	Name     string
	ID       string
	Team     string
	Role     string // admin | member | observer
	Timezone struct {
		UTCOffset int
		Name      string
	}
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
	Hour     int64
	Min      int64
	TimeZone string
	TeamID   string
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
}

type StandUpRepo interface {
	SaveUpdate(s *StandUp) error
	Get(sid string) (*StandUp, error)
}

//go:generate moq -out userRepo_mock_test.go . UserRepo
type UserRepo interface {
	AddUser(u *User) (string, error)
	GetUser(id string) (*User, error)
	Update(u *User) error
}

//go:generate moq -out teamRepo_mock_test.go . TeamRepo
type TeamRepo interface {
	AddTeam(team Team) (string, error)
	GetTeam(id string) (*Team, error)
	Update(t *Team) error
	Delete(id string) error
}

type ScheduleRepo interface {
	SaveUpdate(teamID string, schedule StandupSchedule) error
	List() ([]*StandupSchedule, error)
}
