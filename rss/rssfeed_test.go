package rss

import (
	"database/sql"
	"fmt"
	"github/MaysHroub/gator/internal/database"
	"github/MaysHroub/gator/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFetchFeed_ValidFetch(t *testing.T) {
	feedTitle := "RSS Feed Example"
	feedLink := "https://www.example.com"
	feedDesc := "This is an example RSS feed"
	dummyFeedResponse := fmt.Sprintf(`
		<rss xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
		<channel>
		<title>%s</title>
		<link>%s</link>
		<description>%s</description>
		<item>
			<title>First Article</title>
			<link>https://www.example.com/article1</link>
			<description>This is the content of the first article.</description>
			<pubDate>Mon, 06 Sep 2021 12:00:00 GMT</pubDate>
		</item>
		<item>
			<title>Second Article</title>
			<link>https://www.example.com/article2</link>
			<description>Here's the content of the second article.</description>
			<pubDate>Tue, 07 Sep 2021 14:30:00 GMT</pubDate>
		</item>
		</channel>
		</rss>`, feedTitle, feedLink, feedDesc)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, dummyFeedResponse)
	}))
	defer svr.Close()

	rssFeed, err := FetchFeed(svr.URL)

	require.NoError(t, err)
	assert.Equal(t, feedTitle, rssFeed.Channel.Title)
	assert.Equal(t, feedLink, rssFeed.Channel.Link)
	assert.Equal(t, feedDesc, rssFeed.Channel.Description)
	assert.Equal(t, 2, len(rssFeed.Channel.Items))
}

func TestScrapeFeeds(t *testing.T) {
	dummyFeedResponse := `
		<rss xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
		<channel>
		<title>RSS Feed Example</title>
		<link>https://www.example.com</link>
		<description>This is an example RSS feed</description>
		<item>
			<title>First Article</title>
			<link>https://www.example.com/article1</link>
			<description>This is the content of the first article.</description>
			<pubDate>Mon, 06 Sep 2021 12:00:00 GMT</pubDate>
		</item>
		<item>
			<title>Second Article</title>
			<link>https://www.example.com/article2</link>
			<description>Here's the content of the second article.</description>
			<pubDate>Tue, 07 Sep 2021 14:30:00 GMT</pubDate>
		</item>
		</channel>
		</rss>`

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, dummyFeedResponse)
	}))
	defer svr.Close()

	feedID := uuid.New()
	mockDB := repository.MockRepository{}
	mockDB.On("GetNextFeedToFetch", mock.Anything).Return(database.GetNextFeedToFetchRow{
		ID: feedID,
		Url: svr.URL,
	}, nil)
	mockDB.On("MarkFeedFetched", mock.Anything, feedID).Return(nil)

	err := ScrapeFeeds(&mockDB)
	require.NoError(t, err) 

	mockDB.AssertCalled(t, "GetNextFeedToFetch", mock.Anything)
	mockDB.AssertCalled(t, "MarkFeedFetched", mock.Anything, feedID)
}

func TestSavePosts_UniquePostsURL(t *testing.T) {
	dummyRSSFeed := &RSSFeed{
		Channel: struct {
			Title       string    `xml:"title"`
			Link        string    `xml:"link"`
			Description string    `xml:"description"`
			Items       []RSSItem `xml:"item"`
		}{
			Title:       "RSS Feed Example",
			Link:        "https://www.example.com",
			Description: "This is an example RSS feed",
			Items: []RSSItem{
				{
					Title:       "First Article",
					Link:        "https://www.example.com/article1",
					Description: "This is the content of the first article.",
					PubDate:     "Mon, 06 Sep 2021 12:00:00 GMT",
				},
			},
		},
	}

	mockDB := repository.MockRepository{}
	feedID := uuid.New()
	pubDateStr := "Mon, 06 Sep 2021 12:00:00 GMT"
	layout := "Mon, 02 Jan 2006 15:04:05 MST"
	parsedTime, _ := time.Parse(layout, pubDateStr)

	paramsMatcher := mock.MatchedBy(func(p database.CreatePostParams) bool {
		return p.Title == "First Article" &&
				p.Description == sql.NullString{String: "This is the content of the first article.", Valid: true} &&
				p.Url == "https://www.example.com/article1" &&
				p.PublishedAt.Time.Equal(parsedTime) &&
				p.FeedID == feedID
	})
	mockDB.On("CreatePost", mock.Anything, paramsMatcher).Return(database.Post{}, nil)

	err := savePosts(dummyRSSFeed, &mockDB, feedID)
	require.NoError(t, err)

	mockDB.AssertCalled(t, "CreatePost", mock.Anything, paramsMatcher)
}