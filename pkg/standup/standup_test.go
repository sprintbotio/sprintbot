package standup_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sprintbot.io/sprintbot/pkg/domain"

	"github.com/sprintbot.io/sprintbot/pkg/standup"
)

func TestService_Schedule(t *testing.T) {

	cases := []struct {
		Name             string
		TeamRepoMock     func() domain.TeamRepo
		ScheduleRepoMock func() domain.ScheduleRepo
		StandupRepoMock  func() domain.StandUpRepo
		RunnerMock       func() standup.Runner
		Validate         func(*testing.T, *domain.TeamRepoMock, *domain.ScheduleRepoMock, *domain.StandUpRepoMock)
	}{
		{
			Name: "test stand up scheduled as expected",
			TeamRepoMock: func() domain.TeamRepo {
				return &domain.TeamRepoMock{}
			},
			ScheduleRepoMock: func() domain.ScheduleRepo {
				return &domain.ScheduleRepoMock{
					ListFunc: func() (schedules []*domain.StandupSchedule, e error) {
						return []*domain.StandupSchedule{}, nil
					},
				}
			},
			StandupRepoMock: func() domain.StandUpRepo {
				return &domain.StandUpRepoMock{}
			},
			RunnerMock: func() standup.Runner {
				return &standup.RunnerMock{}
			},
			Validate: func(t *testing.T, mock *domain.TeamRepoMock, srmock *domain.ScheduleRepoMock, mock3 *domain.StandUpRepoMock) {
				if len(srmock.ListCalls()) != 1 {
					t.Fatal("expected schedule to be listed once")
				}
			},
		},
		{
			Name: "test standup not scheduled when paused",
			TeamRepoMock: func() domain.TeamRepo {
				return &domain.TeamRepoMock{}
			},
			ScheduleRepoMock: func() domain.ScheduleRepo {
				return &domain.ScheduleRepoMock{}
			},
			StandupRepoMock: func() domain.StandUpRepo {
				return &domain.StandUpRepoMock{}
			},
			RunnerMock: func() standup.Runner {
				return &standup.RunnerMock{}
			},
			Validate: func(t *testing.T, mock *domain.TeamRepoMock, mock2 *domain.ScheduleRepoMock, mock3 *domain.StandUpRepoMock) {

			},
		},
	}
	w := &sync.WaitGroup{}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			trm := tc.TeamRepoMock()
			srm := tc.ScheduleRepoMock()
			strm := tc.StandupRepoMock()
			s := standup.NewStandUpService(trm, srm, strm, time.Second*1)
			ctx, c := context.WithCancel(context.Background())
			w.Add(1)
			go func() {
				defer w.Done()
				s.Schedule(ctx, tc.RunnerMock())
			}()
			time.AfterFunc(time.Millisecond*1500, c)
			w.Wait()
			tc.Validate(t, trm.(*domain.TeamRepoMock), srm.(*domain.ScheduleRepoMock), strm.(*domain.StandUpRepoMock))
		})

	}

}
