/*
 * 'User' provides basic functionality to modify the user table.
 */
package user

import(
  _ "github.com/lib/pq"
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

const saltByteLength = 8
const passwordEncryptCost = 12

var UserExistsError = errors.New("Username Exists.")
var UserDoesNotExist = errors.New("Username/Password does not match.")
var UserUnspecifiedError = errors.New("Unknown Error.")

// VerifyUser takes a username and password and checks to see if 1) the user exists and 2) the user's password matches.
func VerifyUser(username string, password string) error {
  users, err := getUser(username)
  if err != nil {
    return err 
  }
  
  // Get the actual user struct for access
  userStruct, err := extractUser(users)
  if err != nil {
    return err
  }
  
  // User exists, make sure passwords match
  return bcrypt.CompareHashAndPassword(userStruct.Password, []byte(password))
}

// CreateUser creates a new user entry in the table. Returns whether or not the user was able to be
// created successfully. If not, set the error.
func CreateUser(username string, password string, email string) error {
  if users, nerr := getUser(username);  nerr != nil || users.Next() {
    if nerr != nil {
      return nerr
    }
    return UserExistsError
  }

  // Encrypt password
  encryptedPassword, nerr :=  bcrypt.GenerateFromPassword([]byte(password), passwordEncryptCost)
  if nerr != nil {
    return nerr
  }
  
  // Send DB query to create the user.
  newUser := UserEntry{Username: username, Password: encryptedPassword, IsAdmin: false, Email: email}
  _ , nerr = model.Database.Exec("INSERT INTO users (username, password, isadmin, email) VALUES ($1, $2, $3, $4);", newUser.Username, newUser.Password, newUser.IsAdmin, newUser.Email)
  if nerr != nil {
    return nerr
  }
  return nil
}

// getUser executes a query to search for the specified user by username.
func getUser(username string) (*sqlx.Rows, error) {
  rows, err := model.Database.Queryx("SELECT * FROM users WHERE username = $1;", username)
  if err != nil {
    return nil, err
  }
  return rows, nil
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