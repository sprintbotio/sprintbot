package domain


type User struct {
	Admin bool
	Org string
	Name string
	ID string
	Team string
}

func (u User)IsAdmin()bool  {
	return u.Admin
}


type Team struct{
	Name string
	Owner string
	ID string
	Members []string
}