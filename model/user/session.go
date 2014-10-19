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
  IsAdmin bool
}


var NoSessionError = errors.New("No session exists.")

// 'AddSessionForUser' grants the given user a new session key with the specified expiration date.
func AddSessionForUser(inUser *UserEntry, sessionKey string, expirationTime *time.Time, isAdmin bool) error {
  // Just add the key. If for whatever reason an error occurs, oh well.
  const dateTimeLayout = "2006-01-02 15:04:05"
  _ , err := model.Database.Exec("INSERT INTO sessions (sessionid, userid, expiration, isadmin) VALUES (?, ?, ?, ?)", sessionKey, inUser.UserId, expirationTime.Format(dateTimeLayout), isAdmin)
  if err != nil {
    return err
  }
  return nil
}

// 'VerifySession' takes in a session key and checks to see if it matches any of the user's session keys. If successful, we return nil
// otherwise we return an appropriate error.
func VerifySession(sessionKey string, userId int64, isAdmin bool) error {
  row := model.Database.QueryRowx("SELECT * FROM sessions WHERE sessionid = ? AND userid = ? AND isadmin = ?", 
    sessionKey, userId, isAdmin)
  sessionRow := SessionEntry{}
  err := row.StructScan(&sessionRow)
  if err != nil {
    return err
  }
  
  // If the session is expired, delete it and fail.
  currentTime := time.Now()
  if sessionRow.Expiration.Before(currentTime) {
    RemoveSessionForUser(sessionKey, userId, isAdmin)
    return NoSessionError
  }
  
  return nil
}

// 'RemoveSessionForUser' takes in a session key and attempts to remove the key.
func RemoveSessionForUser(sessionKey string, userId int64, isAdmin bool) error {
  _, err := model.Database.Exec("DELETE FROM sessions WHERE sessionid = ? AND userid = ? AND isadmin = ?", sessionKey, userId, isAdmin)
  if err != nil {
    return err
  }
  return nil
}

func getSession(sessionKey string, userId int64, isAdmin bool) (*sqlx.Rows, error) {
  rows, err := model.Database.Queryx("SELECT * FROM sessions WHERE sessionid = ? AND userid = ? AND isadmin = ?", 
    sessionKey, userId,isAdmin)
  if err != nil {
    return nil, err
  }
  return rows, nil
}