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
	dbClient, dbClientCancel := NewMongoClient(app.Env)
	app.DbClient = dbClient

	appCancelQueue = append(appCancelQueue, dbClientCancel)

	return app
}

func EndOfAppRunning() {
	for _, cancel := range appCancelQueue {
		cancel()
	}
}
