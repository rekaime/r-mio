package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/api/route"
	"github.com/rekaime/r-mio/application"
	"time"
)

func main() {
	app := application.App()
	env := app.Env

	defer application.EndOfAppRunning()

	db := app.DbClient.Database(env.DbName)

	ginEngine := gin.Default()
	timeout := time.Duration(3) * time.Second

	route.Run(ginEngine, env, db, timeout)
	ginEngine.Run(fmt.Sprintf("%s:%d", env.Ip, env.Port))
}
