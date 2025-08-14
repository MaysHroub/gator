// Package rss provides structures and functions for parsing and representing RSS feeds.
package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"github/MaysHroub/gator/internal/gatorapi"
	"time"
)

const timeout = 10 * time.Second


type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := gatorapi.NewClient(timeout)
	data, err := client.Get(feedURL)
	if err != nil {
		return nil, err
	}
	
	var rssFeed RSSFeed
	if err = xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, fmt.Errorf("failed to parse response data: %w", err)
	}

	return &rssFeed, nil
}