var LiveStats = {};

////////////////////////////
// Team Parsing
////////////////////////////
LiveStats.GetTeam = function(data, index) {
    return data.teams[index];
}

LiveStats.GetTeamGold = function(team) {
    return team.gold;
}

LiveStats.GetTeamTowers = function(team) {
    return team.towers;
}

LiveStats.GetTeamCurrentDragons = function(team) {
    return team.currentDragons;
}

LiveStats.GetTeamBarons = function(team) {
    return team.barons;
}

LiveStats.GetTeamKills = function(team) {
    return team.kills;
}

LiveStats.GetTeamSeriesWins = function(team) {
    return team.series;   
}

////////////////////////////
// Player Parsing
////////////////////////////
LiveStats.GetPlayerFromTeam = function(team, playerIndex) {
    return team.players[playerIndex];
}

LiveStats.GetPlayerName = function(player) {
    if (player.name == "") {
        return "Unknown";
    }
    return player.name;
}

LiveStats.GetPlayerChampionName = function(player) {
    return player.champion;
}

LiveStats.GetPlayerOverviewStats = function(player) {
    var statObject = {
        kills: player.kills,
        deaths: player.deaths,
        assists: player.assists,
        creeps: player.creeps
    };
    return statObject;
}