package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/christopherhanke/bootdev_gator/internal/database"
	"github.com/christopherhanke/bootdev_gator/internal/rss"
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
		fmt.Printf("%s\n", item.Title)
	}
	return nil
}
