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