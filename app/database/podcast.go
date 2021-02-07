package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Podcast represents an individual podcast added by a user.
type Podcast struct {
	gorm.Model
	FeedURL       string
	Title         string
	LastRefreshed *time.Time
	OwnerID       uint `gorm:"index"`
	Owner         User
}

// Refresh refreshes the feed for a given podcast.
func (podcast *Podcast) Refresh() error {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(podcast.FeedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	for _, item := range feed.Items {
		_, err = podcast.AddEpisode(item)
		if err != nil {
			return fmt.Errorf("error adding new episode: %w", err)
		}
	}

	refreshTime := time.Now().UTC()
	podcast.LastRefreshed = &refreshTime
	Instance.Save(podcast)
	return nil
}

// AddEpisode adds a new Episode instance to a podcast from a feed item.
func (podcast *Podcast) AddEpisode(item *gofeed.Item) (*Episode, error) {
	if len(item.Enclosures) != 1 {
		return nil, errors.New("item has unexpected number of enclosures")
	}
	instance := &Episode{
		Title:         item.Title,
		Description:   item.Description,
		PublishedTime: *item.PublishedParsed,
		GUID:          item.GUID,
		URL:           item.Link,
		AudioURL:      item.Enclosures[0].URL,
		PodcastID:     podcast.ID,
	}

	tx := Instance.Create(instance)
	if tx.Error != nil {
		return nil, fmt.Errorf("error inserting episode into database: %w", tx.Error)
	}

	log.Info().Str("title", instance.Title).Uint("id", instance.ID).Uint("podcast_id", instance.PodcastID).Msg("New episode created")

	return instance, nil
}

// NewPodcast creates a new podcast.
func NewPodcast(ownerID uint, feedURL string) (*Podcast, error) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(feedURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching podcast details: %w", err)
	}
	instance := &Podcast{
		FeedURL: feedURL,
		Title:   feed.Title,
		OwnerID: ownerID,
	}

	tx := Instance.Create(instance)
	if tx.Error != nil {
		return nil, fmt.Errorf("error inserting podcast into database: %w", tx.Error)
	}

	log.Info().Str("title", instance.Title).Uint("id", instance.ID).Str("feed_url", instance.FeedURL).Msg("New podcast created")

	return instance, nil
}

// Episode represents an individual episode of a podcast.
type Episode struct {
	gorm.Model
	Title         string
	Description   string
	PublishedTime time.Time
	GUID          string `gorm:"index"`
	URL           string
	AudioURL      string
	PodcastID     uint `gorm:"index"`
	Podcast       Podcast
}
