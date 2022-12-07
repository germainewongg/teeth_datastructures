package model

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Patients []*User
	ErrorLog *log.Logger
}

type User struct {
	ID             string
	Name           string
	Username       string
	HashedPassword []byte
	Email          string
	Admin          bool
}

func (m *Users) LoadUsers() error {
	file, err := os.ReadFile("internal/model/users.json")
	if err != nil {
		m.ErrorLog.Print("Failed to open users.json file")
		return err
	}

	var patients []*User
	if err := json.Unmarshal(file, &patients); err != nil {
		m.ErrorLog.Print("Error reading users from json")
		return err
	}
	m.Patients = patients
	return nil
}

func (m *Users) find(id string) (*User, error) {
	for _, patient := range m.Patients {
		if patient.ID == id {
			return patient, nil
		}
	}

	err := errors.New("user not found")
	return nil, err
}

func (m *Users) CreateUser(name, username, password, email string) (string, error) {
	id := uuid.NewString()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		m.ErrorLog.Print("Error while hashing password")
	}
	newUser := &User{
		ID:             id,
		Name:           name,
		HashedPassword: hashedPassword,
		Username:       username,
		Email:          email,
		Admin:          false,
	}

	m.Patients = append(m.Patients, newUser)

	// Write to the json file
	file, _ := json.MarshalIndent(m.Patients, "", " ")
	err = os.WriteFile("internal/model/users.json", file, 0644)
	if err != nil {
		m.ErrorLog.Print("Writing to file failed. User craetion")
		return "", err
	}

	return id, nil
}

func (m *Users) CreateAdmin() (string, error) {
	id := uuid.NewString()
	password := "Password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		m.ErrorLog.Print("Error while hashing password")
	}
	newUser := &User{
		ID:             id,
		Name:           "Admin",
		HashedPassword: hashedPassword,
		Username:       "admin",
		Admin:          true,
	}

	m.Patients = append(m.Patients, newUser)

	// Write to the json file
	file, _ := json.MarshalIndent(m.Patients, "", " ")
	err = os.WriteFile("internal/model/users.json", file, 0644)
	if err != nil {
		m.ErrorLog.Print("Writing to file failed. User craetion")
		return "", err
	}

	return id, nil
}

func (m *Users) GetUser(id string) (*User, error) {
	err := m.LoadUsers()
	if err != nil {
		m.ErrorLog.Print("Getuser failed to load file.")
		return nil, err
	}

	for _, patient := range m.Patients {
		if patient.ID == id {
			return patient, nil
		}
	}

	err = errors.New("getuser: user not found")
	m.ErrorLog.Print("GetUser: user not found")
	return nil, err

}

func (m *Users) UpdateUser(id, name, username, email string) error {
	err := m.LoadUsers()
	if err != nil {
		m.ErrorLog.Print("updateUser: file open failure")
		return err
	}

	m.ErrorLog.Printf("USERID: %v", id)
	userUpdate, err := m.find(id)
	if err != nil {
		m.ErrorLog.Print("Failed to update user.")
		m.ErrorLog.Print(err.Error())
		return err
	}
	userUpdate.Name = name
	userUpdate.Email = email
	userUpdate.Username = username
	return nil
}

func (m *Users) DeleteUser(id string) {
	m.LoadUsers()
	for i, patient := range m.Patients {
		if patient.ID == id {
			m.Patients = append(m.Patients[:i], m.Patients[i+1:]...)
			break
		}
	}

	file, _ := json.MarshalIndent(m.Patients, "", " ")
	err := os.WriteFile("internal/model/users.json", file, 0644)
	if err != nil {
		m.ErrorLog.Print("Writing to file failed. User deletion")
	}

}

// USER LOGIN CHECKS
func (m *Users) Authenticate(username, password string) (bool, *User) {
	var authenticateUser *User

	for _, patient := range m.Patients {
		if patient.Username == username {
			authenticateUser = patient
			break
		}
	}

	storedPassword := authenticateUser.HashedPassword
	err := bcrypt.CompareHashAndPassword(storedPassword, []byte(password))
	if err != nil {
		return false, nil
	}

	return true, authenticateUser
}

// Checks if the admin user has already been created.
func (m *Users) FindAdmin() bool {
	for _, patient := range m.Patients {
		if patient.Username == "admin" {
			return true
		}
	}

	return false
}
