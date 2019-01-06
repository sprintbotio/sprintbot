package user

import (
	"github.com/pkg/errors"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
)

type Service struct {
	userRepo domain.UserRepo
}

func NewService(ur domain.UserRepo) *Service {
	return &Service{userRepo: ur}
}

func (us *Service) RegisterAdmin(adminName, uid, space string) (string, error) {
	admin := domain.User{Admin: true}
	admin.Name = adminName
	admin.Team = space
	admin.Role = "admin"
	admin.ID = uid
	id, err := us.userRepo.AddUser(&admin)
	if err != nil {
		return "", errors.Wrap(err, "failed to register admin")
	}
	return id, nil
}

func (us *Service) UpdateTZ(uid, zone string) error {
	u, err := us.userRepo.GetUser(uid)
	if err != nil {
		return errors.Wrap(err, "failed to find user could not update timezone")
	}
	u.Timezone = zone
	if err := us.userRepo.Update(u); err != nil {
		return errors.Wrap(err, "failed to update timezone for user")
	}
	return nil
}

func (us *Service) ResolveUser(id, name string) (*domain.User, error) {
	// transfer to new object?
	u, err := us.userRepo.GetUser(id)
	if err != nil && domain.IsNotFoundErr(err) {
		u = &domain.User{
			Name: name,
			ID:   id,
			Team: "general",
			Role: "general",
		}
		return u, nil
	}
	return u, err
}

func (us *Service) GetUser(id string) (*domain.User, error) {
	return us.userRepo.GetUser(id)
}

func (us *Service) DeleteUser(id string) error {
	return us.userRepo.Delete(id)
}
