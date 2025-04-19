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

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		log.Fatalf("error in login: %v\n", err)
		return err
	}

	err = s.cfg.SetUser(name)
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

func resetHandler(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		log.Printf("reset unsuccessful: %v\n", err)
		return err
	}

	log.Printf("reset successful\n")

	return nil
}

func userListHandler(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	userList, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Printf("could not fetch user list: %v\n", err)
		return err
	}

	for _, user := range userList {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
