/*
 * 'Model' provides basic database functionality. We force ourselves to use PostgreSQL.
 */
package model

import(
  "log"
  _ "github.com/lib/pq"
  "database/sql"
  "code.google.com/p/gcfg"
)

type DatabaseConfig struct {
  Database struct {
    DatabaseUrl string
    DatabaseUsername string
    DatabasePassword string
    DatabaseName string
  }
}

var Database *sql.DB

// 'InitializeDatabase' sets up a connection with all the databases this program will need.
func InitializeDatabase() {
  // Main database
  var config DatabaseConfig
  err := gcfg.ReadFileInto(&config, "config/db.config")
  if err != nil {
    log.Fatal(err)
  }
  
  // Connect
  Database, err = sql.Open("postgres", "postgres://" + config.Database.DatabaseUsername + ":" + config.Database.DatabasePassword +
    "@" + config.Database.DatabaseUrl + "/" + config.Database.DatabaseName + "?sslmode=verify-full")
  if err != nil {
    log.Fatal(err)
  }
}