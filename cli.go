package main

import (
	"fmt"

	"github.com/christopherhanke/bootdev_gator/internal/config"
)

type state struct {
	cfg *config.Config
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
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set: %v\n", s.cfg.CurrentUserName)
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
