// Package subscription holds the data types for subscriptions and podcast episodes
package subscription

import (
	"time"

	"gorm.io/gorm/clause"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	// ErrSubNotFound is the error returned for when searching subscriptions by tag returned no results
	ErrSubNotFound = errors.New("Subscription tag was not found")
)

// Subscription is the data type for podcast subscriptions
type Subscription struct {
	gorm.Model
	Name     string
	RSSURL   string    `gorm:"unique"`
	Tag      string    `gorm:"unique"`
	Episodes []Episode `gorm:"foreignKey:SubscriptionID"`
}

// Episode contains the information about a single podcast episode
type Episode struct {
	ID             uint `gorm:"primarykey"`
	SubscriptionID int
	Title          string
	URL            string `gorm:"unique"`
	PublishedAt    time.Time
}

// AllTypes returns all the database types defined in this package
func AllTypes() []interface{} {
	return []interface{}{
		Subscription{}, Episode{},
	}
}

// SubHandler is a sub data handler, in charge of data access methods for the subscription data types
type SubHandler struct {
	db *gorm.DB
}

// Handler creates a new subscription handler from the given gorm object
func Handler(db *gorm.DB) SubHandler {
	return SubHandler{
		db: db,
	}
}

// NewSubscription adds the given subscription to the database as a new entry
func (s SubHandler) NewSubscription(sub *Subscription) error {
	if err := s.db.Create(sub).Error; err != nil {
		return errors.WithMessagef(err, "Failed saving subscription [%s/%s]", sub.Tag, sub.Name)
	}
	return nil
}

// GetSubscriptions returns an array of every subscription
func (s SubHandler) GetSubscriptions() ([]Subscription, error) {
	subs := []Subscription{}
	if err := s.db.Find(&subs).Error; err != nil {
		return subs, errors.WithMessage(err, "Failed retrieving all subscriptions")
	}
	return subs, nil
}

// GetSubscriptionByTag returns a subscription by its tag
func (s SubHandler) GetSubscriptionByTag(tag string) (Subscription, error) {
	sub := Subscription{}
	if err := s.db.Where("tag = ?", tag).First(&sub).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return sub, ErrSubNotFound
		}
		return sub, errors.WithMessagef(err, "Failed getting subscription by tag [%s]", tag)
	}
	return sub, nil
}

// GetSubscriptionsByTags returns all subscribtions with the given tags
func (s SubHandler) GetSubscriptionsByTags(tags ...string) ([]Subscription, error) {
	subs := []Subscription{}
	if err := s.db.Preload(clause.Associations).Where("tag IN (?)", tags).Find(&subs).Error; err != nil {
		return subs, errors.WithMessagef(err, "Failed getting subscriptions by tags [%v]", tags)
	}
	return subs, nil
}

// UpdateSubscription saves changes to this subscription to the database
func (s SubHandler) UpdateSubscription(sub Subscription) error {
	if err := s.db.Save(&sub).Error; err != nil {
		return errors.WithMessagef(err, "Failed updating subscription [%s/%s]", sub.Tag, sub.Name)
	}

	return nil
}

// AddEpisodes adds the given episodes to the database
func (s SubHandler) AddEpisodes(sub Subscription, eps []Episode) error {
	// Add subID to the episodes
	for index := range eps {
		eps[index].SubscriptionID = int(sub.ID)
	}
	// Save them all
	if err := s.db.Save(&eps).Error; err != nil {
		return errors.WithMessage(err, "Failed saving new episodes")
	}
	return nil
}

// RemoveEpisodes removes the given episodes from the database
func (s SubHandler) RemoveEpisodes(eps []Episode) error {
	if err := s.db.Delete(eps).Error; err != nil {
		return errors.WithMessagef(err, "Failed removing [%d] episodes", len(eps))
	}
	return nil
}

// RemoveSubscriptions removes the given subscriptions from the database
func (s SubHandler) RemoveSubscriptions(subs []Subscription) error {
	if err := s.db.Delete(subs).Error; err != nil {
		return errors.WithMessagef(err, "Failed removing [%d] subscriptions", len(subs))
	}
	return nil
}
