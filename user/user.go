package user

type User struct {
	name string
}

func NewUser(name string) *User {
	p := new(User)
	p.name = name
	return p
}
