package domain


type User struct {
	Admin bool
	Org string
	Name string
	ID string
	Team string
	Timezone struct{
		UTCOffset int
		Name string
	}
}

func (u User)IsAdmin()bool  {
	return u.Admin
}


type Team struct{
	Name string
	Owner string
	ID string
	Members []string
	Users []*User
}


//go:generate moq -out userRepo_mock_test.go . UserRepo
type UserRepo interface {
	AddUser(u *User)(string, error)
	GetUser(id string)(*User, error)
	Update(u *User)error
}

//go:generate moq -out teamRepo_mock_test.go . TeamRepo
type TeamRepo interface {
	AddTeam(team Team)(string,error)
	GetTeam(id string)(*Team, error )
	Update(t *Team)error
	Delete(id string)error
}