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
  "log"
)

type UserEntry struct {
  UserId int64
  Username string
  IsAdmin bool
  Password []byte
}

const saltByteLength = 8
const passwordEncryptCost = 12

var UserExistsError = errors.New("Username Exists.")
var UserDoesNotExist = errors.New("Username/Password does not match.")
var UserUnspecifiedError = errors.New("Unknown Error.")

// VerifyUser takes a username and password and checks to see if 1) the user exists and 2) the user's password matches.
func VerifyUser(username string, password string) error {
  users, err := GetUser(username)
  if err != nil {
    return err 
  }
  
  // Get the actual user struct for access
  userStruct, err := ExtractUser(users)
  if err != nil {
    return err
  }
  
  // User exists, make sure passwords match
  return bcrypt.CompareHashAndPassword(userStruct.Password, []byte(password))
}

// CreateUser creates a new user entry in the table. Returns whether or not the user was able to be
// created successfully. If not, set the error.
func CreateUser(username string, password string) error {
  if users, nerr := GetUser(username);  nerr != nil || users.Next() {
    if nerr != nil {
      return nerr
    }
    return UserExistsError
  }

  // Encrypt password
  encryptedPassword, nerr :=  bcrypt.GenerateFromPassword([]byte(password), passwordEncryptCost)
  log.Print(encryptedPassword)
  if nerr != nil {
    return nerr
  }
  
  // Send DB query to create the user.
  newUser := UserEntry{Username: username, Password: encryptedPassword, IsAdmin: false}
  _ , nerr = model.Database.Exec("INSERT INTO users (username,password,isadmin) VALUES ($1, $2, $3);", newUser.Username, newUser.Password, newUser.IsAdmin)
  if nerr != nil {
    return nerr
  }
  return nil
}

// GetUser executes a query to search for the specified user by username.
func GetUser(username string) (*sqlx.Rows, error) {
  rows, err := model.Database.Queryx("SELECT * FROM users WHERE username = $1;", username)
  if err != nil {
    return nil, err
  }
  return rows, nil
}

// ExtractUser extracts a user from the result of 'GetUser'. We asssume that there only ever exists one user with a given username.
func ExtractUser(inRows *sqlx.Rows) (*UserEntry, error) {
  newUser := UserEntry{}
  var err error
  if inRows.Next() {
    err = inRows.StructScan(&newUser)
  } else {
    return nil, UserUnspecifiedError
  }
  return &newUser, err
}