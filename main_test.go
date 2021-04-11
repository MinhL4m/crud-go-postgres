package main_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/MinhL4m/crud-go-postgres"
)

// This is in another package
// need to import "github.com/MinhL4m/crud-go-postgres" to work
var a main.App

func TestMain(m *testing.M) {
	err := godotenv.Load(".env")
	if err != nil {
		os.Exit(1)
	}

	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)`
