package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/cmd"
	"github.com/rekaime/r-mio/mongo"
)

type AppCancelFunc func()

var appCancelQueue []AppCancelFunc

type Application struct {
	Env      *Env
	Cmd      *cmd.Cmd
	DbClient mongo.Client
}

func init() {
	appCancelQueue = make([]AppCancelFunc, 0)
}

func App() *Application {
	app := &Application{}

	app.Env = NewEnv()
	app.Cmd = cmd.NewCmd()
	db, dbCancel := NewMongo(app.Env)
	app.DbClient = db

	appCancelQueue = append(appCancelQueue, dbCancel)

	return app
}

func endOfAppRunning() {
	for _, cancel := range appCancelQueue {
		cancel()
	}
}

func (app *Application) Run() {
	defer endOfAppRunning()

	db := app.DbClient.Database(env.DbName)

	router := gin.Default()

	env := app.Env
	router.Run(fmt.Sprintf("%s:%d", env.Ip, env.Port))
}
