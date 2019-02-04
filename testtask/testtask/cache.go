package testtask

import "testtask/store"

type ItemLocationCache interface {
	Put(itemId store.ItemId, locations []store.Location) error
	Get(itemId store.ItemId) ([]store.Location, error)
}
