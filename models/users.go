package models

type User struct {
	ID   int64
	Name string
	Bio  string
}

func NewUser(id int64, name, bio string) *User {
	return &User{
		ID:   id,
		Name: name,
		Bio:  bio,
	}
}
