package connection

import (
	"github.com/tarantool/go-tarantool"
	"os"
	"time"
)

func GetTntConnection() *tarantool.Connection {
	user, ok := os.LookupEnv("TARANTOOL_USER")
	if !ok {
		user = "guest"
	}
	addr, _ := os.LookupEnv("TARANTOOL_ADDR")

	var opts = tarantool.Opts{
		User:          user,
		Timeout:       500 * time.Millisecond,
		MaxReconnects: 1,
	}
	client, err := tarantool.Connect(addr, opts)
	if err != nil {
		panic(err)
	}
	//_, err = client.Ping()
	//if err != nil {
	//	panic(err)
	//}
	return client
}
