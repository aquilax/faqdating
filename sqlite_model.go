package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteModel struct {
	db *sqlx.DB
}

func NewSQLiteModel() *SQLiteModel {
	return &SQLiteModel{}
}

func (m *SQLiteModel) Init(config *Config) error {
	var err error
	m.db, err = sqlx.Connect(config.Database, config.Dsn)
	return err
}
