package repository

import (
	"context"
	
	"go.mongodb.org/mongo-driver/bson"
	"github.com/rekaime/r-mio/mongo"
)

const (
	MusicCollection = "musics"
)

const (
	FieldId = "_id"
	FieldPath = "path"
	FieldName = "name"
	FieldAuthor = "author"
	FieldAlbum = "album"
)

type Music struct {
	Id     string `bson:"_id"`
	Path   string `bson:"path"`
	Name   string `bson:"name"`
	Author string `bson:"author"`
	Album  string `bson:"album"`
}

type MusicRepository interface {
	FindById(ctx context.Context, id string) (*Music, error)
	GetIdList(ctx context.Context) (*[]string, error)
}

type musicRepository struct {
	collection mongo.Collection
}

func (repo *musicRepository) FindById(ctx context.Context, id string) (*Music, error) {
	var music Music
	filter := bson.M{FieldId: id}
	err := repo.collection.FindOne(ctx, filter).Decode(&music)
	if err != nil {
		return nil, err
	}
	return &music, nil
}

func (repo *musicRepository) GetIdList(ctx context.Context) (*[]string, error) {
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
		idList = append(idList, music.Id)
	}
	return &idList, nil
}

func NewMusicRepository(database mongo.Database) MusicRepository {
	return &musicRepository{collection: database.Collection(MusicCollection)}
}
