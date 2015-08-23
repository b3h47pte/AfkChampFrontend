var LiveStats = {};

////////////////////////////
// Team Parsing
////////////////////////////
LiveStats.GetTeam = function(data, index) {
    return data.teams[index];
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