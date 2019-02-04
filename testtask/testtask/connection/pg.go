package connection

import (
	"database/sql"
	"os"
)
import _ "github.com/lib/pq"

func GetPgConnection() *sql.DB {
	dsn, _ := os.LookupEnv("PG_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(4)
	//if err := db.Ping(); err != nil {
	//	panic(err)
	//}
	return db
}
