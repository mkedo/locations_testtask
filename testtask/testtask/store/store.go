package store

type ItemId = int64
type LocationId = int64

type ItemLocations interface {
	Put(itemId ItemId, locationIds []LocationId) error
	Get(itemId ItemId) ([]Location, error)
}

//type Locations interface {
//	Lookup(str string) (Location, error)
//}

type Coordinates struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
}

type Location struct {
	ID          LocationId  `json:"ID"`
	Location    string      `json:"Location"`
	Coordinates Coordinates `json:"Coordinates"`
}
