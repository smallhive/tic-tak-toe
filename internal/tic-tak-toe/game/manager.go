package game

import (
	"context"
	"sync"

	"github.com/smallhive/tic-tak-toe/internal/logger"
)

type SessionCompleteChan chan *Session

// Manager controls game sessions
type Manager struct {
	worldCounter int64

	mutex    *sync.RWMutex
	sessions map[string]*Session

	endedSessions SessionCompleteChan
}

func NewManager() *Manager {
	m := &Manager{
		sessions: make(map[string]*Session),
		mutex:    &sync.RWMutex{},
	}

	m.endedSessions = make(SessionCompleteChan, 1)

	go m.sessionCloser()

	return m
}

func (m *Manager) Session() *Session {
	s := NewSession(m.endedSessions)

	m.mutex.Lock()
	m.sessions[s.id] = s
	m.mutex.Unlock()

	return s
}

func (m *Manager) sessionCloser() {
	for {
		session, ok := <-m.endedSessions
		if !ok {
			break
		}

		logger.Info(context.Background(), "Session", session.ID(), "completed")

		m.mutex.Lock()
		delete(m.sessions, session.ID())
		m.mutex.Unlock()
	}
}
