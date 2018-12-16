package team

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
)

type Service struct {
	userRepo UserRepo
	teamRepo TeamRepo
}

func NewService(ur UserRepo, teamRepo TeamRepo) *Service  {
	return &Service{userRepo:ur, teamRepo:teamRepo}
}

func (ad *Service)id(u domain.User)(string,error)  {
	if u.Name == ""{
		return "", errors.New("no username present")
	}
	if u.Team == ""{
		return "", errors.New("no team present")
	}
	return u.Team+u.Name, nil
}

func (ad *Service)RegisterAdmin(adminName, space string )(string,error)  {
	admin := domain.User{Admin:true}
	admin.Name = adminName
	admin.Team = space
	uid, err := ad.id(admin)
	if err != nil{
		return "", err
	}
	admin.ID = uid
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
	uid, err := ad.id(u)
	if err != nil{
		return err
	}
	u.ID = uid
	id, err := ad.userRepo.AddUser(&u)
	if err != nil{
		return errors.Wrap(err, "failed to add user to team")
	}
	t.Members = append(t.Members, id)
	return ad.teamRepo.Update(t)
}

func (ad *Service)RemoveUserFromTeam(userName , teamID string)error{
	team,err := ad.teamRepo.GetTeam(teamID)
	if err != nil{
		return err
	}
	user := domain.User{ Name:userName,Team:teamID}
	userID, err := ad.id(user)
	if err != nil{
		return err
	}
	var userIndex = -1
	for i, u := range team.Members{
		logrus.Info("removing user", userID, u)
		if u == userID{
			userIndex = i
			break
		}
	}
	if userIndex == -1{
		return errors.New("that user is not a member of the " + team.Name + " team" )
	}
	if userIndex > -1 {
		team.Members = append(team.Members[:userIndex],team.Members[userIndex+1:]...)
	}
	return ad.teamRepo.Update(team)
}



func (ad *Service)IsUserInTeam(name, teamID string)(bool, error)  {
	team, err := ad.PopulateTeam(teamID)
	if err != nil{
		return false, err
	}
	for _, u := range team.Members{
		if name == u{
			return true, nil
		}
	}
	return false, nil
}
