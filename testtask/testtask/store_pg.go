package testtask

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"strings"
	"testtask/store"
)

const dbSchema = "testtask"

type PgStore struct {
	db *sql.DB
}

func NewPgStore(db *sql.DB) *PgStore {
	return &PgStore{
		db: db,
	}
}

func (s *PgStore) Put(itemId store.ItemId, locationIds []store.LocationId) error {
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

	return transact(s.db, func(tx *sql.Tx) error {
		_, err := s.db.Exec(deleteSql, itemId)
		if err != nil {
			log.Println(err)
			return err
		}
		if hasInsert {
			_, err = s.db.Exec(insertSql, valueArgs...)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	})
}

func (s *PgStore) Add(locations []store.Location) error {
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
		valueArgs = append(valueArgs, formatCoordinates(&location.Coordinates))
	}
	insertSql := fmt.Sprintf(`
		INSERT INTO %s.locations (location_id, location, coordinates)
		VALUES %s
		`,
		pq.QuoteIdentifier(dbSchema),
		strings.Join(valueStrings, ","))
	return transact(s.db, func(tx *sql.Tx) error {
		_, err := s.db.Exec(insertSql, valueArgs...)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
}

func (s *PgStore) Get(itemId store.ItemId) ([]store.Location, error) {
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


	rows, err := s.db.Query(selectSql, itemId)
	if err != nil {
		log.Println(err)
		return []store.Location{}, err
	}
	defer rows.Close()
	var locations = make([]store.Location, 0)
	for rows.Next() {
		var location store.Location
		var rawCoordinates string
		err := rows.Scan(&location.ID, &location.Location, &rawCoordinates)
		if err != nil {
			log.Println(err)
			return []store.Location{}, err
		} else {
			coordinates, err := scanCoordinates(rawCoordinates)
			if err == nil {
				location.Coordinates = coordinates
				locations = append(locations, location)
			}

		}
	}
	return locations, nil
}

func transact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_ ,err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
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

func scanCoordinates(str string) (store.Coordinates, error) {
	var x, y float64
	_, err := fmt.Sscanf(str, "(%e,%e)", &x, &y)
	if err != nil {
		return store.Coordinates{}, err
	}
	return store.Coordinates{
		X: x,
		Y: y,
	}, nil
}
func formatCoordinates(coordinates *store.Coordinates) string {
	return fmt.Sprintf("(%f,%f)", coordinates.X, coordinates.Y)
}
