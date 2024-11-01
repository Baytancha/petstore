package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	_ "github.com/lib/pq"
)

var Pets []*Pet

var PrimaryKeyIDx map[int64]*Pet

type Pet struct {

	// category
	Category *Category `json:"category,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// name
	// Example: doggie
	// Required: true
	Name *string `json:"name"`

	// photo urls
	// Required: true
	PhotoUrls []string `json:"photoUrls" xml:"photoUrls"`

	// pet status in the store
	// Enum: ["available","pending","sold"]
	Status string `json:"status,omitempty"`

	// tags
	Tags []*Tag `json:"tags" xml:"tags"`
}

func (c Category) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (t Tag) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (c *Category) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}

func (c *Tag) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}
