package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

func (feed *RSSFeed) String() string {
	xml, _ := xml.Marshal(feed)

	return string(xml)
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create a request for resource %s", feedURL)
	}

	req.Header.Set("User-Agent", "gator/1.0")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("request Timeout")
		}
		return nil, fmt.Errorf("failed to get response from resource %s", feedURL)
	}

	data, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("failed to read response body from resource %s", feedURL)
	}

	var feed RSSFeed

	err = xml.Unmarshal(data, &feed)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml response body from resource %s", feedURL)
	}

	unescapeResponse(&feed)

	return &feed, nil
}

func unescapeResponse(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Link = html.UnescapeString(feed.Channel.Link)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
		feed.Channel.Item[i].Link = html.UnescapeString(feed.Channel.Item[i].Link)
		feed.Channel.Item[i].PubDate = html.UnescapeString(feed.Channel.Item[i].PubDate)
	}
}
