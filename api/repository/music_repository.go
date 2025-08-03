package repository

import (
	"github.com/rekaime/r-mio/mongo"
)

const (
	MusicCollection = "musics"
)

type Music struct {
	Id     string `bson:"_id"`
	Path   string `bson:"path"`
	Name   string `bson:"name"`
	Author string `bson:"author"`
	Album  string `bson:"album"`
}

type MusicRepository interface {
}

type musicRepository struct {
	database   mongo.Database
	collection string
}

func NewMusicReporitory(database mongo.Database) MusicRepository {
	return &musicRepository{database: database, collection: MusicCollection}
}
