package route

import (
	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/application"
	"github.com/rekaime/r-mio/mongo"
	"time"
)

func InitRouter(gin *gin.Engine, app *application.Application, db mongo.Database, timeout time.Duration) {
	publicRoute := gin.Group("/api")
	NewMusicRoute(publicRoute, app.Env, db, timeout)
}
