package chat_test

import (
	"testing"

	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
	"github.com/sprintbot.io/sprintbot/pkg/team"
)

func TestPermissions_CanUserDoCmd(t *testing.T) {
	cases := []struct {
		Name     string
		userRepo func() domain.UserRepo
		teamRepo func() domain.TeamRepo
		validate func(t *testing.T, can bool, err error)
		user     *domain.User
		cmd      []chat.Command
	}{
		{
			Name: "test admin can do everything",
			userRepo: func() domain.UserRepo {
				mu := &domain.UserRepoMock{}
				mu.GetUserFunc = func(id string) (user *domain.User, e error) {
					return &domain.User{
						ID:   id,
						Role: "admin",
					}, nil
				}
				return mu
			},
			teamRepo: func() domain.TeamRepo {
				mt := &domain.TeamRepoMock{}
				mt.GetTeamFunc = func(id string) (i *domain.Team, e error) {
					return &domain.Team{ID: id, Members: []string{"id"}}, nil
				}
				return mt
			},
			user: &domain.User{
				ID:   "id",
				Role: "admin",
			},
			cmd: []chat.Command{
				{
					TeamID: "team",
				},
			},
			validate: func(t *testing.T, can bool, err error) {
				if err != nil {
					t.Fatal("did not expect an error but got one ", err)
				}
				if !can {
					t.Fatal("expected to admin user to be able to perform task but user was denied")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			ts := team.NewService(tc.userRepo(), tc.teamRepo())
			perms := chat.NewPermissions(ts)
			for _, cmd := range tc.cmd {
				can, err := perms.CanUserDoCmd(tc.user, cmd)
				if nil != tc.validate {
					tc.validate(t, can, err)
				}
			}
		})
	}
}
