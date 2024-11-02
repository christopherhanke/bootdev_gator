package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"unicode"

	"github.com/christopherhanke/bootdev_gator/internal/config"
	"github.com/christopherhanke/bootdev_gator/internal/database"
	"github.com/christopherhanke/bootdev_gator/internal/rss"
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
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
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

func handlerUsers(s *state, cmd command) error {
	//Print list of registered users and highlight current users
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users from database: %v", err)
	}
	if len(users) < 1 {
		fmt.Println("no registered users")
		return nil
	}
	for _, user := range users {
		fmt.Printf("* %s ", user)
		if user == s.cfg.CurrentUserName {
			fmt.Print("(current)")
		}
		fmt.Println()
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	rssf, err := rss.FetchFeed(context.Background(), rss.URL)
	if err != nil {
		return err
	}

	fmt.Println("---- RSS Feed ----")
	fmt.Printf("Title: %s\n", rssf.Channel.Title)
	fmt.Printf("Link: %s\n", rssf.Channel.Link)
	fmt.Printf("Description: %s\n", rssf.Channel.Description)
	fmt.Println()
	for key, item := range rssf.Channel.Item {
		fmt.Printf("--- Item %2d ---\n", key)
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Description: %s\n", item.Description)
		fmt.Printf("Date: %s\n", item.PubDate)
		fmt.Println()
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	//prints all feeds in database
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Name: %s, URL: %s, User: %s\n", feed.Name, feed.Url, feed.Name_2)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	//add feed to feeds table in database
	if len(cmd.args) < 2 {
		return fmt.Errorf("commands args is smaller than 2")
	}
	if len(cmd.args) > 2 {
		return fmt.Errorf("commands args is bigger than 2")
	}

	dbUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %s", s.cfg.CurrentUserName)
	}

	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			Name:      cmd.args[0],
			Url:       cmd.args[1],
			UserID:    dbUser.ID,
		},
	)
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			UserID:    dbUser.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return err
	}
	fmt.Println(feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	//takes a single URL and creates a new feed follow for current user
	if len(cmd.args) < 1 {
		return fmt.Errorf("commands arg slice is smaller than 1")
	}
	if len(cmd.args) > 1 {
		return fmt.Errorf("commands arg slice is bigger than 1")
	}
	dbUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user: %s, %v", s.cfg.CurrentUserName, err)
	}
	dbFeed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to get feed: %s, %v", cmd.args[0], err)
	}
	fmt.Printf("Feedname: %s, URL: %s\n", dbFeed.Name, dbFeed.Url)
	feedFollow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			UserID:    dbUser.ID,
			FeedID:    dbFeed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to follow feed: %s, %v", cmd.args[0], err)
	}
	fmt.Printf("Follow feed: %s for user: %s\n", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	//print the feeds current user is following
	dbUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user: %s, %v", s.cfg.CurrentUserName, err)
	}
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), dbUser.ID)
	if err != nil {
		return fmt.Errorf("failed to get feeds for user: %s, %v", s.cfg.CurrentUserName, err)
	}
	for _, feed := range feedFollows {
		fmt.Printf("- %s\n", feed.FeedName)
	}
	return nil
}

type commands struct {
	CommandMap map[string]func(*state, command) error
}

func (c *commands) register(name string, handler func(*state, command) error) {
	//registers a new handler for a command
	_, exists := c.CommandMap[name]
	if exists {
		fmt.Printf("Handler exists: %v", name)
		return
	}
	c.CommandMap[name] = handler
}

func (c *commands) run(s *state, cmd command) error {
	//run a given command with provided state if exists
	handler, exists := c.CommandMap[cmd.name]
	if !exists {
		return fmt.Errorf("command does not exist: %v", cmd.name)
	}
	return handler(s, cmd)
}

func isName(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
