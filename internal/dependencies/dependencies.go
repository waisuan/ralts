package dependencies

import (
	"ralts/internal/config"
)

type Dependencies struct {
	Cfg     *config.Config
	Storage CoreStorageInterface
	Cache   CoreCacheInterface
}

func NewDependencies(cfg *config.Config) *Dependencies {
	storage := NewDB(cfg)

	cache := NewCache(cfg)

	return &Dependencies{
		Cfg:     cfg,
		Storage: storage,
		Cache:   cache,
	}
}

func (deps *Dependencies) Disconnect() {
	deps.Storage.Close()
}
