package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Fetch artist asynchronously
func fetchArtist(apiURL string) (Artist, error) {
	var artist Artist
	resp, err := http.Get(apiURL)
	if err != nil {
		return artist, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&artist)
	if err != nil {
		return artist, err
	}
	return artist, nil
}

// Concurrently fetch multiple artists
func fetchArtists(apiURL string) (Artists, error) {
	artistsChan := make(chan Artists)
	errorChan := make(chan error)

	go func() {
		resp, err := http.Get(apiURL)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		var artists Artists
		err = json.NewDecoder(resp.Body).Decode(&artists)
		if err != nil {
			errorChan <- err
			return
		}
		artistsChan <- artists
	}()

	select {
	case artists := <-artistsChan:
		return artists, nil
	case err := <-errorChan:
		return nil, err
	case <-time.After(5 * time.Second): // Timeout
		return nil, fmt.Errorf("API request timed out")
	}
}

// Concurrently fetch relations, locations, and concert dates
func fetchExtraDetails(artist Artist) (Artist, error) {
	relationsChan := make(chan Artist)
	locationsChan := make(chan Artist)
	datesChan := make(chan Artist)
	errorChan := make(chan error)

	// Fetch relations
	go func() {
		resp, err := http.Get(artist.Relations)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&artist.Relation)
		if err != nil {
			errorChan <- err
			return
		}
		relationsChan <- artist
	}()

	// Fetch locations
	go func() {
		resp, err := http.Get(artist.Locations)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&artist.Location)
		if err != nil {
			errorChan <- err
			return
		}
		locationsChan <- artist
	}()

	// Fetch concert dates
	go func() {
		resp, err := http.Get(artist.Dates)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&artist.Date)
		if err != nil {
			errorChan <- err
			return
		}
		datesChan <- artist
	}()

	for i := 0; i < 3; i++ {
		select {
		case artist = <-relationsChan:
		case artist = <-locationsChan:
		case artist = <-datesChan:
		case err := <-errorChan:
			return artist, err
		case <-time.After(5 * time.Second):
			return artist, fmt.Errorf("Timeout fetching extra details")
		}
	}
	return artist, nil
}
