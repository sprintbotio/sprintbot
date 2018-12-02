package team

import "github.com/sprintbot.io/sprintbot/pkg/domain"

type UserRepo interface {
	AddUser(u domain.User)(string, error)
}
