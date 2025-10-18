package repository

import (
	"context"

	"github.com/rekaime/r-mio/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	MusicCollection = "musics"
)

var (
	MusicFileDir    string
	MusicDownloadFileDir string
)

type MusicItem struct {
	Size        int64    `bson:"size"`
	Path        string   `bson:"path"`
	Title       string   `bson:"title"`
	FileType    string   `bson:"file_type"`
	Artist      []string `bson:"artist"`
	Album       string   `bson:"album"`
	Composer    []string `bson:"composer"`
	AlbumArtist []string `bson:"album_artist"`
}

type MusicStatus struct {
	IsDisabled bool   `bson:"is_disabled"`
	BindID     string `bson:"bind_id"`
}

type Music struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	Item   MusicItem          `bson:"item"`
	Status MusicStatus        `bson:"status"`
}

type MusicRepository interface {
	FindById(ctx context.Context, id string) (*Music, error)
	FindByName(ctx context.Context, name string) (*Music, error)
	FindManyByName(ctx context.Context, name string) ([]*Music, error)
	GetIdList(ctx context.Context) ([]string, error)
	InsertOne(ctx context.Context, music *Music) error
}

type musicRepository struct {
	collection mongo.Collection
}

func (repo *musicRepository) FindById(ctx context.Context, id string) (*Music, error) {
	var music Music
	filter := bson.M{"_id": id}
	err := repo.collection.FindOne(ctx, filter).Decode(&music)
	if err != nil {
		return nil, err
	}
	return &music, nil
}

func (repo *musicRepository) FindByName(ctx context.Context, name string) (*Music, error) {
	var music Music
	filter := bson.M{"item.title": name}
	err := repo.collection.FindOne(ctx, filter).Decode(&music)
	if err != nil {
		return nil, err
	}
	return &music, nil
}

func (repo *musicRepository) FindManyByName(ctx context.Context, name string) ([]*Music, error) {
	musicList := []*Music{}
	filter := bson.M{"item.title": name}
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.Decode(&musicList)
	if err != nil {
		return nil, err
	}
	return musicList, nil
}

func (repo *musicRepository) GetIdList(ctx context.Context) ([]string, error) {
	idList := []string{}
	filter := bson.M{}
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var music Music
		err := cursor.Decode(&music)
		if err != nil {
			return nil, err
		}
		idList = append(idList, music.Id.Hex())
	}
	return idList, nil
}

func (repo *musicRepository) InsertOne(ctx context.Context, music *Music) error {
	_, err := repo.collection.InsertOne(ctx, music)
	if err != nil {
		return err
	}

	return nil
}

func NewMusicRepository(database mongo.Database) MusicRepository {
	return &musicRepository{collection: database.Collection(MusicCollection)}
}
