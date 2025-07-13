package repository

type MusicRepository interface {}

type musicRepository struct {

}

func NewMusicReporitory() MusicRepository {
	return &musicRepository{}
}