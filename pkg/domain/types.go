package domain


type User struct {
	Admin bool
	Org string
	Name string
	ID string
}

func (u User)IsAdmin()bool  {
	return u.Admin
}
