package service

import (
	"encoding/json"
	"fmt"
	"github.com/aaanger/music-library/internal/dto"
	"net/http"
	"net/url"
)

func (s *MusicService) fetchSongFromAPI(group, song string) (*dto.SongDetail, error) {
	urlAPI := fmt.Sprintf("%s/info?group=%s&song=%s", s.apiURL, url.QueryEscape(group), url.QueryEscape(song))

	res, err := http.Get(urlAPI)
	if err != nil {
		s.log.Errorf("Error fetching song from API for group=%s, song=%s: %s", group, song, err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		s.log.Errorf("Fetch song from api error, status: %d", res.StatusCode)
		return nil, fmt.Errorf("fetch song from api error, status: %d", res.StatusCode)
	}

	var songDetail dto.SongDetail

	err = json.NewDecoder(res.Body).Decode(&songDetail)
	if err != nil {
		s.log.Errorf("Error decoding response from API for group=%s, song=%s: %s", group, song, err)
		return nil, err
	}

	s.log.Infof("Successfully fetched song details from API: %+v", songDetail)
	return &songDetail, nil
}
