package application

import (
	"context"
	"fmt"
	"log"

	"github.com/rekaime/r-mio/internal/utils/r-context"
	"github.com/rekaime/r-mio/mongo"
)

func NewMongo(env *Env) (mongo.Client, AppCancelFunc) {
	ctx, cancel := rcontext.CreateTimeoutContext()
	defer cancel()

	var mongoURI string
	if env.DbUser == "" || env.DbPswd == "" {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%d", env.DbUser, env.DbPswd, env.DbHost, env.DbPort)
	} else {
		mongoURI = fmt.Sprintf("mongodb://%s:%d", env.DbHost, env.DbPort)
	}

	client, err := mongo.NewClient(mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client, func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Connection to MongoDB closed...")
	}
}
