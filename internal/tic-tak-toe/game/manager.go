package game

import (
	"fmt"
	"sync"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/network"
)

type SessionCompleteChan chan *Session

type Manager struct {
	worldCounter int64
	session      *Session

	mutex    *sync.RWMutex
	sessions map[int64]*Session

	endedSessions SessionCompleteChan
}

func NewManager() *Manager {
	m := &Manager{
		sessions: make(map[int64]*Session),
		mutex:    &sync.RWMutex{},
	}

	m.endedSessions = make(SessionCompleteChan, 1)

	go m.sessionCloser()

	return m
}

func (m *Manager) Session() *Session {
	if m.session == nil || m.session.IsFull() {
		hub := network.NewHub()
		go hub.Run()

		s := NewSession(hub, m.endedSessions)

		m.mutex.Lock()
		m.sessions[s.id] = s
		m.mutex.Unlock()

		m.session = s
	}

	return m.session
}

func (m *Manager) sessionCloser() {
	for session := range m.endedSessions {
		// if !ok {
		// 	break
		// }

		fmt.Println("Session", session.ID(), "completed")

		m.mutex.Lock()
		delete(m.sessions, session.ID())
		m.mutex.Unlock()
	}
}
