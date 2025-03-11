package service

import (
	"context"
	"errors"
	"github.com/aaanger/music-library/internal/dto"
	"github.com/aaanger/music-library/internal/model"
	"github.com/aaanger/music-library/internal/repository"
	"github.com/sirupsen/logrus"
)

type IMusicService interface {
	AddSong(ctx context.Context, req *dto.AddSongReq) (*model.Song, error)
	GetSongsList(ctx context.Context, req *dto.GetSongsListReq) ([]*model.Song, error)
	GetSongLyrics(ctx context.Context, songID, limit, offset int) ([]*model.Verse, error)
	UpdateSong(ctx context.Context, songID int, req *dto.UpdateSongReq) error
	DeleteSong(ctx context.Context, songID int) error
}

type MusicService struct {
	repo   repository.IMusicRepository
	apiURL string
	log    *logrus.Logger
}

func NewMusicService(repo repository.IMusicRepository, apiURL string, log *logrus.Logger) *MusicService {
	return &MusicService{
		repo:   repo,
		apiURL: apiURL,
		log:    log,
	}
}

func (s *MusicService) AddSong(ctx context.Context, req *dto.AddSongReq) (*model.Song, error) {
	s.log.Infof("AddSong service: adding song - %s group - %s", req.Song, req.Group)

	songDetails, err := s.fetchSongFromAPI(req.Group, req.Song)
	if err != nil {
		return nil, err
	}

	s.log.Debugf("AddSong service: fetched from API - %+v", songDetails)

	song := model.Song{
		Song:        req.Song,
		Group:       req.Group,
		ReleaseDate: songDetails.ReleaseDate,
		Text:        songDetails.Text,
		Link:        songDetails.Link,
	}

	savedSong, err := s.repo.AddSong(ctx, &song)
	if err != nil {
		return nil, err
	}
	s.log.Infof("AddSong service: song successfully saved with ID - %d", savedSong.ID)

	err = s.repo.AddLyrics(ctx, savedSong.ID, savedSong.Text)
	if err != nil {
		return nil, err
	}

	return savedSong, nil
}

func (s *MusicService) GetSongsList(ctx context.Context, req *dto.GetSongsListReq) ([]*model.Song, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	s.log.Debugf("GetSongsList service: filters - %+v", req)

	songs, err := s.repo.GetSongsList(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(songs) == 0 {
		return nil, errors.New("couldn't find songs")
	}

	return songs, nil
}

func (s *MusicService) GetSongLyrics(ctx context.Context, songID, limit, offset int) ([]*model.Verse, error) {
	s.log.Debugf("GetSongLyrics service: songID=%d, limit=%d, offset=%d", songID, limit, offset)

	verses, err := s.repo.GetSongLyrics(ctx, songID, limit, offset)
	if err != nil {
		return nil, err
	}

	if verses == nil {
		return nil, errors.New("can't find lyrics for this song")
	}

	return verses, nil
}

func (s *MusicService) UpdateSong(ctx context.Context, songID int, req *dto.UpdateSongReq) error {
	s.log.Debugf("UpdateSong service: updating song with id %d with data - %+v", songID, req)
	return s.repo.UpdateSong(ctx, songID, req)
}

func (s *MusicService) DeleteSong(ctx context.Context, songID int) error {
	s.log.Infof("DeleteSong service: deleting song ID=%d", songID)
	return s.repo.DeleteSong(ctx, songID)
}
