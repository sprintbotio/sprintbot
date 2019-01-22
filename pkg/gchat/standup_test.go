package gchat_test

import (
	"reflect"
	"testing"

	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/domain"
	"github.com/sprintbot.io/sprintbot/pkg/gchat"
	"github.com/sprintbot.io/sprintbot/pkg/standup"
	"github.com/sprintbot.io/sprintbot/pkg/user"
)

func TestStandUpUseCase_LogStandUpStatus(t *testing.T) {
	teamId := "someteam"
	cases := []struct {
		Name        string
		SchduleRepo func() domain.ScheduleRepo
		TeamRepo    func() domain.TeamRepo
		StandUpRepo func() domain.StandUpRepo
		ChatCmd     func() chat.Command
		Event       func() *gchat.Event
		UserRepo    func() domain.UserRepo
		ExpectError bool
		Validate    func(t *testing.T, resp string, err error)
	}{
		{
			Name: "test logging stand up fails when not a direct message to the bot",
			SchduleRepo: func() domain.ScheduleRepo {
				return &domain.ScheduleRepoMock{}
			},
			StandUpRepo: func() domain.StandUpRepo {
				m := &domain.StandUpRepoMock{}
				return m
			},
			TeamRepo: func() domain.TeamRepo {
				return &domain.TeamRepoMock{}
			},
			UserRepo: func() domain.UserRepo {
				return &domain.UserRepoMock{}
			},
			ChatCmd: func() chat.Command {
				return chat.Command{
					TeamID: teamId,
				}
			},
			Event: func() *gchat.Event {
				return &gchat.Event{
					Space: struct {
						Name        string `json:"name"`
						DisplayName string `json:"displayName"`
						Type        string `json:"type"`
					}{Name: "test", DisplayName: "test", Type: "ROOM"},
				}
			},
			ExpectError: true,
			Validate: func(t *testing.T, resp string, err error) {
				if resp != "" {
					t.Fatalf("expected the response to be empty but it was %s", resp)
				}
				_, ok := err.(*domain.NotDirectMessageErr)
				if !ok {
					t.Fatalf("expectd an NotDirectMessageErr response but got %s", reflect.TypeOf(resp))
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
				tc.Validate(t, resp, err)
			}

		})
	}
}
