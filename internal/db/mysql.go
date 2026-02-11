package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func NewMySQL(dns string) *sqlx.DB {
	db, err := sqlx.Connect("mysql", dns)

	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db
}
