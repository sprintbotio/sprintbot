package team

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
)

type Service struct {
	userRepo domain.UserRepo
	teamRepo domain.TeamRepo
}

func NewService(ur domain.UserRepo, teamRepo domain.TeamRepo) *Service  {
	return &Service{userRepo:ur, teamRepo:teamRepo}
}


func (ad *Service)RemoveTeam(id string)error  {
	return ad.teamRepo.Delete(id)
}

func (ad *Service)IsUserAdminForTeam(userID, teamName string)(bool,error)  {
	u, err := ad.userRepo.GetUser(userID)
	if err != nil{
		return false, err
	}
	return u.Team == teamName && u.Admin == true, nil
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
	users := []*domain.User{}
	for _, m := range t.Members{
		u, err := ad.userRepo.GetUser(m)
		if err != nil{
			return nil, err
		}
		userNames = append(userNames, u.Name)
		users = append(users, u)
	}
	t.Members = userNames
	t.Users = users
	return t, nil
}

func (ad *Service)AddUserToTeam(name, uid, teamid string)error  {
	t, err := ad.teamRepo.GetTeam(teamid)
	if err != nil{
		return err
	}
	u := domain.User{Admin:false}
	u.Name = name
	u.Team = teamid
	u.ID = uid
	id, err := ad.userRepo.AddUser(&u)
	if err != nil{
		return errors.Wrap(err, "failed to add user to team")
	}
	t.Members = append(t.Members, id)
	return ad.teamRepo.Update(t)
}

func (ad *Service)RemoveUserFromTeam(uid , teamID string)error{
	team,err := ad.teamRepo.GetTeam(teamID)
	if err != nil{
		return err
	}
	var userIndex = -1
	for i, u := range team.Members{
		logrus.Info("removing user", uid, u)
		if u == uid{
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

func (ad *Service)GetTeamForUser(userID string)(*domain.Team, error){
	u, err := ad.userRepo.GetUser(userID)
	if err != nil{
		return nil,err
	}
	t, err := ad.teamRepo.GetTeam(u.Team)
	if err != nil{
		return  nil,err
	}
	return t , nil
}



func (ad *Service)IsUserInTeam(uid, teamID string)(bool, error)  {
	team, err := ad.teamRepo.GetTeam(teamID)
	if err != nil{
		return false, err
	}
	for _, u := range team.Members{
		if uid == u{
			return true, nil
		}
	}
	return false, nil
}
