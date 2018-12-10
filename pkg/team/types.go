package team

import "github.com/sprintbot.io/sprintbot/pkg/domain"

//go:generate moq -out userRepo_mock_test.go . UserRepo
type UserRepo interface {
	AddUser(u *domain.User)(string, error)
}

//go:generate moq -out teamRepo_mock_test.go . TeamRepo
type TeamRepo interface {
	AddTeam(team domain.Team)(string,error)
	GetTeam(id string)(*domain.Team, error )
}
