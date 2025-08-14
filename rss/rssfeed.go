// Package rss provides structures and functions for parsing and representing RSS feeds.
package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
)

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
	client := gatorapi.NewClient(feedURL)
	resp, err := client.MakeRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to make request to %s: %w", feedURL, err)
	}
	defer resp.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rssFeed RSSFeed
	if err = xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, fmt.Errorf("failed to parse response data: %w", err)
	}

	return &rssFeed, nil
}