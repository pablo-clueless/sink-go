package store

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Email struct {
	ID         string    `json:"id"`
	To         string    `json:"to"`
	From       string    `json:"from"`
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	HTML       string    `json:"html"`
	ReceivedAt time.Time `json:"receivedAt"`
}

type Store struct {
	mu     sync.RWMutex
	emails map[string][]Email
}

func New() *Store {
	return &Store{emails: make(map[string][]Email)}
}

func (s *Store) Add(email Email) Email {
	email.ID = uuid.NewString()
	email.ReceivedAt = time.Now()
	s.mu.Lock()
	s.emails[email.To] = append(s.emails[email.To], email)
	s.mu.Unlock()
	return email
}

func (s *Store) GetByAddress(addr string) []Email {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.emails[addr]
}

func (s *Store) DeleteByAddress(addr string) {
	s.mu.Lock()
	delete(s.emails, addr)
	s.mu.Unlock()
}
