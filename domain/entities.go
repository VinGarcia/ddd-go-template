package domain

import "time"

type Venue struct {
	ID       string
	Name     string
	Location Location
}

type Location struct {
	Address          string          `json:"address"`
	CrossStreet      string          `json:"crossStreet"`
	Latitude         float64         `json:"lat"`
	Longitude        float64         `json:"lng"`
	LabeledLatLngs   []LabeledCoords `json:"labeledLatLngs"`
	Distance         int             `json:"distance"`
	PostalCode       string          `json:"postalCode"`
	CountryCode      string          `json:"cc"`
	City             string          `json:"city"`
	State            string          `json:"state"`
	Country          string          `json:"country"`
	FormattedAddress []string        `json:"formattedAddress"`
}

type LabeledCoords struct {
	Label string  `json:"label"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
}

type User struct {
	ID   int    `ksql:"id"`
	Name string `ksql:"name"`
	Age  int    `ksql:"age"`

	CreatedAt *time.Time `ksql:"created_at"`
	UpdatedAt *time.Time `ksql:"updated_at"`
}
