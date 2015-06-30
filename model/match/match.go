/*
 * 'match' provides basic functionality to perform basic SQL operations on the 'matches' database as well as the 
 * various other databases that are joined with matches (i.e. match_team and match_players).
 */
package match

import (
	"AfkChampFrontend/model"
	_ "github.com/go-sql-driver/mysql"
    "time"
    "strings"
    "errors"
)

type MatchRow struct {
    MatchId int
    ParentSeriesId int
    MatchTime time.Time
    IsGameFinished bool
}

// Information needed by the match controller
// to display the initial page.
// Team Data (Names), Series Data, Date
// Note that "teams" will have the team names separated by a comma.
type BasicMatchInformationScanned struct {
    Teams string
    CurrentGame int
    TotalGames int
    MatchDate time.Time
}

type BasicMatchInformationReturn struct {
    TeamOne string
    TeamTwo string
    CurrentGame int
    TotalGames int
    MatchDate time.Time
}

func QueryBasicInformation(id int64) (*BasicMatchInformationReturn, error) {
	// Find the ONE match.
	rows, err := model.Database.Queryx(`SELECT GROUP_CONCAT(teams.teamshorthand) AS teams,
                                            matches.matchnum AS currentgame,
                                            series.bestoutof AS totalgames,
                                            matches.matchtime AS matchdate
                                        FROM matches
                                            INNER JOIN series ON matches.parentseriesid = seriesid
                                            INNER JOIN match_team ON match_team.matchid = matches.matchid
                                            INNER JOIN teams ON match_team.teamid = teams.teamid
                                        WHERE matches.matchid = ?`, id);
	if err != nil {
		return nil, err
	}

	// Create structs
	for rows.Next() {
		newObj := BasicMatchInformationScanned{}
		err = rows.StructScan(&newObj)
		if err != nil {
			return nil, err
		}
        
        // Split team name by comma.
        teams := strings.Split(newObj.Teams, ",")
        
        retObj := BasicMatchInformationReturn {
            TeamOne: teams[0],
            TeamTwo: teams[1],
            CurrentGame: newObj.CurrentGame,
            TotalGames: newObj.TotalGames,
            MatchDate: newObj.MatchDate,
        }
        return &retObj, nil
	}

    // Shouldn't get unless no match is found.
	return nil, errors.New("No corresponding match ID found")
}