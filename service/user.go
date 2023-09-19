package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User contains the user's information
type User struct {
	Username string
	Hash     string
	Role     string
}

// NewUser returns a new user
func NewUser(username string, password string, role string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	user := &User{
		Username: username,
		Hash:     string(hash),
		Role:     role,
	}

	return user, nil
}

// IsCorrectPassword returns true if the password is correct
func (u *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(password))
	return err == nil
}

// Clone returns a clone of the user
func (u *User) Clone() *User {
	return &User{
		Username: u.Username,
		Hash:     u.Hash,
		Role:     u.Role,
	}
}
