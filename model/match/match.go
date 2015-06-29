/*
 * 'match' provides basic functionality to perform basic SQL operations on the 'matches' database as well as the 
 * various other databases that are joined with matches (i.e. match_team and match_players).
 */
package model

import (
	"AfkChampFrontend/model"
	_ "github.com/go-sql-driver/mysql"
)