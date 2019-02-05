package testtask

import "testtask/store"

// Кеш, хранящий результаты получения Location для объявлений.
type ItemLocationCache interface {
	Put(itemId store.ItemId, locations []store.Location) error
	Get(itemId store.ItemId) ([]store.Location, error)
}
