package domain

type Venue struct {
	ID       string
	Name     string
	Location Location
}

type Location struct {
	Address          string            `json:"address"`
	CrossStreet      string            `json:"crossStreet"`
	Latitude         float64           `json:"lat"`
	Longitude        float64           `json:"lng"`
	LabeledLatLngs   []LabeledLatLangs `json:"labeledLatLngs"`
	Distance         int               `json:"distance"`
	PostalCode       string            `json:"postalCode"`
	CountryCode      string            `json:"cc"`
	City             string            `json:"city"`
	State            string            `json:"state"`
	Country          string            `json:"country"`
	FormattedAddress []string          `json:"formattedAddress"`
}

// TODO(vinicius.garcia): Find a better name?
type LabeledLatLangs struct {
	Label string  `json:"label"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
}
