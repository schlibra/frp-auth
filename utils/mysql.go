package utils

import (
	"database/sql"
	"fmt"
	"frp-auth/model"
	"log"
)

func LoadMySQL(cfg model.Config) *sql.DB {
	username := cfg.MySQL.Username
	password := cfg.MySQL.Password
	host := cfg.MySQL.Host
	port := cfg.MySQL.Port
	database := cfg.MySQL.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
