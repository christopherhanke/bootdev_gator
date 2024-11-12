package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/christopherhanke/bootdev_gator/internal/database"
	"github.com/christopherhanke/bootdev_gator/internal/rss"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:            nextFeed.ID,
	})
	if err != nil {
		return err
	}
	rssf, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	fmt.Printf("Tried at %s\n", time.Now().String())
	for _, item := range rssf.Channel.Item {
		parseTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err
		}
		fmt.Printf("Post: %s\npublished at: %v\n\n", item.Title, parseTime)
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: sql.NullTime{Time: parseTime, Valid: true},
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			fmt.Printf("There has been an error: %s\n", err)
		}
	}
	return nil
}
