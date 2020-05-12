package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	pb "github.com/cshealy/sync-sandbox/proto"

	log "github.com/sirupsen/logrus"

	"github.com/cshealy/sync-sandbox/data"
)

// dao to communicate with spotify
type SpotifyDAO struct {
	data.DAO
	spotifyUsername string
	spotifyPassword string
	SpotifyPlaylist string
}

// spotify auth response
type spotifyAuth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// spotify playlist response
type spotifyPlaylist struct {
	Collaborative bool   `json:"collaborative"`
	Description   string `json:"description"`
	ExternalUrls  struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  interface{} `json:"href"`
		Total int         `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
	Name  string `json:"name"`
	Owner struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"owner"`
	Public     interface{} `json:"public"`
	SnapshotID string      `json:"snapshot_id"`
	Tracks     struct {
		Href  string `json:"href"`
		Items []struct {
			AddedAt time.Time `json:"added_at"`
			AddedBy struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"added_by"`
			IsLocal bool `json:"is_local"`
			Track   struct {
				Album struct {
					AlbumType        string   `json:"album_type"`
					AvailableMarkets []string `json:"available_markets"`
					ExternalUrls     struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href   string `json:"href"`
					ID     string `json:"id"`
					Images []struct {
						Height int    `json:"height"`
						URL    string `json:"url"`
						Width  int    `json:"width"`
					} `json:"images"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"album"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
				AvailableMarkets []string `json:"available_markets"`
				DiscNumber       int      `json:"disc_number"`
				DurationMs       int      `json:"duration_ms"`
				Explicit         bool     `json:"explicit"`
				ExternalIds      struct {
					Isrc string `json:"isrc"`
				} `json:"external_ids"`
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href        string `json:"href"`
				ID          string `json:"id"`
				Name        string `json:"name"`
				Popularity  int    `json:"popularity"`
				PreviewURL  string `json:"preview_url"`
				TrackNumber int    `json:"track_number"`
				Type        string `json:"type"`
				URI         string `json:"uri"`
			} `json:"track"`
		} `json:"items"`
		Limit    int         `json:"limit"`
		Next     string      `json:"next"`
		Offset   int         `json:"offset"`
		Previous interface{} `json:"previous"`
		Total    int         `json:"total"`
	} `json:"tracks"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}

// NewDAO creates a new spotify data access object for interacting with spotify
func NewDAO(spotifyUsername, spotifyPassword, spotifyPlaylistId string) (*SpotifyDAO, error) {
	spotifyDAO := &SpotifyDAO{
		spotifyUsername: spotifyUsername,
		spotifyPassword: spotifyPassword,
		SpotifyPlaylist: spotifyPlaylistId,
	}
	err := spotifyDAO.getBearerToken()

	if err != nil {
		return nil, err
	}

	return spotifyDAO, nil
}

// GetSpotifyToken fetches a token from spotify
func (s *SpotifyDAO) getBearerToken() error {
	url := "https://accounts.spotify.com/api/token"
	method := "POST"

	payload := strings.NewReader("grant_type=client_credentials")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}

	req.SetBasicAuth(s.spotifyUsername, s.spotifyPassword)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	defer res.Body.Close()
	var spotifyAuth *spotifyAuth
	err = json.NewDecoder(res.Body).Decode(&spotifyAuth)
	if err != nil {
		return err
	}

	// determine if the response was successful
	if res.StatusCode >= 200 || res.StatusCode <= 299 {
		s.BearerToken = spotifyAuth.AccessToken

		// refresh the bearer token if the current utc time is after token expiration
		s.TokenExpiration = time.Now().UTC().Add(time.Minute * time.Duration(spotifyAuth.ExpiresIn))
	} else {
		// TODO: add more details to error response
		return errors.New(fmt.Sprintf("failed to get spotify bearer token -- status_code: %d", res.StatusCode))
	}

	return nil
}

// GetPlaylist will get the playlist metadata passed in through docker-compose
func (s *SpotifyDAO) GetPlaylist() (*pb.SpotifyPlaylist, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", s.SpotifyPlaylist)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.BearerToken))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	log.Info(res)
	log.Info(res.StatusCode)
	defer res.Body.Close()
	var spotifyPlaylist *spotifyPlaylist
	err = json.NewDecoder(res.Body).Decode(&spotifyPlaylist)
	if err != nil {
		return nil, err
	}

	// determine if the response was successful
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("failed to get playlist -- status_code: %d", res.StatusCode))
	}

	// parse the spotify playlist response
	var spotifyPlaylistTracks []*pb.SpotifyPlaylistTrack
	for _, track := range spotifyPlaylist.Tracks.Items {

		// collect all of the artists for this track
		var trackArtists []*pb.SpotifyPlaylistArtist
		for _, artist := range track.Track.Artists {
			trackArtists = append(trackArtists, &pb.SpotifyPlaylistArtist{
				Name: artist.Name,
			})
		}

		// add a new track to the playlist response
		spotifyPlaylistTracks = append(spotifyPlaylistTracks, &pb.SpotifyPlaylistTrack{
			Name:    track.Track.Name,
			Artists: trackArtists,
		})

	}

	return &pb.SpotifyPlaylist{
		Tracks: spotifyPlaylistTracks,
	}, nil
}

// refreshToken checks if the current time is after spotify's expiration time,
// then fetches a new bearer token if needed
func (s *SpotifyDAO) refreshToken() {
	if s.TokenExpiration.Before(time.Now().UTC()) {
		s.getBearerToken()
	}
}
