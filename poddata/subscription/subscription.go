// Package subscription holds the data types for subscriptions and podcast episodes
package subscription

import (
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
	RSSURL   string
	Tag      string    `gorm:"unique"`
	Episodes []Episode `gorm:"foreignKey:SubscriptionID"`
}

// Episode contains the information about a single podcast episode
type Episode struct {
	ID             uint `gorm:"primarykey"`
	SubscriptionID int
	Title          string
	URL            string
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
