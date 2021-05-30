package aws

import (
	"database/sql"
	"fmt"

	"feed-service/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	var err error
	c := config.NewConfig()

	DB, err = sql.Open(c.Dialect, "postgres://"+c.Username+":"+c.Password+"@"+c.Host+"/"+c.Database+"?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}
