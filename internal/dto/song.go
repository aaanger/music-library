package dto

type AddSongReq struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

type GetSongsListReq struct {
	Song        *string `json:"song,omitempty"`
	Group       *string `json:"group,omitempty"`
	ReleaseDate *string `json:"releaseDate,omitempty"`
	Limit       int
	Offset      int
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type UpdateSongReq struct {
	Song        *string `json:"song,omitempty"`
	Group       *string `json:"group,omitempty"`
	ReleaseDate *string `json:"releaseDate,omitempty"`
	Text        *string `json:"text,omitempty"`
	Link        *string `json:"link,omitempty"`
}
