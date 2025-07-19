package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/cmd"
)

type Application struct {
	Env *Env
	Cmd *cmd.Cmd
}

func App() *Application {
	app := &Application{}

	app.Env = NewEnv()
	app.Cmd = cmd.NewCmd()

	return app
}

func (app *Application) Run() { 
	router := gin.Default()

	env := app.Env
	router.Run(fmt.Sprintf("%s:%d", env.R_IP, env.R_PORT))
}