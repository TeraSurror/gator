package main

import (
	"context"
	"fmt"

	"github.com/TeraSurror/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("user %s not logged in", s.cfg.CurrentUserName)
		}
		return handler(s, cmd, user)
	}
}
