package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// helper function making an HTTP request, decoding JSON, handling timeouts
func fetchData[T any](apiURL string) (T, error) {
	var data T
	dataChan := make(chan T)
	errorChan := make(chan error)

	go func() {
		resp, err := http.Get(apiURL)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			errorChan <- err
			return
		}
		dataChan <- data
	}()

	select {
	case data := <-dataChan:
		return data, nil
	case err := <-errorChan:
		return data, err
	case <-time.After(5 * time.Second):
		return data, fmt.Errorf("API request timed out")
	}
}

// // Fetch artist asynchronously
// func fetchArtist(apiURL string) (Artist, error) {
// 	var artist Artist
// 	resp, err := http.Get(apiURL)
// 	if err != nil {
// 		return artist, err
// 	}
// 	defer resp.Body.Close()

// 	err = json.NewDecoder(resp.Body).Decode(&artist)
// 	if err != nil {
// 		return artist, err
// 	}
// 	return artist, nil
// }

// // Concurrently fetch multiple artists
// func FetchArtists(apiURL string) (Artists, error) {
// 	var artists Artists
// 	client := http.Client{Timeout: 5 * time.Second} // Προσθήκη timeout για ασφάλεια
// 	resp, err := client.Get(apiURL)
// 	if err != nil {
// 		return artists, err
// 	}
// 	defer resp.Body.Close()

// 	err = json.NewDecoder(resp.Body).Decode(&artists)
// 	if err != nil {
// 		return artists, fmt.Errorf("Error decoding API response: %v", err)
// 	}

// 	return artists, nil
// }

// Fetch artist data from given API URL
func fetchArtist(apiURL string) (Artist, error) {

	// var artist Artist

	// // get request to fetch artist data
	// resp, err := http.Get(apiURL)
	// if err != nil {
	// 	return artist, err
	// }
	// defer resp.Body.Close() // ensure response body is closed after function exits

	// // decode JSON response into the artist struct
	// err = json.NewDecoder(resp.Body).Decode(&artist)
	// if err != nil {
	// 	return artist, err
	// }
	// return artist, nil
	return fetchData[Artist](apiURL)
}

// Concurrently fetch multiple artists
func FetchArtists(apiURL string) (Artists, error) {

	// artistsChan := make(chan Artists)
	// errorChan := make(chan error)

	// // go routines are used to fetch data asynchronously and avoid blocking execution
	// go func() {
	// 	// send get request to fetch artists' data
	// 	resp, err := http.Get(apiURL)
	// 	if err != nil {
	// 		errorChan <- err
	// 		return
	// 	}
	// 	defer resp.Body.Close()

	// 	var artists Artists
	// 	// decode JSON response into the artists struct
	// 	err = json.NewDecoder(resp.Body).Decode(&artists)
	// 	if err != nil {
	// 		errorChan <- err
	// 		return
	// 	}
	// 	artistsChan <- artists
	// }()

	// // select is used to handle multiple possible responses:
	// select {
	// case artists := <-artistsChan:
	// 	return artists, nil
	// case err := <-errorChan:
	// 	return nil, err
	// case <-time.After(5 * time.Second): // Timeout, preventing indefinite waiting.
	// 	return nil, fmt.Errorf("API request timed out")
	// }
	return fetchData[Artists](apiURL)
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
