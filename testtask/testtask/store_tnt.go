package testtask

import (
	"encoding/json"
	"github.com/tarantool/go-tarantool"
	"log"
	"testtask/store"
)

type TntStore struct {
	con *tarantool.Connection
}

func NewTntStore(con *tarantool.Connection) *TntStore {
	return &TntStore{
		con: con,
	}
}

func (s *TntStore) Put(itemId store.ItemId, locationIds []store.LocationId) error {
	var ids = []int64(locationIds)
	_, err := s.con.CallAsync("put_item_locations", []interface{}{
		itemId,
		ids,
	}).Get()
	if err != nil {
		return err
	}
	return nil
}

func (s *TntStore) Get(itemId store.ItemId) ([]store.Location, error) {
	var locations = make([]store.Location, 0)
	var locationsStrArray [][]string
	err := s.con.CallTyped("get_item_locations", []interface{}{itemId}, &locationsStrArray)
	if err != nil {
		return nil, err
	}
	for _, row := range locationsStrArray[0] {
		var location store.Location

		//TODO: проверить что json правильной структуры
		if err := json.Unmarshal([]byte(row), &location); err == nil {
			locations = append(locations, location)
		} else {
			log.Println(err)
		}
	}
	return locations, nil
}

func (s *TntStore) Add(locations []store.Location) error {
	if len(locations) == 0 {
		return nil
	}
	bulk := make([]interface{}, len(locations))
	for i := 0; i < len(locations); i++ {
		data, err := json.Marshal(locations[i])
		if err != nil {
			return err
		}
		bulk[i] = []interface{}{
			locations[i].ID,
			data,
		}
	}
	// { {location_id, location_data}, ... }
	_, err := s.con.CallAsync("add_locations", []interface{}{bulk}).Get()
	if err != nil {
		return err
	}
	return nil
}
