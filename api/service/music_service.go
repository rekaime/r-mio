package service

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	"github.com/rekaime/r-mio/api/repository"
	rcontext "github.com/rekaime/r-mio/internal/utils/r-context"
)

type MusicService interface {
	GetMusicList() ([]string, error)
	GetMusicById(string) (*repository.Music, error)
	HasMusic(string) bool
	HandleDownloadMusic(MusicDir string, MusicDownloadDir string) error
	getDownloadMusicFilepath(string) ([]string, error)
	getMusicMetadata(string) (*repository.Music, error)
	moveMusicFile(string, string) error
	ReadLocalMusicCover(string) ([]byte, error)
	ReadLocalMusic(string) (io.ReadSeekCloser, error)
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
func (service *musicService) HandleDownloadMusic(MusicDir string, MusicDownloadDir string) error {
	ctx, cancel := rcontext.CreateTimeoutContext()
	defer cancel()

	audioFiles, err := service.getDownloadMusicFilepath(MusicDownloadDir)
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

		id, err := service.musicRepository.InsertOne(ctx, info)
		if err != nil {
			log.Printf("HandleDownloadMusic(2): %v\n", err)
			continue
		}

		destPath := filepath.Join(MusicDir, id)
		err = service.moveMusicFile(filePath, destPath)
		if err != nil {
			log.Printf("HandleDownloadMusic(3): %s\n=> err: %v\n", filePath, err)
			err = service.musicRepository.DeleteOne(ctx, id)
			if err != nil {
				log.Printf("HandleDownloadMusic(4): %v\n", err)
			}
			continue
		}
	}
	return nil
}

// 获取下载目录下的所有音频
func (service *musicService) getDownloadMusicFilepath(musicDownloadFileDir string) ([]string, error) {
	var audioFiles []string
	audioExtensions := map[string]bool{
		".mp3":  true,
		".flac": true,
		".wav":  true,
		".ogg":  true,
		".m4a":  true,
		".aac":  true,
	}

	entries, err := os.ReadDir(musicDownloadFileDir)
	if err != nil {
		return audioFiles, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(musicDownloadFileDir, entry.Name())

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

func (service *musicService) moveMusicFile(originPath string, destPath string) error {
	return os.Rename(originPath, destPath)
}

func (service *musicService) ReadLocalMusicCover(path string) (cover []byte, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tg, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	cover = tg.Picture().Data
	
	return cover, nil
}

func (service *musicService) ReadLocalMusic(path string) (io.ReadSeekCloser, error) {
	return os.Open(path)
}

func NewMusicService(musicRepository repository.MusicRepository) MusicService {
	return &musicService{musicRepository}
}
