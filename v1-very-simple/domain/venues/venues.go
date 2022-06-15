package venues

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vingarcia/ddd-go-template/v1-very-simple/domain"
)

type Service struct {
	logger domain.LogProvider
	rest   domain.RestProvider
	cache  domain.CacheProvider

	baseURL  string
	clientID string
	secret   string
}

func NewService(
	logger domain.LogProvider,
	rest domain.RestProvider,
	cache domain.CacheProvider,
	baseURL string,
	clientID string,
	secret string,
) Service {
	return Service{
		logger:   logger,
		cache:    cache,
		rest:     rest,
		baseURL:  baseURL,
		clientID: clientID,
		secret:   secret,
	}
}

func (s Service) GetVenues(ctx context.Context, latitude string, longitude string) ([]domain.Venue, error) {
	url := fmt.Sprintf("%s/venues/search?client_id=%s&client_secret=%s&v=20210514&ll=%s,%s", s.baseURL, s.clientID, s.secret, latitude, longitude)
	resp, err := s.rest.Get(ctx, url, domain.RequestData{})
	if err != nil {
		s.logger.Error(ctx, "error-retrieving-venues-from-foursquare-by-coordinates", domain.LogBody{
			"latitude":  latitude,
			"longitude": longitude,
			"payload":   string(resp.Body),
		})
		return nil, domain.InternalErr("error-retrieving-venues-from-foursquare", map[string]interface{}{
			"latitude":  latitude,
			"longitude": longitude,
		})
	}

	var respBody struct {
		Meta struct {
			Code      int    `json:"code"`
			RequestID string `json:"requestId"`
		}
		Response struct {
			Venues []domain.Venue `json:"venues"`
		} `json:"response"`
	}
	err = json.Unmarshal(resp.Body, &respBody)
	if err != nil {
		return nil, fmt.Errorf("error parsing foursquare venues as JSON: %s", err)
	}

	return respBody.Response.Venues, nil
}

func (s Service) GetVenue(ctx context.Context, venueID string) ([]byte, error) {
	var cachedVenue []byte
	err := s.cache.Get(ctx, venueID, &cachedVenue)
	if err == nil {
		// Log IDs, not payloads whenever possible, except when errors happen, then log everything.
		s.logger.Debug(ctx, "fetching-venue-from-cache", domain.LogBody{
			"venue_id": venueID,
		})
		return cachedVenue, nil
	}

	url := fmt.Sprintf("%s/venues/%s?client_id=%s&client_secret=%s&v=20210514", s.baseURL, venueID, s.clientID, s.secret)
	resp, err := s.rest.Get(ctx, url, domain.RequestData{})
	if err != nil {
		s.logger.Error(ctx, "error-fetching-venue-by-latitude-from-foursquare", domain.LogBody{
			"venue_id": venueID,
			"error":    err.Error(),
			"payload":  string(resp.Body),
		})
		return nil, err
	}

	s.logger.Debug(ctx, "adding-venue-to-cache", domain.LogBody{
		"venue_id": venueID,
	})
	err = s.cache.Set(ctx, venueID, resp.Body)
	return resp.Body, err
}
