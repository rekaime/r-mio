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

type MusicItem struct {
	Size        int64    `bson:"size" json:"size"`
	Path        string   `bson:"path" json:"path"`
	Title       string   `bson:"title" json:"title"`
	FileType    string   `bson:"file_type" json:"file_type"`
	Artist      []string `bson:"artist" json:"artist"`
	Album       string   `bson:"album" json:"album"`
	Composer    []string `bson:"composer" json:"composer"`
	AlbumArtist []string `bson:"album_artist" json:"album_artist"`
}

type MusicStatus struct {
	IsDisabled bool   `bson:"is_disabled" json:"is_disabled"`
	BindID     string `bson:"bind_id" json:"bind_id"`
}

type Music struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Item   MusicItem          `bson:"item" json:"item"`
	Status MusicStatus        `bson:"status" json:"status"`
}

type MusicRepository interface {
	FindById(ctx context.Context, id string) (*Music, error)
	FindByName(ctx context.Context, name string) (*Music, error)
	FindManyByName(ctx context.Context, name string) ([]*Music, error)
	GetIdList(ctx context.Context) ([]string, error)
	InsertOne(ctx context.Context, music *Music) (string, error)
	DeleteOne(ctx context.Context, id string) error
}

type musicRepository struct {
	collection mongo.Collection
}

func (repo *musicRepository) FindById(ctx context.Context, id string) (*Music, error) {
	var music Music
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	err = repo.collection.FindOne(ctx, filter).Decode(&music)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
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

func (repo *musicRepository) InsertOne(ctx context.Context, music *Music) (string, error) {
	result, err := repo.collection.InsertOne(ctx, music)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (repo *musicRepository) DeleteOne(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	_, err = repo.collection.DeleteOne(ctx, filter)
	return err
}

func NewMusicRepository(database mongo.Database) MusicRepository {
	return &musicRepository{collection: database.Collection(MusicCollection)}
}
