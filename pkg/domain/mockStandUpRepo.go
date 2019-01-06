// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package domain

import (
	"sync"
)

var (
	lockStandUpRepoMockFindByTeam sync.RWMutex
	lockStandUpRepoMockGet        sync.RWMutex
	lockStandUpRepoMockSaveUpdate sync.RWMutex
)

// StandUpRepoMock is a mock implementation of StandUpRepo.
//
//     func TestSomethingThatUsesStandUpRepo(t *testing.T) {
//
//         // make and configure a mocked StandUpRepo
//         mockedStandUpRepo := &StandUpRepoMock{
//             FindByTeamFunc: func(tid string) (*StandUp, error) {
// 	               panic("TODO: mock out the FindByTeam method")
//             },
//             GetFunc: func(sid string) (*StandUp, error) {
// 	               panic("TODO: mock out the Get method")
//             },
//             SaveUpdateFunc: func(s *StandUp) error {
// 	               panic("TODO: mock out the SaveUpdate method")
//             },
//         }
//
//         // TODO: use mockedStandUpRepo in code that requires StandUpRepo
//         //       and then make assertions.
//
//     }
type StandUpRepoMock struct {
	// FindByTeamFunc mocks the FindByTeam method.
	FindByTeamFunc func(tid string) (*StandUp, error)

	// GetFunc mocks the Get method.
	GetFunc func(sid string) (*StandUp, error)

	// SaveUpdateFunc mocks the SaveUpdate method.
	SaveUpdateFunc func(s *StandUp) error

	// calls tracks calls to the methods.
	calls struct {
		// FindByTeam holds details about calls to the FindByTeam method.
		FindByTeam []struct {
			// Tid is the tid argument value.
			Tid string
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Sid is the sid argument value.
			Sid string
		}
		// SaveUpdate holds details about calls to the SaveUpdate method.
		SaveUpdate []struct {
			// S is the s argument value.
			S *StandUp
		}
	}
}

// FindByTeam calls FindByTeamFunc.
func (mock *StandUpRepoMock) FindByTeam(tid string) (*StandUp, error) {
	if mock.FindByTeamFunc == nil {
		panic("StandUpRepoMock.FindByTeamFunc: method is nil but StandUpRepo.FindByTeam was just called")
	}
	callInfo := struct {
		Tid string
	}{
		Tid: tid,
	}
	lockStandUpRepoMockFindByTeam.Lock()
	mock.calls.FindByTeam = append(mock.calls.FindByTeam, callInfo)
	lockStandUpRepoMockFindByTeam.Unlock()
	return mock.FindByTeamFunc(tid)
}

// FindByTeamCalls gets all the calls that were made to FindByTeam.
// Check the length with:
//     len(mockedStandUpRepo.FindByTeamCalls())
func (mock *StandUpRepoMock) FindByTeamCalls() []struct {
	Tid string
} {
	var calls []struct {
		Tid string
	}
	lockStandUpRepoMockFindByTeam.RLock()
	calls = mock.calls.FindByTeam
	lockStandUpRepoMockFindByTeam.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *StandUpRepoMock) Get(sid string) (*StandUp, error) {
	if mock.GetFunc == nil {
		panic("StandUpRepoMock.GetFunc: method is nil but StandUpRepo.Get was just called")
	}
	callInfo := struct {
		Sid string
	}{
		Sid: sid,
	}
	lockStandUpRepoMockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	lockStandUpRepoMockGet.Unlock()
	return mock.GetFunc(sid)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedStandUpRepo.GetCalls())
func (mock *StandUpRepoMock) GetCalls() []struct {
	Sid string
} {
	var calls []struct {
		Sid string
	}
	lockStandUpRepoMockGet.RLock()
	calls = mock.calls.Get
	lockStandUpRepoMockGet.RUnlock()
	return calls
}

// SaveUpdate calls SaveUpdateFunc.
func (mock *StandUpRepoMock) SaveUpdate(s *StandUp) error {
	if mock.SaveUpdateFunc == nil {
		panic("StandUpRepoMock.SaveUpdateFunc: method is nil but StandUpRepo.SaveUpdate was just called")
	}
	callInfo := struct {
		S *StandUp
	}{
		S: s,
	}
	lockStandUpRepoMockSaveUpdate.Lock()
	mock.calls.SaveUpdate = append(mock.calls.SaveUpdate, callInfo)
	lockStandUpRepoMockSaveUpdate.Unlock()
	return mock.SaveUpdateFunc(s)
}

// SaveUpdateCalls gets all the calls that were made to SaveUpdate.
// Check the length with:
//     len(mockedStandUpRepo.SaveUpdateCalls())
func (mock *StandUpRepoMock) SaveUpdateCalls() []struct {
	S *StandUp
} {
	var calls []struct {
		S *StandUp
	}
	lockStandUpRepoMockSaveUpdate.RLock()
	calls = mock.calls.SaveUpdate
	lockStandUpRepoMockSaveUpdate.RUnlock()
	return calls
}
