var LiveStats = {};

////////////////////////////
// Global Stats Parsing
////////////////////////////
LiveStats.GetTime = function(data) {
    return data.global.time;
}

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

LiveStats.GetTeamName = function(team) {
    return team.fullName;
}

////////////////////////////
// Player Parsing
////////////////////////////
LiveStats.GetPlayerFromTeam = function(team, playerIndex) {
    return team.players[playerIndex];
}

LiveStats.GetPlayerName = function(player) {
    if (!player.name || player.name == "") {
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

////////////////////////////
// Draft Parsing
////////////////////////////
LiveStats.GetPicksForTeam = function(data, teamIdx) {
    return data.picks[teamIdx];
}

LiveStats.GetBansForTeam = function(data, teamIdx) {
    return data.bans[teamIdx];
}

////////////////////////////
// Miscellaneous Parsing
////////////////////////////
LiveStats.GetIsDraftMode = function(data) {
    return (data.mode == 1);
}

LiveStats.GetIsGameMode = function(data) {
    return (data.mode == 0);
}