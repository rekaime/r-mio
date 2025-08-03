package route

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/application"
	"github.com/rekaime/r-mio/mongo"
)

func NewMusicRoute(group *gin.RouterGroup, env *application.Env, db mongo.Database, timeout time.Duration) {
	// musicRepo := NewMusicRoute(db)
}