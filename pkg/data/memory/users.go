package memory

import (
	"github.com/sprintbot.io/sprintbot/pkg/domain"
	"sync"
)

type UserRespository struct {
	*sync.Mutex
	data map[string]domain.User
}

func NewUserRespository()*UserRespository  {
	return &UserRespository{data: map[string]domain.User{}, Mutex: &sync.Mutex{}}
}

func (ur *UserRespository)AddUser(u domain.User)(string, error){
	ur.Lock()
	defer ur.Unlock()
	id := u.Org + "/"+u.Name
	ur.data[id] = u
	return id, nil
}
