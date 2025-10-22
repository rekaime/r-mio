package service

import (
	
	"github.com/rekaime/r-mio/internal/utils/r-context"
	"github.com/rekaime/r-mio/api/repository"
)

type ConfigService interface {
	Get() (*repository.Config, error)
}

type configService struct {
	configRepository repository.ConfigRepository
}

func (s *configService) Get() (*repository.Config, error) {
	ctx, cancel := rcontext.CreateTimeoutContext()
	defer cancel()
	
	return s.configRepository.Get(ctx)
}

func NewConfigService(configRepository repository.ConfigRepository) ConfigService {
	return &configService{
		configRepository: configRepository,
	}
}
