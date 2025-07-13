package service

import (
	"github.com/rekaime/r-mio/api/repository"
)

type MusicService interface {}

type musicService struct {
	musicRepository repository.MusicRepository
}

func NewMusicService(musicRepository repository.MusicRepository) MusicService {
	return &musicService{musicRepository}
}