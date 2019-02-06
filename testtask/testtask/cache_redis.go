package testtask

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"testtask/store"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

// Промежуточная структура для сериализации/десереализации.
type locationCollection struct {
	Locations []store.Location
}

// Кеш на Redis, хранящий результаты получения Location для объявлений.
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) Put(itemId store.ItemId, locations []store.Location) error {
	key := formatCacheKey(itemId)
	err := r.client.Set(key, locationCollection{Locations: locations}, 1*time.Hour).Err()
	return err
}

func (r *RedisCache) Get(itemId store.ItemId) ([]store.Location, error) {
	key := formatCacheKey(itemId)
	var locations locationCollection
	err := r.client.Get(key).Scan(&locations)
	if err == redis.Nil {
		return nil, ItemLocationCacheMiss
	} else if err != nil {
		log.Println(err)
		return nil, err
	} else {
		return locations.Locations, nil
	}
}

func (r *RedisCache) Invalidate(itemId store.ItemId) error {
	key := formatCacheKey(itemId)
	return r.client.Del(key).Err()
}

func (l locationCollection) MarshalBinary() ([]byte, error) {
	return json.Marshal(l)
}

func (l *locationCollection) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, l); err != nil {
		return err
	}
	return nil
}

// Формирует имя ключа который указывает где хранится массив Location
// для указанного объявления.
func formatCacheKey(itemId store.ItemId) string {
	return fmt.Sprintf("item:%v:locations", itemId)
}
