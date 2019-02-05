package store

import "context"

type ItemId = int64
type LocationId = int64

// Хранилище привязанных к объявлению адресов.
type ItemLocations interface {
	PutContext(ctx context.Context, itemId ItemId, locationIds []LocationId) error
	GetContext(ctx context.Context, itemId ItemId) ([]Location, error)
}

//type Locations interface {
//	Lookup(str string) (Location, error)
//}

type Coordinates struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
}

type Location struct {
	ID          LocationId  `json:"ID"`          // идентификатор локации/адреса
	Location    string      `json:"Location"`    // адрес
	Coordinates Coordinates `json:"Coordinates"` // координаты связанные с этим адресом
}
