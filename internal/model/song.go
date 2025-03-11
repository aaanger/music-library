package model

type Song struct {
	ID          int
	Song        string `json:"song" binding:"required"`
	Group       string `json:"group" binding:"required"`
	ReleaseDate string `json:"release_date" time_format:"2006-01-02"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
