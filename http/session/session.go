package session

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type Session struct {
	Session *sessions.Session
	Request *http.Request
	Writer  http.ResponseWriter
}

// Save the current session.
func (s *Session) Save() error {
	return s.Session.Save(s.Request, s.Writer)
}

// Get a value from the current session.
func (s *Session) Get(name interface{}) interface{} {
	return s.Session.Values[name]
}

// GetOnce gets a value from the current session and then deletes it.
func (s *Session) GetOnce(name interface{}) interface{} {
	if x, ok := s.Session.Values[name]; ok {
		s.Delete(name)
		return x
	}
	return nil
}

// Set a value onto the current session. If a value with that name
// already exists it will be overridden with the new value.
func (s *Session) Set(name, value interface{}) {
	s.Session.Values[name] = value
}

// Delete a value from the current session.
func (s *Session) Delete(name interface{}) {
	delete(s.Session.Values, name)
}

// Clear the current session
func (s *Session) Clear() {
	for k := range s.Session.Values {
		s.Delete(k)
	}
}
