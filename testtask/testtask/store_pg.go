package testtask

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"strings"
	"sync"
	"testtask/store"
	"time"
)

const dbSchema = "testtask"

type PgStore struct {
	db          *sql.DB
	insertMutex *sync.Mutex
}

// Хранилище в БД postgres.
// Адреса и связь с объявлениями хранятся в самой БД.
func NewPgStore(db *sql.DB) *PgStore {
	return &PgStore{
		db:          db,
		insertMutex: &sync.Mutex{},
	}
}

func (s *PgStore) PutContext(ctx context.Context, itemId store.ItemId, locationIds []store.LocationId) error {
	var insertSql string
	var valueArgs []interface{}
	hasInsert := len(locationIds) > 0

	if hasInsert {
		valueStrings := make([]string, 0, len(locationIds))
		valueArgs = make([]interface{}, 0, len(locationIds)*2)
		paramCount := 1
		for _, locationId := range locationIds {
			valueStrings = append(
				valueStrings,
				fmt.Sprintf("($%d, $%d)", paramCount, paramCount+1))
			paramCount = paramCount + 2
			valueArgs = append(valueArgs, itemId)
			valueArgs = append(valueArgs, locationId)
		}
		insertSql = fmt.Sprintf(
			`INSERT INTO %s.item_locations (%s, %s)
				VALUES %s
				ON CONFLICT (item_id,location_id) DO NOTHING`,
			pq.QuoteIdentifier(dbSchema),
			pq.QuoteIdentifier("item_id"),
			pq.QuoteIdentifier("location_id"),
			strings.Join(valueStrings, ","))
	}

	deleteSql := fmt.Sprintf(
		"DELETE FROM %s.item_locations WHERE %s = $1",
		pq.QuoteIdentifier(dbSchema),
		pq.QuoteIdentifier("item_id"))

	// записывам однопоточно, чтобы избежать serialize error
	s.insertMutex.Lock()
	defer s.insertMutex.Unlock()
	err := retryOnSerializeError(3, func() error {
		return transact(s.db, func(tx *sql.Tx) error {
			_, err := tx.Exec(deleteSql, itemId)
			if err != nil {
				log.Println(err)
				return err
			}
			if hasInsert {
				_, err = tx.Exec(insertSql, valueArgs...)
				if err != nil {
					log.Println(err)
					return err
				}
			}
			return nil
		})
	})
	return err
}

func retryOnSerializeError(maxRetry int, query func() error) error {
	if maxRetry <= 0 {
		panic("maxRetry must be positive non-zero integer")
	}
	var err error
	for i := 0; i < maxRetry; i++ {
		err = query()
		if err != nil {
			if pqerr, ok := err.(*pq.Error); ok {
				if pqerr.Code.Name() == "serialization_failure" {
					continue
				} else {
					return err
				}
			}
			return err
		}
		break
	}
	return err
}

func (s *PgStore) Add(ctx context.Context, locations []store.Location) error {
	if len(locations) == 0 {
		return nil
	}
	valueStrings := make([]string, 0, len(locations))
	valueArgs := make([]interface{}, 0, len(locations)*3)
	paramCount := 1
	for _, location := range locations {
		valueStrings = append(
			valueStrings,
			fmt.Sprintf("($%d, $%d, $%d)", paramCount, paramCount+1, paramCount+2))
		paramCount = paramCount + 3
		valueArgs = append(valueArgs, location.ID)
		valueArgs = append(valueArgs, location.Location)
		valueArgs = append(valueArgs, pgCoordinates(location.Coordinates))
	}
	insertSql := fmt.Sprintf(`
		INSERT INTO %s.locations (location_id, location, coordinates)
		VALUES %s
		`,
		pq.QuoteIdentifier(dbSchema),
		strings.Join(valueStrings, ","))
	return transact(s.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(insertSql, valueArgs...)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
}

func (s *PgStore) GetContext(ctx context.Context, itemId store.ItemId) ([]store.Location, error) {
	selectSql := fmt.Sprintf(`
		SELECT 
			l.location_id, l.location, l.coordinates::varchar
		FROM 
			%s.item_locations loc
		INNER JOIN 
			%s.locations l ON (l."location_id" = loc."location_id")
		WHERE 
			loc.item_id = $1
	`,
		pq.QuoteIdentifier(dbSchema),
		pq.QuoteIdentifier(dbSchema))

	context, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	//context := context.Background()
	rows, err := s.db.QueryContext(context, selectSql, itemId)
	if err != nil {
		log.Println(err)
		return []store.Location{}, err
	}
	defer rows.Close()
	return FetchLocations(rows)
}

// Конвертирует *sql.Rows в []store.Location.
// В запросе, первыми 3 полями, должны быть выбраны следующие поля (порядок важен):
//     location_id, location, coordinates::varchar
func FetchLocations(rows *sql.Rows) ([]store.Location, error) {
	locations := make([]store.Location, 0)
	for rows.Next() {
		var location store.Location
		var rawCoordinates pgCoordinates
		err := rows.Scan(&location.ID, &location.Location, &rawCoordinates)
		if err != nil {
			log.Println(err)
			return []store.Location{}, err
		} else {
			location.Coordinates = store.Coordinates(rawCoordinates)
			locations = append(locations, location)
		}
	}
	if err := rows.Err(); err != nil {
		return locations, err
	}
	return locations, nil
}

// Выполняет функцию txFunc внутри транзакции.
// Откатывает транзакцию в случае ошибки, коммитит в случае успеха.
func transact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}

// Промежуточная структура чтобы конвертировать координату
// в формат понятный postgres и обратно.
type pgCoordinates store.Coordinates

func (c pgCoordinates) Value() (driver.Value, error) {
	return fmt.Sprintf("(%f,%f)", c.X, c.Y), nil
}

func (c *pgCoordinates) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	default:
		return errors.New("Incompatible type for pgCoordinates")
	}
	var x, y float64
	_, err := fmt.Sscanf(source, "(%e,%e)", &x, &y)
	if err != nil {
		return err
	}
	*c = pgCoordinates{
		X: x,
		Y: y,
	}
	return nil
}
