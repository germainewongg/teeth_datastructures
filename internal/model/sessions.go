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

func (s *Sessions) Store(newSession *Session) {
	errorLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	s.Sessions = append(s.Sessions, newSession)
	file, _ := json.MarshalIndent(s.Sessions, "", " ")
	err := os.WriteFile("internal/model/sessions.json", file, 0644)
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
