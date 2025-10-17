package main

import (
	"fmt"
	"log"
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

	route.InitRouter(ginEngine, env, db, timeout)
	err := ginEngine.Run(fmt.Sprintf("%s:%d", env.Ip, env.Port))
	if err != nil {
		log.Fatal(err)
	}
}
