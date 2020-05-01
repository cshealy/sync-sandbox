package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// dao to communicate with spotify
type SpotifyDAO struct {
	BearerToken string
}

// TODO: get spotify dao

// spotify auth response
type spotifyAuth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func NewSpotifyDAO() (*SpotifyDAO, error) {
	spotifyDAO := &SpotifyDAO{}
	err := spotifyDAO.getSpotifyToken()

	if err != nil {
		return nil, err
	}

	return spotifyDAO, nil
}

// getSpotifyToken fetches a token from spotify
func (s SpotifyDAO) getSpotifyToken() error {
	url := "https://accounts.spotify.com/api/token"
	method := "POST"

	payload := strings.NewReader("grant_type=client_credentials")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Basic OGMyNzMzNmI0YWViNGU0ZWFhZmYzODNlZWE4ZTc0NGY6NmVlYWU0NDE4ZmFmNDVlODliZTU5ZDc0ODE5MzdkOTA=")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	defer res.Body.Close()
	var spotifyAuth *spotifyAuth
	err = json.NewDecoder(res.Body).Decode(&spotifyAuth)
	if err != nil {
		return err
	}
	s.BearerToken = spotifyAuth.AccessToken
	return nil
}
