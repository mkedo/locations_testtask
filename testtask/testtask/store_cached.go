package testtask

import (
	"context"
	"log"
	"testtask/store"
)

type CachedStore struct {
	pg    *PgStore
	cache ItemLocationCache
}

type CachedStoreOptions struct {
	PgStore           *PgStore
	ItemLocationCache ItemLocationCache
}

// Кеширующее хранилище.
// Если информация по адресам объявления есть в кеше, то она выдается из кеша.
// Если информации нет в кеше, то данные запрашиваются из хранилища на postgres.
func NewCachedStore(config *CachedStoreOptions) *CachedStore {
	return &CachedStore{pg: config.PgStore, cache: config.ItemLocationCache}
}

func (s *CachedStore) PutContext(ctx context.Context, itemId store.ItemId, locationIds []store.LocationId) error {
	err := s.pg.PutContext(ctx, itemId, locationIds)
	if err == nil {
		if err := s.cache.Invalidate(itemId); err != nil {
			log.Println(err)
		}
	}
	return err
}

func (s *CachedStore) GetContext(ctx context.Context, itemId store.ItemId) ([]store.Location, error) {
	locations, err := s.cache.Get(itemId)
	if err != nil {
		if err != ItemLocationCacheMiss {
			log.Println(err)
		}
	} else {
		return locations, nil
	}

	locations, err = s.pg.GetContext(ctx, itemId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = s.cache.Put(itemId, locations)
	if err != nil {
		log.Println(err)
	}
	return locations, nil
}
