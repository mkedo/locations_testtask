package testtask

import "testtask/store"

type CachedStore struct {
	pg    *PgStore
	cache ItemLocationCache
}

type CachedStoreOptions struct {
	PgStore           *PgStore
	ItemLocationCache ItemLocationCache
}

func NewCachedStore(config *CachedStoreOptions) *CachedStore {
	return &CachedStore{pg: config.PgStore, cache: config.ItemLocationCache}
}

func (s *CachedStore) Put(itemId store.ItemId, locationIds []store.LocationId) error {
	return s.pg.Put(itemId, locationIds)
}

func (s *CachedStore) Get(itemId store.ItemId) ([]store.Location, error) {
	locations, err := s.cache.Get(itemId)
	if locations != nil && err == nil {
		return locations, nil
	}

	locations, err = s.pg.Get(itemId)
	if err != nil {
		return nil, err
	}

	err = s.cache.Put(itemId, locations)
	if err != nil {
	}
	return locations, nil
}
