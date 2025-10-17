package route

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/application"
	"github.com/rekaime/r-mio/mongo"
	"github.com/rekaime/r-mio/api/controller"
	"github.com/rekaime/r-mio/api/service"
	"github.com/rekaime/r-mio/api/repository"
)

func NewMusicRoute(group *gin.RouterGroup, env *application.Env, db mongo.Database, timeout time.Duration) {
	repo := repository.NewMusicRepository(db)
	service := service.NewMusicService(repo)
	controller := controller.NewMusicController(service)
	group.GET("/music-list", controller.GetMusicList)
	group.GET("/music/:id", controller.GetMusicById)
}