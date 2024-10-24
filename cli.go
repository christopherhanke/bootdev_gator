package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"unicode"

	"github.com/christopherhanke/bootdev_gator/internal/config"
	"github.com/christopherhanke/bootdev_gator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	//Login user given in commands args slice
	if len(cmd.args) < 1 {
		return fmt.Errorf("commands arg slice is smaller than 1")
	}
	if len(cmd.args) > 1 {
		return fmt.Errorf("commands arg slice is bigger than 1")
	}
	_, err := s.db.GetUsers(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("user is not registered: %v", cmd.args[0])
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set: %v\n", s.cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	//Register user given in commands args slice
	if len(cmd.args) < 1 {
		return fmt.Errorf("commands arg slice is smaller than 1")
	}
	if len(cmd.args) > 1 {
		return fmt.Errorf("commands arg slice is bigger than 1")
	}
	if len(cmd.args[0]) <= 1 {
		return fmt.Errorf("no valid username given: %v (to short)", cmd.args[0])
	}
	if !isName(cmd.args[0]) {
		return fmt.Errorf("no valid username given: %v (no char)", cmd.args[0])
	}
	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			Name:      cmd.args[0],
		},
	)
	if err != nil {
		return fmt.Errorf("user already exists. couldn't create: %v", cmd.args[0])
	}
	s.cfg.SetUser(user.Name)
	fmt.Println(user.ID, user.CreatedAt.Time, user.UpdatedAt.Time, user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	//Reset database state, deleting all registered users
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete table users: %v", err)
	}
	return nil
}

type commands struct {
	CommandMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	//registers a new handler for a command
	_, exists := c.CommandMap[name]
	if exists {
		fmt.Printf("Handler exists: %v", name)
		return
	}
	c.CommandMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	//run a given command with provided state if exists
	_, exists := c.CommandMap[cmd.name]
	fmt.Printf("run: %v\n", cmd.name)
	if !exists {
		return fmt.Errorf("command does not exist: %v", cmd.name)
	}

	f := c.CommandMap[cmd.name]
	return f(s, cmd)
}

func isName(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
