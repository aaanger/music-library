package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aaanger/music-library/internal/dto"
	"github.com/aaanger/music-library/internal/model"
	"github.com/sirupsen/logrus"
	"strings"
)

type IMusicRepository interface {
	AddSong(ctx context.Context, req *model.Song) (*model.Song, error)
	AddLyrics(ctx context.Context, songID int, lyrics string) error
	GetSongsList(ctx context.Context, req *dto.GetSongsListReq) ([]*model.Song, error)
	GetSongLyrics(ctx context.Context, songID, limit, offset int) ([]*model.Verse, error)
	UpdateSong(ctx context.Context, songID int, req *dto.UpdateSongReq) error
	DeleteSong(ctx context.Context, songID int) error
}

type MusicRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewMusicRepository(db *sql.DB, log *logrus.Logger) *MusicRepository {
	return &MusicRepository{
		db:  db,
		log: log,
	}
}

func (r *MusicRepository) AddSong(ctx context.Context, song *model.Song) (*model.Song, error) {
	row := r.db.QueryRowContext(ctx, `INSERT INTO songs (song, artist, release_date, link) VALUES($1, $2, $3, $4) RETURNING id;`,
		song.Song, song.Group, song.ReleaseDate, song.Link)

	err := row.Scan(&song.ID)
	if err != nil {
		r.log.Errorf("AddSong repository error: %s", err)
		return nil, err
	}

	r.log.Infof("Successfully added song to DB: %+v", song)
	return song, nil
}

func (r *MusicRepository) AddLyrics(ctx context.Context, songID int, lyrics string) error {
	verses := strings.Split(lyrics, "\n\n")
	query := `INSERT INTO verses (song_id, verse_number, verse_lyrics) VALUES($1, $2, $3)`

	for i, verse := range verses {
		_, err := r.db.ExecContext(ctx, query, songID, i+1, verse)
		if err != nil {
			r.log.Errorf("AddLyrics repository error: %s", err)
			return err
		}
	}

	r.log.Infof("Successfully added %d verses for song id %d", len(verses), songID)
	return nil
}

func (r *MusicRepository) GetSongsList(ctx context.Context, req *dto.GetSongsListReq) ([]*model.Song, error) {
	query := `SELECT id, song, artist, release_date, link FROM songs`

	keys := make([]string, 0)
	values := make([]interface{}, 0)
	arg := 1

	if req.Song != nil {
		keys = append(keys, fmt.Sprintf("song=$%d", arg))
		values = append(values, *req.Song)
		arg++
	}
	if req.Group != nil {
		keys = append(keys, fmt.Sprintf("artist=$%d", arg))
		values = append(values, *req.Group)
		arg++
	}
	if req.ReleaseDate != nil {
		keys = append(keys, fmt.Sprintf("release_date=$%d", arg))
		values = append(values, *req.ReleaseDate)
		arg++
	}

	if len(keys) > 0 {
		query += " WHERE " + strings.Join(keys, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", arg, arg+1)

	values = append(values, req.Limit, req.Offset)

	r.log.Debugf("GetSongsList repository filters: song - %v group - %v releaseDate - %v limit - %v offset - %v", req.Song, req.Group, req.ReleaseDate, req.Limit, req.Offset)
	r.log.Debugf("GetSongsList repository: executing sql query: %s", query)
	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		r.log.Errorf("GetSongsList repository error: %s", err)
		return nil, err
	}
	defer rows.Close()

	var songs []*model.Song

	for rows.Next() {
		var song model.Song
		err = rows.Scan(&song.ID, &song.Song, &song.Group, &song.ReleaseDate, &song.Link)
		if err != nil {
			r.log.Errorf("GetSongsList repository error: %s", err)
			return nil, err
		}

		songs = append(songs, &song)
	}

	r.log.Infof("Successfully got songs list: %+v", songs)
	return songs, nil
}

func (r *MusicRepository) GetSongLyrics(ctx context.Context, songID, limit, offset int) ([]*model.Verse, error) {
	r.log.Debugf("GetSongLyrics repository: songID - %d limit - %d offset - %d", songID, limit, offset)

	rows, err := r.db.QueryContext(ctx, `SELECT verse_number, verse_lyrics FROM verses WHERE song_id=$1 ORDER BY verse_number LIMIT $2 OFFSET $3`,
		songID, limit, offset)
	if err != nil {
		r.log.Errorf("GetSongLyrics repository error: %s", err)
		return nil, err
	}
	defer rows.Close()

	var verses []*model.Verse

	for rows.Next() {
		var verse model.Verse

		err := rows.Scan(&verse.Number, &verse.Lyrics)
		if err != nil {
			r.log.Errorf("GetSongLyrics repository error: %s", err)
			return nil, err
		}

		verses = append(verses, &verse)
	}

	r.log.Infof("Successfully got lyrics for song id %d: %+v", songID, verses)
	return verses, nil
}

func (r *MusicRepository) UpdateSong(ctx context.Context, songID int, req *dto.UpdateSongReq) error {
	keys := make([]string, 0)
	values := make([]interface{}, 0)
	arg := 1

	if req.Song != nil {
		keys = append(keys, fmt.Sprintf("song=$%d", arg))
		values = append(values, *req.Song)
		arg++
	}
	if req.Group != nil {
		keys = append(keys, fmt.Sprintf("artist=$%d", arg))
		values = append(values, *req.Group)
		arg++
	}
	if req.ReleaseDate != nil {
		keys = append(keys, fmt.Sprintf("release_date=$%d", arg))
		values = append(values, *req.ReleaseDate)
		arg++
	}

	joinKeys := strings.Join(keys, ", ")

	r.log.Debugf("UpdateSong repository input parameters: song - %v group - %v releaseDate - %v", req.Song, req.Group, req.ReleaseDate)

	query := fmt.Sprintf("UPDATE songs SET %s WHERE id=$%d", joinKeys, arg)

	values = append(values, songID)

	_, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		r.log.Errorf("UpdateSong repository error: %s", err)
		return err
	}

	r.log.Infof("Successfully updated song with id %d", songID)
	return nil
}

func (r *MusicRepository) DeleteSong(ctx context.Context, songID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM songs WHERE id = $1;`, songID)
	if err != nil {
		r.log.Errorf("DeleteSong repository error: %s", err)
		return err
	}

	r.log.Infof("Successfully deleted song with id %d", songID)
	return nil
}
