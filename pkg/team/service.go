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

func (ad *Service)PopulateTeam(id string)(*domain.Team, error)  {
	t , err := ad.teamRepo.GetTeam(id)
	if err != nil{
		return nil, err
	}
	userNames := []string{}
	for _, m := range t.Members{
		u, err := ad.userRepo.GetUser(m)
		if err != nil{
			return nil, err
		}
		userNames = append(userNames, u.Name)
	}
	t.Members = userNames
	return t, nil
}

func (ad *Service)AddUserToTeam(name, teamid string)error  {
	t, err := ad.teamRepo.GetTeam(teamid)
	if err != nil{
		return err
	}
	u := domain.User{Admin:false}
	u.Name = name
	u.Team = teamid
	id, err := ad.userRepo.AddUser(&u)
	if err != nil{
		return errors.Wrap(err, "failed to add user to team")
	}
	t.Members = append(t.Members, id)
	return ad.teamRepo.Update(t)
}
