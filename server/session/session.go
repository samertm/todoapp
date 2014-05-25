package session

// Sessions keep track of the map between session ids and persons

import (
	"github.com/samertm/todoapp/engine"
)

// maps session ids to
type Session struct {
	// create a new person for the session
	Register chan string
	// sets the person to the username in the database
	Set chan Set
	// get person associated with session, on Out
	Get chan string
	// clean up session
	Delete chan string
	// output for chan 'get'
	Out chan *engine.Person

	sessions map[string]*engine.Person
}

type Set struct {
	SessionID string
	Name      string
}

func New() *Session {
	return &Session{
		Register: make(chan string),
		Set:      make(chan Set),
		Get:      make(chan string),
		Delete:   make(chan string),
		Out:      make(chan *engine.Person),
		sessions: make(map[string]*engine.Person),
	}
}

func (s Session) Run() {
	for {
		select {
		case i := <-s.Register:
			s.sessions[i] = &engine.Person{}
		case set := <-s.Set:
			p := engine.PersonStore[set.Name]
			if p != nil {
				s.sessions[set.SessionID] = p
			}
		case i := <-s.Get:
			// TODO is there a more robust way of doing this?
			if s.sessions[i] == nil {
				s.sessions[i] = &engine.Person{}
			}
			s.Out <- s.sessions[i]
		case i := <-s.Delete:
			delete(s.sessions, i)
		}
	}
}
