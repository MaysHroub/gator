// Package rss provides structures and functions for parsing and representing RSS feeds.
package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"github/MaysHroub/gator/internal/gatorapi"
	"github/MaysHroub/gator/internal/repository"
	"html"
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

func FetchFeed(feedURL string) (*RSSFeed, error) {
	client := gatorapi.NewClient(timeout)
	data, err := client.Get(feedURL)
	if err != nil {
		return nil, err
	}
	
	var rssFeed RSSFeed
	if err = xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, fmt.Errorf("failed to parse response data: %w", err)
	}

	unescapeRssFeedFields(&rssFeed)

	return &rssFeed, nil
}

func unescapeRssFeedFields(rssFeed *RSSFeed) {
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for i := range rssFeed.Channel.Items {
		rssFeed.Channel.Items[i].Title = html.UnescapeString(rssFeed.Channel.Items[i].Title)
		rssFeed.Channel.Items[i].Description = html.UnescapeString(rssFeed.Channel.Items[i].Description)
	}
}

func ScrapeFeeds(db repository.Repository) error {
	feed, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	db.MarkFeedFetched(context.Background(), feed.ID)

	rssFeed, err := FetchFeed(feed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("Items of %s with link %s:\n", rssFeed.Channel.Title, rssFeed.Channel.Link)

	for _, item := range rssFeed.Channel.Items {
		fmt.Println(item.Title)
	}

	return nil 
}