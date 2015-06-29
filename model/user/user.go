/*
 * 'User' provides basic functionality to modify the user table.
 */
package user

import(
  _ "github.com/go-sql-driver/mysql"
  "AfkChampFrontend/model"
  "github.com/jmoiron/sqlx"
  "code.google.com/p/go.crypto/bcrypt"
  "errors"
)

type UserEntry struct {
  UserId int64
  Username string
  IsAdmin bool
  Password []byte
  Email string
}

const MinPasswordLength = 8
const MaxEmailLength = 254
const MaxUsernameLength= 50

const saltByteLength = 8
const passwordEncryptCost = 12

var UserExistsError = errors.New("Username Exists.")
var UserDoesNotExist = errors.New("Username/Password does not match.")
var UserUnspecifiedError = errors.New("Unknown Error.")
var UserLoggedInError = errors.New("User already logged in.")

// RetrieveUser takes a user ID and returns the user entry corresponding to that ID.
func RetrieveUser(userId int64) (*UserEntry, error) {
  row := model.Database.QueryRowx("SELECT * FROM users WHERE userid = ?", userId)
  userRow := UserEntry{}
  err := row.StructScan(&userRow)
  if err != nil {
    return nil, err
  }
  return &userRow, nil
}

// VerifyUser takes a username and password and checks to see if 1) the user exists and 2) the user's password matches.
func VerifyUser(username string, password string) (*UserEntry, error) {
  users, err := getUser(username)
  if err != nil {
    return nil, err 
  }
  
  // Get the actual user struct for access
  userStruct, err := extractUser(users)
  if err != nil {
    return nil, err
  }
  
  // User exists, make sure passwords match
  return userStruct, bcrypt.CompareHashAndPassword(userStruct.Password, []byte(password))
}

// CreateUser creates a new user entry in the table. Returns whether or not the user was able to be
// created successfully. If not, set the error.
func CreateUser(username string, password string, email string) (*UserEntry, error) {
  if users, nerr := getUser(username);  nerr != nil || users.Next() {
    if nerr != nil {
      return nil, nerr
    }
    return nil, UserExistsError
  }

  // Encrypt password
  encryptedPassword, nerr := bcrypt.GenerateFromPassword([]byte(password), passwordEncryptCost)
  if nerr != nil {
    return nil, nerr
  }
  
  // Send DB query to create the user.
  newUser := UserEntry{Username: username, Password: encryptedPassword, IsAdmin: false, Email: email}
  retRow, nerr := model.Database.Exec("INSERT INTO users (username, password, isadmin, email) VALUES (?, ?, ?, ?)", newUser.Username, newUser.Password, newUser.IsAdmin, newUser.Email)
  if nerr != nil {
    return nil, nerr
  }
  newUser.UserId, nerr = retRow.LastInsertId()
  if nerr != nil {
    return nil, nerr
  }
  return &newUser, nil
}

// UpdateUser takes in a user ID and updates entries accordingly. Note that we DO NOT update the password.
func UpdateUser(userId int64, userData *UserEntry) error {
  _ , nerr := model.Database.Exec("UPDATE users SET username = ?, isadmin = ?, email = ? WHERE userid = ?", userData.Username, userData.IsAdmin, userData.Email, userId)
  if nerr != nil {
    return nerr
  }
  return nil
}

// getUser executes a query to search for the specified user by username.
func getUser(username string) (*sqlx.Rows, error) {
  rows, err := model.Database.Queryx("SELECT * FROM users WHERE username = ?", username)
  if err != nil {
    return nil, err
  }
  return rows, nil
}

// 'IndexUsers' let's us return a list of users with the list having an arbitrary length and an arbitrary offset.
func IndexUsers(offset int, count int) ([]UserEntry, error) {
  requestUsers := make([]UserEntry, 0, 0)
  
  // Find the users
  rows, err := model.Database.Queryx("SELECT * FROM users ORDER BY userid ASC LIMIT ?, ?", offset, count)
  if err != nil {
    return nil, err
  }
  
  // Create the GameRow structs and then return it
  for rows.Next() {
    newObj := UserEntry{}
    err = rows.StructScan(&newObj)
    if err != nil {
      continue
    }
    requestUsers = append(requestUsers, newObj)
  }
  
  return requestUsers, nil
}

// 'DeleteUser' removes the user from the database
func DeleteUser(userid int64) error {
  _, err := model.Database.Exec("DELETE FROM users WHERE userid = ? AND isadmin = ?", userid, false)
  if err != nil {
    return err
  }
  return nil
}

// extractUser extracts a user from the result of 'getUser'. We asssume that there only ever exists one user with a given username.
func extractUser(inRows *sqlx.Rows) (*UserEntry, error) {
  newUser := UserEntry{}
  var err error
  if inRows.Next() {
    err = inRows.StructScan(&newUser)
  } else {
    return nil, UserUnspecifiedError
  }
  return &newUser, err
}