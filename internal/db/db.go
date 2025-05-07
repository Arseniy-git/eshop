package db

import (
	"database/sql"
	"eshop/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init(dsn string) {
	var err error
	DB, err = sql.Open("mysql", config.GetDSN())
	if err != nil {
		log.Fatal("DB open error:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}

	log.Println("Connected to DB")
}
