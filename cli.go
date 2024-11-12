package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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

// Login user given in commands args slice
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 || len(cmd.args) > 1 {
		return fmt.Errorf("usage: %v <username>", cmd.name)
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

// Register user given in commands args slice
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 || len(cmd.args) > 1 {
		return fmt.Errorf("usage: %v <username>", cmd.name)
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

// Reset database state, deleting all registered users
func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete table users: %v", err)
	}
	return nil
}

// Print list of registered users and highlight current users
func handlerUsers(s *state, cmd command) error {
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

// aggregate feeds, needs one arg of time between fetching feeds, e.g. "1m" for one minute
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 || len(cmd.args) > 1 {
		return fmt.Errorf("usage: %v <time between fetches, e.g. 1m>", cmd.name)
	}
	time_between_req, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds everey %v\n", time_between_req)
	ticker := time.NewTicker(time_between_req)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

// prints all feeds in database
func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Name: %s, URL: %s, User: %s\n", feed.Name, feed.Url, feed.Name_2)
	}
	return nil
}

// add feed to feeds table in database
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 || len(cmd.args) > 2 {
		return fmt.Errorf("usage: %v <name> <url>", cmd.name)
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

// takes a single URL and creates a new feed follow for current user
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 || len(cmd.args) > 1 {
		return fmt.Errorf("usage: %v <url>", cmd.name)
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

// print the feeds current user is following
func handlerFollowing(s *state, cmd command, user database.User) error {
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

// unfollow given feed by url for logged in user
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 || len(cmd.args) > 1 {
		return fmt.Errorf("usage: %v <url>", cmd.name)
	}
	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		Name: user.Name,
		Url:  cmd.args[0],
	})
	if err != nil {
		return err
	}
	return nil
}

// browse the feeds from logged in user, argument for number of posts [default=2]
func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error
	if len(cmd.args) == 1 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("usage: %v <Optional: limit number of posts>", cmd.name)
		}
	} else {
		limit = 2
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		Name:  user.Name,
		Limit: int32(limit),
	})
	if err != nil {
		return err
	}
	for num, post := range posts {
		fmt.Printf("%2v - %v\n", num, post.Title)
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
