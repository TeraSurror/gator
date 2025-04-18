package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TeraSurror/gator/internal/database"
	"github.com/google/uuid"
)

func loginHandler(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	err := s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("could not set user: %w", err)
	}

	log.Printf("User %s has been logged in\n", s.cfg.CurrentUserName)
	return nil
}

func registerHandler(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	user := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		log.Printf("user creation failed: %v\n", err)
		return err
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("could not set user: %w", err)
	}

	log.Printf("User has been created: %+v\n", createdUser)
	return nil
}
