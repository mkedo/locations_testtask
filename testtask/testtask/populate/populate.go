package populate

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"testtask/store"
	"time"
)

func GetRandomLocations(startIdx int64, endIndex int64) []store.Location {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if endIndex < startIdx {
		panic(errors.New("endIndex is less than startIdx"))
	}
	locationsNum := endIndex - startIdx + 1
	result := make([]store.Location, locationsNum)
	for i, _ := range result {
		id := startIdx + int64(i)
		addLen := 50 + rand.Int31n(755)
		result[i] = store.Location{
			ID:          id,
			Location:    "location " + fmt.Sprintf("%d", id) + randStringBytes(int(addLen)),
			Coordinates: randCoordinate(r),
		}
	}
	return result
}

func randCoordinate(r *rand.Rand) store.Coordinates {
	const roundTo = 10000
	x := math.Round(((r.Float64()*180*2)-180)*roundTo) / roundTo
	y := math.Round(((r.Float64()*90*2)-90)*roundTo) / roundTo
	return store.Coordinates{
		X: x,
		Y: y,
	}
}

func GetRandomLocationIds(locationRangeStart, locationRangeEnd int64) ([]store.LocationId) {
	if locationRangeEnd < locationRangeStart {
		panic(errors.New("locationRangeEnd is less than locationRangeStart"))
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	locationsPerItem := r.Int31n(7)
	rangeSize := int32(locationRangeEnd - locationRangeStart + 1)
	if locationsPerItem > rangeSize {
		locationsPerItem = rangeSize
	}

	locationIds := make([]store.LocationId, 0, locationsPerItem)
	generated := make(map[int64]bool)
	for i := int32(0); i < locationsPerItem; i++ {
		if locationsPerItem == rangeSize {
			id := locationRangeStart + int64(i)
			locationIds = append(locationIds, id)
		} else {
			for {
				id := locationRangeStart + r.Int63n(int64(rangeSize))
				if !generated[id] {
					generated[id] = true
					locationIds = append(locationIds, id)
					break
				}
			}
		}
	}
	return locationIds
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
