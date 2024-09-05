package app

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitMySQL() {
	db = initMySQL()
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
}

func GetConnection() *sql.DB {
	return db
}

func initMySQL() *sql.DB {
	fmt.Println("Initializing connection to MySQL...")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting to MySQL...")
	err = conn.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to MySQL")

	return conn
}
