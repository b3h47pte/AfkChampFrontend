var LiveStatsUtility = {};

////////////////////////////
// Team Utility
////////////////////////////
LiveStatsUtility.GetTeamImage = function(teamData) {
    return "/images/teams/cloud9.png";
}

////////////////////////////
// Player Utility
////////////////////////////
LiveStatsUtility.GetPlayerProfilePicture = function(playerName) {
    return "/images/players/c9-balls.jpg";
}

////////////////////////////
// Champion Utility
////////////////////////////
LiveStatsUtility.GetChampionProfilePicture = function(championName) {
    return "/images/champions/" + championName + "_Square_0.png";
}

LiveStatsUtility.GetLargeChampionProfilePicture = function(championName) {
    return "/images/champions/" + championName + "_0.jpg";
}

////////////////////////////
// Draft Utility
////////////////////////////
LiveStatsUtility.ConstructPlayerPickItem = function(pick, player) {
    return {
        champ: LiveStatsUtility.GetLargeChampionProfilePicture(pick),
        player: LiveStatsUtility.GetPlayerProfilePicture(player),
        playerName: LiveStats.GetPlayerName(player)
    };
}

LiveStatsUtility.ConstructBanItem = function(ban) {
    return LiveStatsUtility.GetChampionProfilePicture(ban);
}
