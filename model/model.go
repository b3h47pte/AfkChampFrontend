/*
 * 'Model' provides basic database functionality. We force ourselves to use PostgreSQL.
 */
package model

import(
  "log"
  _ "github.com/go-sql-driver/mysql"
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

// 'init' sets up a connection with all the databases this program will need.
func init() {
  // Main database
  var config DatabaseConfig
  err := gcfg.ReadFileInto(&config, "config/db.config")
  if err != nil {
    log.Fatal(err)
  }
  
  // Connect
  fullDBUrl := config.Database.DatabaseUsername + ":" + config.Database.DatabasePassword +
    "@tcp(" + config.Database.DatabaseUrl + ":" + config.Database.DatabasePort + ")/" + config.Database.DatabaseName
  if config.Database.DatabaseUseSSL {
    fullDBUrl += "?tls=true"
  }
  // Need date time support
  fullDBUrl += "?parseTime=true"
  
  Database, err = sqlx.Connect("mysql", fullDBUrl)
  if err != nil {
    log.Fatal(err)
  }
}