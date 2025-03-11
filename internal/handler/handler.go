package handler

import (
	_ "github.com/aaanger/music-library/docs"
	"github.com/aaanger/music-library/internal/dto"
	"github.com/aaanger/music-library/internal/service"
	"github.com/aaanger/music-library/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type MusicHandler struct {
	service service.IMusicService
	log     *logrus.Logger
}

func NewMusicHandler(service service.IMusicService, log *logrus.Logger) *MusicHandler {
	return &MusicHandler{
		service: service,
		log:     log,
	}
}

// AddSong godoc
// @Summary Добавление новой песни
// @Tags Songs
// @Produce json
// @Param song body dto.AddSongReq true "Данные для добавления песни"
// @Success 200 {object} model.Song "Данные песни"
// @Failure 400 {string} string "Неверное тело запроса"
// @Failure 500 {string} string "Ошибка добавления песни на сервере"
// @Router /api/v1/add [post]
func (h *MusicHandler) AddSong(c *gin.Context) {
	var req dto.AddSongReq

	err := c.BindJSON(&req)
	if err != nil {
		h.log.Debugf("AddSong handler: Invalid input parameters: %s", err)
		response.Error(c, http.StatusBadRequest, "invalid input parameters")
		return
	}

	h.log.Infof("AddSong handler request: song - %s, group - %s", req.Song, req.Group)

	song, err := h.service.AddSong(c, &req)
	if err != nil {
		h.log.Errorf("AddSong failure: %s", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.Infof("AddSong handler successful response: %+v", song)
	response.JSON(c, song)
}

// GetSongsList godoc
// @Summary Получение данных библиотеки с фильтрацией по всем полям и пагинацией
// @Tags Songs
// @Produce json
// @Param song query string false "Фильтр по названию песни"
// @Param group query string false "Фильтр по названию исполнителя"
// @Param release_date query int false "Фильтр по дате выпуска"
// @Param limit query int false "Количество песен на странице" default(10)
// @Param page query int false "Номер страницы" default(1)
// @Success 200 {array} model.Song
// @Failure 400 {string} string "Некорректный фильтр или параметры пагинации"
// @Failure 500 {string} string "Ошибка получения данных"
// @Router /api/v1/songs [get]
func (h *MusicHandler) GetSongsList(c *gin.Context) {
	var req dto.GetSongsListReq

	song := c.Query("song")
	req.Song = &song

	group := c.Query("group")
	req.Group = &group

	releaseDate := c.Query("release_date")
	req.ReleaseDate = &releaseDate

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		h.log.Debugf("GetSongsList handler: invalid limit query: %s", err)
		response.Error(c, http.StatusBadRequest, "invalid limit")
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		h.log.Debugf("GetSongsList handler: invalid limit query: %s", err)
		response.Error(c, http.StatusBadRequest, "invalid page")
		return
	}

	h.log.Debugf("GetSongsList handler request: song - %s, group - %s, releaseDate - %s, limit - %v, page - %v", song, group, releaseDate, limit, page)

	req.Limit = limit
	req.Offset = (page - 1) * limit

	songs, err := h.service.GetSongsList(c, &req)
	if err != nil {
		h.log.Errorf("GetSongsList failure: %s", err)
		response.Error(c, http.StatusInternalServerError, "failed to get songs")
		return
	}

	h.log.Infof("GetSongsList handler successful response: %+v", songs)
	response.JSON(c, songs)
}

// GetSongLyrics godoc
// @Summary Получение текста песни с пагинацией по куплетам
// @Tags Songs
// @Produce json
// @Param songID path int true "ID песни"
// @Param limit query int false "количество куплетов на странице" default(3)
// @Param page query int false "номер страницы" default(1)
// @Success 200 {array} model.Verse
// @Failure 400 {string} string "Неверный ID песни или некорректные параметры пагинации"
// @Failure 500 {string} string "Ошибка получения текста песни"
// @Router /api/v1/{songID}/lyrics [get]
func (h *MusicHandler) GetSongLyrics(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("songID"))
	if err != nil {
		h.log.Debugf("GetSongLyrics handler: invalid song id: %s", err)
		response.Error(c, http.StatusBadRequest, "invalid song id")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "3"))
	if err != nil || limit <= 0 {
		h.log.Debugf("GetSongLyrics handler: invalid limit: %s", err)
		response.Error(c, http.StatusBadRequest, "invalid limit")
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		h.log.Debugf("GetSongLyrics handler: invalid page: %s", err)
		response.Error(c, http.StatusBadRequest, "invalid page")
		return
	}

	h.log.Debugf("GetSongsLyrics handler request: songID - %v, limit - %v, page - %v", songID, limit, page)

	offset := (page - 1) * limit

	verses, err := h.service.GetSongLyrics(c, songID, limit, offset)
	if err != nil {
		h.log.Errorf("GetSongLyrics failure: %s", err)
		response.Error(c, http.StatusInternalServerError, "failed to get text")
		return
	}

	h.log.Infof("GetSongsLyrics handler successful response: %+v", verses)
	response.JSON(c, verses)
}

// UpdateSong godoc
// @Summary Изменение данных песни
// @Tags Songs
// @Produce json
// @Param songID path int true "ID песни"
// @Param song body dto.UpdateSongReq true "Данные для изменения"
// @Success 200 {string} string "Данные песни изменены"
// @Failure 400 {string} string "Неверное тело запроса или ID песни"
// @Failure 500 {string} string "Ошибка изменения песни на сервере"
// @Router /api/v1/{songID} [put]
func (h *MusicHandler) UpdateSong(c *gin.Context) {
	var req dto.UpdateSongReq

	err := c.BindJSON(&req)
	if err != nil {
		h.log.Debugf("UpdateSong handler: invalid input parameters: %s", err)
		response.Error(c, http.StatusBadRequest, "invalid input parameters")
		return
	}

	songID, err := strconv.Atoi(c.Param("songID"))
	if err != nil {
		h.log.Debugf("UpdateSong handler: invalid song id: %s", err)
		response.Error(c, http.StatusBadRequest, "invalid song id")
		return
	}

	h.log.Debugf("UpdateSong handler request: songID - %v", songID)

	err = h.service.UpdateSong(c, songID, &req)
	if err != nil {
		h.log.Errorf("UpdateSong failure: %s", err)
		response.Error(c, http.StatusInternalServerError, "failed to update song")
		return
	}

	h.log.Infof("UpdateSong handler successful response")
	response.JSON(c, "successfully updated song")
}

// DeleteSong godoc
// @Summary Удаление песни
// @Tags Songs
// @Produce json
// @Param songID path int true "ID песни"
// @Success 200 {string} string "Песня удалена"
// @Failure 400 {string} string "Неверный ID песни"
// @Failure 500 {string} string "Ошибка удаления песни"
// @Router /api/v1/{songID} [delete]
func (h *MusicHandler) DeleteSong(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("songID"))
	if err != nil {
		h.log.Debugf("DeleteSong handler: invalid song id: %s", err)
		response.Error(c, http.StatusBadRequest, "invalid song id")
		return
	}

	err = h.service.DeleteSong(c, songID)
	if err != nil {
		h.log.Errorf("DeleteSong failure: %s", err)
		response.Error(c, http.StatusInternalServerError, "failed to delete song")
		return
	}

	h.log.Infof("DeleteSong handler successful response")
	response.JSON(c, "successfully deleted song")
}
