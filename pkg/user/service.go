package user

import (
	"github.com/pkg/errors"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
)

type Service struct {
	userRepo domain.UserRepo
}

func NewService(ur domain.UserRepo)*Service  {
	return &Service{userRepo:ur}
}

func (us *Service)RegisterAdmin(adminName, uid, space string )(string,error)  {
	admin := domain.User{Admin:true}
	admin.Name = adminName
	admin.Team = space
	admin.ID = uid
	id, err := us.userRepo.AddUser(&admin)
	if err != nil{
		return "", errors.Wrap(err, "failed to register admin")
	}
	return id, nil
}

func (us *Service)UpdateTZ(uid, zone string)error  {
	u , err  := us.userRepo.GetUser(uid)
	if err != nil{
		return errors.Wrap(err, "failed to find user could not update timezone")
	}
	u.Timezone = struct {
		UTCOffset int
		Name      string
	}{UTCOffset: 1, Name:zone }
	if err := us.userRepo.Update(u); err != nil{
		return errors.Wrap(err, "failed to update timezone for user")
	}
	return nil
}

