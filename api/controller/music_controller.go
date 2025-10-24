package controller

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/api/repository"
	"github.com/rekaime/r-mio/api/service"
)

type MusicResponse struct {
	Id    string           `json:"id"`
	Music []byte           `json:"music"`
	Cover []byte           `json:"cover"`
	Info  repository.Music `json:"info"`
}

type MusicController interface {
	GetMusicList(ctx *gin.Context)
	GetMusicInfoById(ctx *gin.Context)
	GetMusicFileById(ctx *gin.Context)
	GetMusicCoverById(ctx *gin.Context)
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

func (controller *musicController) GetMusicInfoById(ctx *gin.Context) {
	id := ctx.Param("id")
	music, err := controller.musicService.GetMusicById(id)
	if err != nil {
		InternalError(ctx)
		return
	}

	Success(ctx, *music)
}

func (controller *musicController) GetMusicFileById(ctx *gin.Context) {
	id := ctx.Param("id")
	music, err := controller.musicService.GetMusicById(id)
	if err != nil {
		InternalError(ctx)
		return
	}

	config, err := controller.configService.Get()
	if err != nil {
		InternalError(ctx)
		return
	}
	path := filepath.Join(config.MusicDir, music.Item.Path)
	audio, err := controller.musicService.ReadLocalMusic(path)
	if err != nil {
		InternalError(ctx)
		return
	}

	Stream(ctx, audio)
}

func (controller *musicController) GetMusicCoverById(ctx *gin.Context) {
	id := ctx.Param("id")
	music, err := controller.musicService.GetMusicById(id)
	if err != nil {
		InternalError(ctx)
		return
	}

	config, err := controller.configService.Get()
	if err != nil {
		InternalError(ctx)
		return
	}
	path := filepath.Join(config.MusicDir, music.Item.Path)
	cover, err := controller.musicService.ReadLocalMusicCover(path)
	if err != nil {
		InternalError(ctx)
		return
	}

	Data(ctx, cover)
}

func (controller *musicController) HandleDownloadMusic(ctx *gin.Context) {
	config, err := controller.configService.Get()
	if err != nil {
		InternalError(ctx)
		return
	}
	err = controller.musicService.HandleDownloadMusic(config.MusicDir, config.MusicDownloadDir)
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
