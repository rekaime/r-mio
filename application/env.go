package application

import (
	"log"
	"github.com/spf13/viper"
)

type Env struct {
	R_IP 	string  `mapstructure:"R_IP"`
	R_PORT	int 	`mapstructure:"R_PORT"`
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
	log.Println(env)

	return env
}