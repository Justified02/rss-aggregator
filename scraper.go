package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Justified02/rssagg/internal/database"
)

// 1. Define RSS structs
type RSSFeed struct {
	Channel struct {
		Title string    `xml:"title"`
		Items []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// 2. fetchFeed fetches a URL, reads it and parses XML (converts it to Go structs)
func fetchFeed(ctx context.Context, feedURL string) (RSSFeed, error) {
	var feed RSSFeed

	// 1. Create HTTP request with context - create an internet request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil) // Like opening chrome and typing a URL
	if err != nil {
		return feed, err
	}

	// 2. Send request
	resp, err := http.DefaultClient.Do(req) // send the request to the internet
	if err != nil {
		return feed, err
	}
	defer resp.Body.Close() // Always close the response

	// 3. Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return feed, err
	}

	// 4. Parse XML
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return feed, nil
	}

	return feed, nil
}

func startScraping(db *database.Queries, concurrency int, interval time.Duration) {
	ticker := time.NewTicker(interval) // fires every `interval` (e.g., 1 minute)
	defer ticker.Stop()

	for {
		// 1. Wait for the next tick
		<-ticker.C

		// 2. Get feeds to fetch from DB
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("error getting feeds: %v", err)
			continue
		}

		var wg sync.WaitGroup

		// 3. Launch a goroutine for each feed
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, &wg, feed)
		}

		// 4. Wait for all goroutines to finish
		wg.Wait()
	}
}

// 3. Scrape a single feed
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done() // tells WaitGroup this goroutine is done

	// Mark feed as fetched
	_, err := db.MarkFeedsToFetch(context.Background(), feed.ID)
	if err != nil {
		log.Printf("failed to mark feed fetched: %v", err)
		return
	}

	// Fetch RSS
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("failed to fetch feed %s: %v", feed.Url, err)
		return
	}

	// Loop over items
	for _, item := range rssFeed.Channel.Items {
		pubTime, err := time.Parse(time.RFC1123, item.PubDate)

		// Quick safety check: if parsing fails, we pass a null time
    	isTimeValid := err == nil 

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			FeedID:      feed.ID,
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: sql.NullTime{Time: pubTime, Valid: isTimeValid}, // helper function
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("failed to create post: %v", err)
		}
	}
}
