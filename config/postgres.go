package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func MustPostgres() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost",
		5432,
		"postgres",
		"password",
		"golang",
	)

	return sql.Open("postgres", psqlInfo)
}
