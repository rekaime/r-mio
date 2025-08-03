package application

import (
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

func EndOfAppRunning() {
	for _, cancel := range appCancelQueue {
		cancel()
	}
}
