package service

import (
	"github.com/rekaime/r-mio/internal/utils/r-context"
	"github.com/rekaime/r-mio/api/repository"
)

type MusicService interface {
	GetMusicList() (*[]string, error)
	GetMusicById(id string) (*repository.Music, error)
}

type musicService struct {
	musicRepository repository.MusicRepository
}

func (service *musicService) GetMusicList() (*[]string, error) {
	ctx, cancel := rcontext.CreateTimeoutContext()
	defer cancel()
	return service.musicRepository.GetIdList(ctx)
}

func (service *musicService) GetMusicById(id string) (*repository.Music, error) {
	ctx, cancel := rcontext.CreateTimeoutContext()
	defer cancel()
	return service.musicRepository.FindById(ctx, id)
}

func NewMusicService(musicRepository repository.MusicRepository) MusicService {
	return &musicService{musicRepository}
}