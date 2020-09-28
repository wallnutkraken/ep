// Package poddata is the DAL for ep. It's responsible for persistent data storage. Information such as
// the podcasts added, the episodes currently synced, et cetera.
package poddata

import (
	"github.com/pkg/errors"
	"github.com/wallnutkraken/ep/poddata/subscription"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Data is the wrapper struct for our GORM object, and contains methods to interact with the persistent storage.
type Data struct {
	db *gorm.DB
}

// New creates a new instance of the Data access object
func New(confPath string) (Data, error) {
	// Open gorm with the given data path
	db, err := gorm.Open(sqlite.Open(confPath), &gorm.Config{})
	if err != nil {
		return Data{}, errors.WithMessagef(err, "Failed opening sqlite file at [%s]", confPath)
	}
	// Wrap the gorm object in our Data object
	data := Data{
		db: db,
	}
	// And call migrate to automigrate
	err = data.migrate()
	if err != nil {
		return data, errors.WithMessage(err, "Failed migrating database")
	}
	return data, nil
}

// migrate collects all the database data types and calls gorm's AutoMigrate method
// to migrate the schema to the database
func (d Data) migrate() error {
	allDataTypes := []interface{}{}
	allDataTypes = append(allDataTypes, subscription.AllTypes()...)

	return d.db.AutoMigrate(allDataTypes...)
}

// Subscriptions returns the sub-handler for subscriptions which contains subscription-related methods
func (d Data) Subscriptions() subscription.SubHandler {
	return subscription.Handler(d.db)
}
