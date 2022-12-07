package model

type Users struct {
	patients map[string][]*User
}

type User struct {
	ID             int
	Name           string
	Age            int
	Username       string
	HashedPassword string
	Email          string
}

func (m *Users) CreateUser(name, username, password, email string, age int) (int, error) {
	return 0, nil
}
