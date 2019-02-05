package testtask

import (
	"context"
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
	return s.pg.PutContext(ctx, itemId, locationIds)
}

func (s *CachedStore) GetContext(ctx context.Context, itemId store.ItemId) ([]store.Location, error) {
	locations, err := s.cache.Get(itemId)
	if locations != nil && err == nil {
		return locations, nil
	}

	locations, err = s.pg.GetContext(ctx, itemId)
	if err != nil {
		return nil, err
	}

	err = s.cache.Put(itemId, locations)
	if err != nil {
	}
	return locations, nil
}
