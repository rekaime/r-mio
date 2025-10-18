package service

import (
	"path/filepath"
	"strings"
	"os"
	"log"

	"github.com/dhowden/tag"
	"github.com/rekaime/r-mio/internal/utils/r-context"
	"github.com/rekaime/r-mio/api/repository"
)

type MusicService interface {
	GetMusicList() ([]string, error)
	GetMusicById(id string) (*repository.Music, error)
	HasMusic(name string) bool
	HandleDownloadMusic() error
	getDownloadMusicFilepath() ([]string, error)
	getMusicMetadata(path string) (*repository.Music, error)
	moveMusicFile(path string) error
}

type musicService struct {
	musicRepository repository.MusicRepository
}

func (service *musicService) GetMusicList() ([]string, error) {
	ctx, cancel := rcontext.CreateTimeoutContext()
	defer cancel()
	return service.musicRepository.GetIdList(ctx)
}

func (service *musicService) GetMusicById(id string) (*repository.Music, error) {
	ctx, cancel := rcontext.CreateTimeoutContext()
	defer cancel()
	return service.musicRepository.FindById(ctx, id)
}

func (service *musicService) HasMusic(name string) bool {
	ctx, cancel := rcontext.CreateTimeoutContext()
	defer cancel()

	music, err := service.musicRepository.FindByName(ctx, name)
	if err != nil {
		return false
	}
	return music != nil
}

// 尝试将下载目录下的音频文件移动到音乐目录 同时入库
func (service *musicService) HandleDownloadMusic() error {
	ctx, cancel := rcontext.CreateTimeoutContext()
	defer cancel()

	audioFiles, err := service.getDownloadMusicFilepath()
	if err != nil {
		return err
	}

	for i, filePath := range audioFiles {
		log.Printf("处理文件 %d/%d: %s\n", i+1, len(audioFiles), filePath)

		info, err := service.getMusicMetadata(filePath)
		if err != nil {
			log.Printf("HandleDownloadMusic(1): %v\n", err)
			continue
		}

		has := service.HasMusic(info.Item.Title)
		if has {
			log.Printf("文件 \"%s\" 已存在\n", info.Item.Title)
			continue
		}

		err = service.musicRepository.InsertOne(ctx, info)
		if err != nil {
			log.Printf("HandleDownloadMusic(2): %v\n", err)
			continue
		}

		err = service.moveMusicFile(filePath)
		if err != nil {
			log.Printf("HandleDownloadMusic(3): %s\n=> err: %v\n", filePath, err)
			continue
		}
	}
	return nil
}

// 获取下载目录下的所有音频
func (service *musicService) getDownloadMusicFilepath() ([]string, error) {
	var audioFiles []string
	audioExtensions := map[string]bool{
		".mp3":  true,
		".flac": true,
		".wav":  true,
		".ogg":  true,
		".m4a":  true,
		".aac":  true,
	}

	entries, err := os.ReadDir(repository.MusicDownloadFileDir)
	if err != nil {
		return audioFiles, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(repository.MusicDownloadFileDir, entry.Name())

		ext := strings.ToLower(filepath.Ext(path))
		if audioExtensions[ext] {
			audioFiles = append(audioFiles, path)
		}
	}

	return audioFiles, err
}

// 音频文件 path -> repository.Music
func (service *musicService) getMusicMetadata(path string) (*repository.Music, error) {
	var info repository.Music
	info.Item.Path = filepath.Base(path)
	info.Item.FileType = strings.ToLower(filepath.Ext(path))
	info.Status.IsDisabled = false

	f, err := os.Open(path)
	if err != nil {
		return &info, err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return &info, err
	}
	info.Item.Size = stat.Size()

	// 解析元数据
	tg, err := tag.ReadFrom(f)
	if err != nil {
		info.Item.Title = strings.TrimSuffix(filepath.Base(path), info.Item.FileType)
		return &info, err
	}

	info.Item.Title = tg.Title()
	if info.Item.Title == "" {
		info.Item.Title = strings.TrimSuffix(filepath.Base(path), info.Item.FileType)
	}
	
	split := "/"
	info.Item.Artist = strings.Split(tg.Artist(), split)
	info.Item.Composer = strings.Split(tg.Composer(), split)
	info.Item.AlbumArtist = strings.Split(tg.AlbumArtist(), split)

	info.Item.Album = tg.Album()

	return &info, nil
}

func (service *musicService) moveMusicFile(path string) error {
	basename := filepath.Base(path)
	newPath := filepath.Join(repository.MusicFileDir, basename)
	return os.Rename(path, newPath)
}

func NewMusicService(musicRepository repository.MusicRepository) MusicService {
	return &musicService{musicRepository}
}