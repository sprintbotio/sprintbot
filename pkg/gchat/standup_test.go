package gchat_test

import (
	"testing"

	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
	"github.com/sprintbot.io/sprintbot/pkg/gchat"
	"github.com/sprintbot.io/sprintbot/pkg/standup"
	"github.com/sprintbot.io/sprintbot/pkg/user"
)

func TestStandUpUseCase_LogStandUpStatus(t *testing.T) {
	userID := "someuser"
	teamId := "someteam"
	cases := []struct {
		Name        string
		UserRepo    func() domain.UserRepo
		TeamRepo    func() domain.TeamRepo
		StandUpRepo func() domain.StandUpRepo
		ChatCmd     func() chat.Command
		Event       func() *gchat.Event
		ExpectError bool
		Validate    func(t *testing.T, resp string)
	}{
		{
			Name: "test stand up logged for current day when no stand up present",
			UserRepo: func() domain.UserRepo {
				m := &domain.UserRepoMock{
					GetUserFunc: func(id string) (i *domain.User, e error) {
						return &domain.User{
							ID:       userID,
							Timezone: "Europe/Dublin",
						}, nil
					},
				}
				return m
			},
			StandUpRepo: func() domain.StandUpRepo {
				m := &domain.StandUpRepoMock{
					GetFunc: func(sid string) (up *domain.StandUp, e error) {
						return nil, domain.NewNotFoundErr("standup")
					},
				}
				return m
			},
			ChatCmd: func() chat.Command {
				return chat.Command{
					TeamID: teamId,
				}
			},
		},
		//{
		//	Name: "test stand up logged for next day when stand up already complete",
		//},
		//{
		//	Name: "test stand up logged for current day when stand up present but not started",
		//},
		//{
		//	Name: "test stand up logged when stand up status logged for current day when stand up in progress",
		//},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			us := user.NewService(tc.UserRepo())
			ss := standup.NewStandUpService(tc.TeamRepo(), nil, tc.StandUpRepo())
			standUp := gchat.NewStandUpUseCase(ss, us)
			resp, err := standUp.LogStandUpStatus(tc.ChatCmd(), tc.Event())
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an err but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatal("did not expect an err but got one ", err)
			}
			if tc.Validate != nil {
				tc.Validate(t, resp)
			}

		})
	}
}
