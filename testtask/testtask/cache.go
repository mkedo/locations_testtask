package testtask

import "testtask/store"

const ItemLocationCacheMiss = CacheError("store: cache miss")

type CacheError string

func (e CacheError) Error() string { return string(e) }

// Кеш, хранящий результаты получения Location для объявлений.
type ItemLocationCache interface {
	Put(itemId store.ItemId, locations []store.Location) error
	Get(itemId store.ItemId) ([]store.Location, error)
	Invalidate(itemId store.ItemId) error
}
