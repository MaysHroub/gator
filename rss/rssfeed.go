// Package rss provides structures and functions for parsing and representing RSS feeds.
package rss

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github/MaysHroub/gator/internal/database"
	"github/MaysHroub/gator/internal/gatorapi"
	"github/MaysHroub/gator/internal/repository"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
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

	err = savePosts(rssFeed, db, feed.ID)

	return err 
}

func savePosts(rssFeed *RSSFeed, db repository.Repository, feedID uuid.UUID) error {
	ctx := context.Background()
	for _, item := range rssFeed.Channel.Items {
		parsedPubDate, err := parsePublishDate(item.PubDate)
		if err != nil {
			return err
		}
		createPostParams := database.CreatePostParams {
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Description: sql.NullString{String: item.Description, Valid: true},
			Url: item.Link,
			PublishedAt: sql.NullTime{Time: parsedPubDate, Valid: !parsedPubDate.IsZero()},
			FeedID: feedID,
		}
		_, err = db.CreatePost(ctx, createPostParams)
		if err == nil {
			continue
		}
		if strings.Contains(err.Error(), "unique constraint") && strings.Contains(err.Error(), "posts_url_key") {
			continue // ignore the error
		}
		return err
	}
	return nil 
}

func parsePublishDate(pubDate string) (time.Time, error) {
	layout := "Mon, 02 Jan 2006 15:04:05 MST"
	parsedPubDate, err := time.Parse(layout, pubDate)
	return parsedPubDate, err
}