package config

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

// DB holds the database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
    dsn := "paraco_user:ValueLeafBang1!@tcp(127.0.0.1:3306)/paraco_db"
    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal("Could not connect to the database:", err)
    }

    log.Println("Successfully connected to the database!")
}
