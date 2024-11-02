package main

import (
	"context"
	"fmt"

	"github.com/christopherhanke/bootdev_gator/internal/database"
)

func middlerwareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("no user logged in: %v", err)
		}
		return handler(s, cmd, user)
	}
}
