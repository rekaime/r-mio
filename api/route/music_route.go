package route

import (
	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/api/controller"
	"github.com/rekaime/r-mio/api/repository"
	"github.com/rekaime/r-mio/api/service"
	"github.com/rekaime/r-mio/application"
	"github.com/rekaime/r-mio/mongo"
	"time"
)

func NewMusicRoute(group *gin.RouterGroup, env *application.Env, db mongo.Database, timeout time.Duration) {
	musicRepo := repository.NewMusicRepository(db)
	configRepo := repository.NewConfigRepository(db)
	musicService := service.NewMusicService(musicRepo)
	configService := service.NewConfigService(configRepo)
	controller := controller.NewMusicController(controller.Params{
		MusicService:  musicService,
		ConfigService: configService,
	})

	group.GET("/music-list", controller.GetMusicList)
	group.GET("/music/:id", controller.GetMusicById)
	group.POST("/music/handle-download", controller.HandleDownloadMusic)
}
