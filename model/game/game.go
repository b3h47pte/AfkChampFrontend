/*
 * 'game' provides basic functionality to modify the game table. Nothing fancy here besides some get/set/modify SQL queries.
 */
package game

import(
  _ "github.com/go-sql-driver/mysql"
  "AfkChampFrontend/model"
)

type GameRow struct {
  GameId int
  GameName string
  // Used in the URL
  GameShortName string
}

// 'GetGames' returns a portion of the available games. We are provided with a range of games to get assuming the games
// are sorted by game ID (number index).
func GetGames(offset int, count int) ([]GameRow, error) {
  requestGames := make([]GameRow, count, count)
  
  // Find the games
  rows, err := model.Database.Queryx("SELECT * FROM games LIMIT ?, ?", offset, count)
  if err != nil {
    return nil, err
  }
  
  // Create the GameRow structs and then return it
  objectCount := 0
  for rows.Next() {
    newObj := GameRow{}
    err = rows.StructScan(&newObj)
    if err != nil {
      continue
    }
    requestGames[objectCount] = newObj
    objectCount++
  }
  
  return requestGames, nil
}
