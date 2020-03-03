package db

type QI interface {
	UserQI() UserQI
}

type UserQI interface {
	SaveUser(user User) error
	GetUser(login string) (*User, error)
}
