package main

import (
	"Yadier01/neon/cmd/api"
	db "Yadier01/neon/db/sqlc"
	"Yadier01/neon/util"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func main() {

	config, err := util.LoadConfig("./")
	if err != nil {
		log.Fatal("cannot load config ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(config, store)

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

}
