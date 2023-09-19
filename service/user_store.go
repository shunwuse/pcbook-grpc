package service

import "sync"

type UserStore interface {
	// Save saves a user to the store
	Save(user *User) error
	// Find finds a user by username
	Find(username string) (*User, error)
}

type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*User
}

// NewInMemoryUserStore returns a new InMemoryUserStore
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*User),
	}
}

// Save saves a user to the store
func (s *InMemoryUserStore) Save(user *User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.users[user.Username] != nil {
		return ErrAlreadyExists
	}

	s.users[user.Username] = user.Clone()

	return nil
}

// Find finds a user by username
func (s *InMemoryUserStore) Find(username string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user := s.users[username]
	if user == nil {
		return nil, nil
	}

	return user.Clone(), nil
}
