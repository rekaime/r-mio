package application

import (
	"github.com/spf13/viper"
	"log"
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
