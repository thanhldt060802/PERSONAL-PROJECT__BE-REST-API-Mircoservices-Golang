package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
	"thanhldt060802/config"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var DB *bun.DB

func InitPostgesConnection() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.AppConfig.DBUser, config.AppConfig.DBPassword, config.AppConfig.DBHost, config.AppConfig.DBPort, config.AppConfig.DBName,
	)

	pgdb, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Connect to PostgreSQL with Bun ORM failed: ", err)
	}

	DB = bun.NewDB(pgdb, pgdialect.New())

	if err := DB.Ping(); err != nil {
		log.Fatal("Ping to database failed: ", err)
	}
	log.Println("Connected to PostgreSQL with Bun ORM successful")
}
