package model

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Session struct {
	Username string
	Expiry   time.Time
}

type Sessions struct {
	Sessions []*Session
}

func (s *Session) Expired() bool {
	return s.Expiry.Before(time.Now())
}

func (s *Sessions) LoadSessions() error {
	file, err := os.ReadFile("internal/model/sessions.json")
	if err != nil {
		return err
	}

	var sessions []*Session
	if err := json.Unmarshal(file, &sessions); err != nil {
		return err
	}
	s.Sessions = sessions
	return nil
}

func (s *Sessions) Store(newSession *Session) {
	errorLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	err := s.LoadSessions()
	if err != nil {
		errorLog.Print("Failed to load existing sessions")
		return
	}
	s.Sessions = append(s.Sessions, newSession)
	file, _ := json.MarshalIndent(s.Sessions, "", " ")
	err = os.WriteFile("internal/model/sessions.json", file, 0644)
	if err != nil {
		errorLog.Print("Failed to write to json")
		return
	}

}

func (s *Sessions) RemoveSession(sessionToken string) {
	errorLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	for i, session := range s.Sessions {
		if session.Username == sessionToken {
			s.Sessions = append(s.Sessions[:i], s.Sessions[i+1:]...)
			errorLog.Print("Removed session")
			break
		}
	}
	file, _ := json.MarshalIndent(s.Sessions, "", " ")
	err := os.WriteFile("internal/model/sessions.json", file, 0644)
	if err != nil {
		errorLog.Print("Failed to write to json")
	}
}

func (s *Sessions) Validate(sessionToken string) (bool, *Session) {
	for _, session := range s.Sessions {
		if session.Username == sessionToken {
			return true, session
		}
	}

	return false, nil
}
