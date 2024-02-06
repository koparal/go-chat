package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Table    string `json:"table"`
}

type Database struct {
	DB         *sql.DB
	connString string
}

func New(dbConf Config) (*Database, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Table)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	return &Database{DB: db, connString: connString}, nil
}

func (d *Database) Close() {
	d.DB.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.DB
}
