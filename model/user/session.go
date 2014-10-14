/*
 * 'Session' provides basic functionality to keep track of user sessions.
 */
package user

import(
  _ "github.com/go-sql-driver/mysql"
  "AfkChampFrontend/model"
  "github.com/jmoiron/sqlx"
  "time"
  "errors"
)

type SessionEntry struct {
  SessionId string
  UserId int64
  Expiration time.Time
}


var NoSessionError = errors.New("No session exists.")

// 'AddSessionForUser' grants the given user a new session key with the specified expiration date.
func AddSessionForUser(inUser *UserEntry, sessionKey string, expirationTime *time.Time) error {
  // Just add the key. If for whatever reason an error occurs, oh well.
  const dateTimeLayout = "2006-01-02 15:04:05"
  _ , err := model.Database.Exec("INSERT INTO sessions (sessionid, userid, expiration) VALUES (?, ?, ?)", sessionKey, inUser.UserId, expirationTime.Format(dateTimeLayout))
  if err != nil {
    return err
  }
  return nil
}

// 'VerifySession' takes in a session key and checks to see if it matches any of the user's session keys. If successful, we return the user ID.
func VerifySession(sessionKey string) (int64, error) {
  row := model.Database.QueryRowx("SELECT * FROM sessions WHERE sessionid = ?", sessionKey)
  sessionRow := SessionEntry{}
  err := row.StructScan(&sessionRow)
  if err != nil {
    return -1, err
  }
  return sessionRow.UserId, nil
}

// 'RemoveSessionForUser' takes in a session key attempts to remove the key.
func RemoveSessionForUser(sessionKey string) error {
  _, err := model.Database.Exec("DELETE FROM sessions WHERE sessionid = ?", sessionKey)
  if err != nil {
    return err
  }
  return nil
}

func getSession(sessionKey string) (*sqlx.Rows, error) {
  rows, err := model.Database.Queryx("SELECT * FROM sessions WHERE sessionid = ?", sessionKey)
  if err != nil {
    return nil, err
  }
  return rows, nil
}