/*
 * 'Model' provides basic database functionality. We force ourselves to use PostgreSQL.
 */
package model

import(
  "log"
  _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"
  "code.google.com/p/gcfg"
)

type DatabaseConfig struct {
  Database struct {
    DatabaseUrl string
    DatabasePort string
    DatabaseUsername string
    DatabasePassword string
    DatabaseName string
    DatabaseUseSSL bool
  }
}

var Database *sqlx.DB

// 'InitializeDatabase' sets up a connection with all the databases this program will need.
func InitializeDatabase() {
  // Main database
  var config DatabaseConfig
  err := gcfg.ReadFileInto(&config, "config/db.config")
  if err != nil {
    log.Fatal(err)
  }
  
  // Connect
  fullDBUrl := "postgres://" + config.Database.DatabaseUsername + ":" + config.Database.DatabasePassword +
    "@" + config.Database.DatabaseUrl + ":" + config.Database.DatabasePort + "/" + config.Database.DatabaseName
  if config.Database.DatabaseUseSSL {
    fullDBUrl += "?sslmode=verify-full"
  } else {
    fullDBUrl += "?sslmode=disable"
  }
  
  Database, err = sqlx.Connect("postgres", fullDBUrl)
  if err != nil {
    log.Fatal(err)
  }
}