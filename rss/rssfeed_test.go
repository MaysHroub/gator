package rss

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

	rssFeed, err := FetchFeed(context.Background(), svr.URL)

	require.NoError(t, err)
	assert.Equal(t, feedTitle, rssFeed.Channel.Title)
	assert.Equal(t, feedLink, rssFeed.Channel.Link)
	assert.Equal(t, feedDesc, rssFeed.Channel.Description)
	assert.Equal(t, 2, len(rssFeed.Channel.Items))
}