package repository

import (
	"context"
	"sync"
	"github.com/rekaime/r-mio/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	ConfigCollection = "config"
)

var (
	buffer *Config = nil
	mutex  sync.RWMutex
)

type Config struct {
	MusicDir         string `bson:"music_dir"`
	MusicDownloadDir string `bson:"music_download_dir"`
}

type ConfigRepository interface {
	Get(context.Context) (*Config, error)
}

type configRepository struct {
	collection mongo.Collection
}

func (repo *configRepository) Get(ctx context.Context) (*Config, error) {
	mutex.RLock()
	if buffer != nil {
		copy := *buffer
		mutex.RUnlock()
		return &copy, nil
	}
	mutex.RUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	var config Config
	filter := bson.M{}
	err := repo.collection.FindOne(ctx, filter).Decode(&config)
	if err != nil {
		return nil, err
	}
	buffer = &config
	copy := *buffer
	return &copy, nil
}

func NewConfigRepository(database mongo.Database) ConfigRepository {
	return &configRepository{
		collection: database.Collection(ConfigCollection),
	}
}
