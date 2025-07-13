package controller

import (
	"github.com/rekaime/r-mio/api/service"
)

type MusicController interface { 
}

type musicController struct {
	musicService service.MusicService
}

func NewMusicController(musicService service.MusicService) MusicController {
	return musicController{musicService}
}