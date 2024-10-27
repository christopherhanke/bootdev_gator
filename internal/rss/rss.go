package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

const URL = "https://www.wagslane.dev/index.xml"

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	//fetch RSS feed from URL and return as structured data
	req, err := http.NewRequestWithContext(context.Background(), "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "bootdev_gator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssf RSSFeed
	err = xml.Unmarshal(data, &rssf)
	if err != nil {
		return nil, err
	}

	rssf.Channel.Description = html.UnescapeString(rssf.Channel.Description)
	rssf.Channel.Title = html.UnescapeString(rssf.Channel.Title)
	for key, item := range rssf.Channel.Item {
		rssf.Channel.Item[key].Title = html.UnescapeString(item.Title)
		rssf.Channel.Item[key].Description = html.UnescapeString(item.Description)
	}

	return &rssf, nil
}
