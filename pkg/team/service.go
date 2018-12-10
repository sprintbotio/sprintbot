package team

import (
	"github.com/pkg/errors"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
)

type Service struct {
	userRepo UserRepo
	teamRepo TeamRepo
}

func NewService(ur UserRepo, teamRepo TeamRepo) *Service  {
	return &Service{userRepo:ur, teamRepo:teamRepo}
}

func (ad *Service)RegisterAdmin(adminName, space string )(string,error)  {
	admin := domain.User{Admin:true}
	admin.Name = adminName
	admin.Team = space
	id, err := ad.userRepo.AddUser(&admin)
	if err != nil{
		return "", errors.Wrap(err, "failed to register admin")
	}
	return id, nil
}

func (ad *Service)RegisterTeam(team, room, owner string)(string,error) {
	t := domain.Team{
		Name:team,
		Members:[]string{
			owner,
		},
		ID:room,
	}
	if _ , err := ad.teamRepo.AddTeam(t); err != nil{
		return "", errors.Wrap(err, "failed to add new team")
	}
	return t.ID, nil
}

func (ad *Service)Team(id string)(*domain.Team, error)  {
	return ad.teamRepo.GetTeam(id)
}
