package testtask

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/lib/pq"
	"log"
	"strings"
	"testtask/store"
	"time"
)

type PgRedisPersistent struct {
	db    *sql.DB
	redis *redis.Client
}

// Промежуточная структура для сериализации/десереализации идентификаторов адресов
type redisLocationIds struct {
	ids []store.LocationId
}

func (r redisLocationIds) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r.ids)
}
func (r *redisLocationIds) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &r.ids)
}

func NewPgRedisPersistent(db *sql.DB, redis *redis.Client) *PgRedisPersistent {
	return &PgRedisPersistent{
		db:    db,
		redis: redis,
	}
}

func (s *PgRedisPersistent) PutContext(ctx context.Context, itemId store.ItemId, locationIds []store.LocationId) error {
	return s.redis.Set(formatKey(itemId), redisLocationIds{locationIds}, 0).Err()
}

func (s *PgRedisPersistent) GetContext(ctx context.Context, itemId store.ItemId) ([]store.Location, error) {
	var locations = make([]store.Location, 0)
	var redisLocationIds redisLocationIds
	cmd := s.redis.Get(formatKey(itemId))
	err := cmd.Err()
	if err == redis.Nil {
		return locations, nil
	}
	if err != nil {
		return locations, err
	}
	if err := cmd.Scan(&redisLocationIds); err != nil {
		return locations, err
	}
	if len(redisLocationIds.ids) == 0 {
		return locations, nil
	}

	phStrings := make([]string, 0, len(redisLocationIds.ids))
	phValues := make([]interface{}, 0, len(redisLocationIds.ids))
	for i, id := range redisLocationIds.ids {
		phStrings = append(phStrings, fmt.Sprintf("$%d", i+1))
		phValues = append(phValues, id)
	}

	selectSql := fmt.Sprintf(`
		SELECT 
			l.location_id, l.location, l.coordinates::varchar
		FROM 
			%s.locations l
		WHERE 
			l.location_id in (%s)
	`,
		pq.QuoteIdentifier(dbSchema),
		strings.Join(phStrings, ","))

	context, cancel := context.WithTimeout(ctx, 1000 * time.Millisecond)
	defer cancel()

	rows, err := s.db.QueryContext(context, selectSql, phValues...)
	if err != nil {
		log.Println(err)
		return locations, err
	}
	defer rows.Close()
	return FetchLocations(rows)
}

// Формирует имя ключа который указывает где хранятся идентификаторы адресов
// для указанного объявления.
func formatKey(itemId store.ItemId) string {
	return fmt.Sprintf("item:%v:location_ids", itemId)
}
