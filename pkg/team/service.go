package team

import (
	"github.com/pkg/errors"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
)

type Service struct {
	userRepo UserRepo
}

func NewService(ur UserRepo) *Service  {
	return &Service{userRepo:ur}
}

func (ad *Service)RegisterAdmin(adminName, orgName string )error  {
	admin := domain.User{Admin:true}
	admin.Org = orgName
	admin.Name = adminName
	if _, err := ad.userRepo.AddUser(admin); err != nil{
		return errors.Wrap(err, "failed to register admin")
	}
	return nil
}
