package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/api/service"
)

type MusicController interface {
	GetMusicList(ctx *gin.Context)
	GetMusicById(ctx *gin.Context)
	HandleDownloadMusic(ctx *gin.Context)
}

type musicController struct {
	musicService  service.MusicService
	configService service.ConfigService
}

func (controller *musicController) GetMusicList(ctx *gin.Context) {
	musicList, err := controller.musicService.GetMusicList()
	if err != nil {
		InternalError(ctx)
		return
	}
	Success(ctx, musicList)
}

func (controller *musicController) GetMusicById(ctx *gin.Context) {
	id := ctx.Param("id")
	music, err := controller.musicService.GetMusicById(id)
	if err != nil {
		InternalError(ctx)
		return
	}
	Success(ctx, music)
}

func (controller *musicController) HandleDownloadMusic(ctx *gin.Context) {
	err := controller.musicService.HandleDownloadMusic()
	if err != nil {
		InternalError(ctx)
		return
	}
	Success(ctx, nil)
}

type Params struct {
	MusicService  service.MusicService
	ConfigService service.ConfigService
}

func NewMusicController(params Params) MusicController {
	return &musicController{
		musicService:  params.MusicService,
		configService: params.ConfigService,
	}
}
