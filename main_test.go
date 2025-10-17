package main_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/rekaime/r-mio/internal/utils/r-context"
	"github.com/rekaime/r-mio/mongo"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

type Env struct {
	Ip       string `mapstructure:"R_IP"`
	Port     int    `mapstructure:"R_PORT"`
	DbHost   string `mapstructure:"R_DB_HOST"`
	DbPort   int    `mapstructure:"R_DB_PORT"`
	DbUser   string `mapstructure:"R_DB_USER"`
	DbPswd   string `mapstructure:"R_DB_PSWD"`
	DbName   string `mapstructure:"R_DB_NAME"`
	MusicDir string `mapstructure:"R_MUSIC_DIR"`
}

var env *Env

func NewEnv() *Env {
	if env != nil {
		return env
	}

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return env
}

func NewMongo(env *Env) mongo.Client {
	ctx, cancel := rcontext.CreateTimeoutContext()
	defer cancel()

	var mongoURI string
	if env.DbUser != "" && env.DbPswd != "" {
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

	return client
}

func TestRun(t *testing.T) {
	ctx, cancel := rcontext.CreateTimeoutContext()
	defer cancel()

	env := NewEnv()
	client := NewMongo(env)
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database(env.DbName)
	collection := db.Collection("music_info")

	// doc := bson.D{{"name", "测试文档"}, {"value", 123}}
	doc := bson.D{}
	cursor, err := collection.Find(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}
	// t.Logf("Inserted ID: %v", cursor.Next(ctx))
	for cursor.Next(ctx) {
		var user bson.D
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		t.Logf("%v", user)
	}
}
