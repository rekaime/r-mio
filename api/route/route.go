package route

import (
	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/application"
	"github.com/rekaime/r-mio/mongo"
	"time"
)

func InitRouter(gin *gin.Engine, env *application.Env, db mongo.Database, timeout time.Duration) {
	publicRoute := gin.Group("/api")
	NewMusicRoute(publicRoute, env, db, timeout)
}
