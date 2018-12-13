// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package team

import (
	"github.com/sprintbot.io/sprintbot/pkg/domain"
	"sync"
)

var (
	lockUserRepoMockAddUser sync.RWMutex
	lockUserRepoMockGetUser sync.RWMutex
)

// UserRepoMock is a mock implementation of UserRepo.
//
//     func TestSomethingThatUsesUserRepo(t *testing.T) {
//
//         // make and configure a mocked UserRepo
//         mockedUserRepo := &UserRepoMock{
//             AddUserFunc: func(u *domain.User) (string, error) {
// 	               panic("TODO: mock out the AddUser method")
//             },
//             GetUserFunc: func(id string) (*domain.User, error) {
// 	               panic("TODO: mock out the GetUser method")
//             },
//         }
//
//         // TODO: use mockedUserRepo in code that requires UserRepo
//         //       and then make assertions.
//
//     }
type UserRepoMock struct {
	// AddUserFunc mocks the AddUser method.
	AddUserFunc func(u *domain.User) (string, error)

	// GetUserFunc mocks the GetUser method.
	GetUserFunc func(id string) (*domain.User, error)

	// calls tracks calls to the methods.
	calls struct {
		// AddUser holds details about calls to the AddUser method.
		AddUser []struct {
			// U is the u argument value.
			U *domain.User
		}
		// GetUser holds details about calls to the GetUser method.
		GetUser []struct {
			// ID is the id argument value.
			ID string
		}
	}
}

// AddUser calls AddUserFunc.
func (mock *UserRepoMock) AddUser(u *domain.User) (string, error) {
	if mock.AddUserFunc == nil {
		panic("UserRepoMock.AddUserFunc: method is nil but UserRepo.AddUser was just called")
	}
	callInfo := struct {
		U *domain.User
	}{
		U: u,
	}
	lockUserRepoMockAddUser.Lock()
	mock.calls.AddUser = append(mock.calls.AddUser, callInfo)
	lockUserRepoMockAddUser.Unlock()
	return mock.AddUserFunc(u)
}

// AddUserCalls gets all the calls that were made to AddUser.
// Check the length with:
//     len(mockedUserRepo.AddUserCalls())
func (mock *UserRepoMock) AddUserCalls() []struct {
	U *domain.User
} {
	var calls []struct {
		U *domain.User
	}
	lockUserRepoMockAddUser.RLock()
	calls = mock.calls.AddUser
	lockUserRepoMockAddUser.RUnlock()
	return calls
}

// GetUser calls GetUserFunc.
func (mock *UserRepoMock) GetUser(id string) (*domain.User, error) {
	if mock.GetUserFunc == nil {
		panic("UserRepoMock.GetUserFunc: method is nil but UserRepo.GetUser was just called")
	}
	callInfo := struct {
		ID string
	}{
		ID: id,
	}
	lockUserRepoMockGetUser.Lock()
	mock.calls.GetUser = append(mock.calls.GetUser, callInfo)
	lockUserRepoMockGetUser.Unlock()
	return mock.GetUserFunc(id)
}

// GetUserCalls gets all the calls that were made to GetUser.
// Check the length with:
//     len(mockedUserRepo.GetUserCalls())
func (mock *UserRepoMock) GetUserCalls() []struct {
	ID string
} {
	var calls []struct {
		ID string
	}
	lockUserRepoMockGetUser.RLock()
	calls = mock.calls.GetUser
	lockUserRepoMockGetUser.RUnlock()
	return calls
}