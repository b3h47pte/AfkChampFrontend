rocketelo.factory('socketIOService', function($rootScope) {
    var socketService = {}
    socketService.isInit = false;
    socketService.Initialize = function(url, matchId) {
        if (socketService.isInit) {
            return;
        }
        
        // TODO: Get initial data.
        
        socketService.isInit = true;
        socketService.handlers = [];
        socketService.socket = io.connect(url);
        socketService.socket.emit('identify', {match: matchId.toString()});
        socketService.socket.on('liveupdate', socketService.ReceiveLiveStatsData);
    }
    
    socketService.ReceiveLiveStatsData = function(data) {
        for(var i = 0; i < socketService.handlers.length; ++i) {
            socketService.handlers[i](data);
        }
    }
    
    socketService.RegisterHandler = function(handler) {
        socketService.handlers.push(handler);
    }
    
    return socketService;
});

rocketelo.factory('compositeStatsService', function ($rootScope) {
    var compositeStatsService = {};
    compositeStatsService.teamStats = [[null, null, null, null, null], [null, null, null, null, null]];
    compositeStatsService.listeners = [];
    
    compositeStatsService.SetStatsForPlayerOnTeam = function(playerIdx, teamIdx, inStats) {
        if (playerIdx < 0 || playerIdx >= 5 || teamIdx < 0 || teamIdx >= 2) {
            return;
        }
        compositeStatsService.teamStats[teamIdx][playerIdx] = inStats;
    }
    
    compositeStatsService.RegisterListener = function(listener) {
        compositeStatsService.listeners.push(listener);
    }
    
    compositeStatsService.ComputeCompositeStats = function() {
        var compositeStats = {
            kills: 0,
            deaths: 0,
            assists: 0,
            creeps: 0
        };
        for (var t = 0; t < 2; ++t) {
            for (var p = 0; p < 5; ++p) {
                if (!compositeStatsService.teamStats[t][p]) {
                    continue;
                }
                
                // The constant numbers are to make it so that the bars don't start off topped off at the beginning of the game (aka there's room to grow until mid game-ish)
                compositeStats.kills = Math.min(Math.max(compositeStats.kills, compositeStatsService.teamStats[t][p].kills), 5);
                compositeStats.deaths = Math.min(Math.max(compositeStats.deaths, compositeStatsService.teamStats[t][p].deaths), 3);
                compositeStats.assists = Math.min(Math.max(compositeStats.assists, compositeStatsService.teamStats[t][p].assists), 5);
                compositeStats.creeps = Math.min(Math.max(compositeStats.creeps, compositeStatsService.teamStats[t][p].creeps), 50);
            }
        }
        
        for (var i = 0; i < compositeStatsService.listeners.length; ++i) {
            compositeStatsService.listeners[i](compositeStats);
        }
    }
    
    return compositeStatsService;
});

rocketelo.config(function($routeProvider) {
  $routeProvider.when('/draft', {
    templateUrl : '/partials/match/draft.html',
    controller: 'DraftBanController'
  }).when('/match', {
    templateUrl : '/partials/match/match.html',
    controller: 'MatchController'
  }).when('/post', {
    templateUrl : '/partials/match/post.html',
    controller: 'PostMatchController'
  }).otherwise({
    template: " ",
    controller: 'DefaultMatchController'
  });
});

rocketelo.controller('DraftBanController', function($scope, $location, socketIOService) {
    socketIOService.Initialize($scope.apiUrl, $scope.matchId);
    $scope.isDraftBan = true;
});
                     
rocketelo.controller('MatchController', function($scope, $location, socketIOService) {
    socketIOService.Initialize($scope.apiUrl, $scope.matchId);
    $scope.isGameStats = true;
});

rocketelo.controller('PostMatchController', function($scope, $location, socketIOService) {
    socketIOService.Initialize($scope.apiUrl, $scope.matchId);
    $scope.isPostmatch = true;
});

rocketelo.controller('DefaultMatchController', function($scope, $location) {
    $scope.init = function() {
        $location.path('/draft');
    }

    $scope.init();
});

rocketelo.controller('TeamDraftController', function ($scope, socketIOService) {
    $scope.init = function() {
        //TEMPORARY
        $scope.allChampions = [{champ: "/images/champions/Aatrox_0.jpg", player: "/images/players/c9-balls.jpg", playerName: "Balls"},
                               {champ: "/images/champions/Aatrox_0.jpg", player: "/images/players/c9-balls.jpg", playerName: "Balls"},
                               {champ: "/images/champions/Aatrox_0.jpg", player: "/images/players/c9-balls.jpg", playerName: "Balls"},
                               {champ: "/images/champions/Aatrox_0.jpg", player: "/images/players/c9-balls.jpg", playerName: "Balls"},
                               {champ: "/images/champions/Aatrox_0.jpg", player: "/images/players/c9-balls.jpg", playerName: "Balls"},];
        $scope.allBans = ["/images/champions/Ahri_Square_0.png", "/images/champions/Ahri_Square_0.png", "/images/champions/Ahri_Square_0.png"];
    }
    $scope.init();
});

rocketelo.directive('ngChampionDraft', function() {
    return {
        restrict: 'A',
        templateUrl: '/partials/match/templates/champion.html',
        scope: {
            champ: '@',
            player: '@',
            playername: '@'
        }
    }
});

rocketelo.directive('ngBanDraft', function() {
    return {
        restrict: 'A',
        templateUrl: '/partials/match/templates/ban.html',
        scope: {
            champ: '@'
        }
    }
});

rocketelo.controller('GameTimelineController', function($scope, socketIOService) {
    $scope.time = 0;
    
    $scope.init = function() {
        socketIOService.RegisterHandler($scope.ReceiveData);
    }
    
    $scope.ReceiveData = function(data) {
        console.log("Update time: " + $scope.time);
        $scope.time = LiveStats.GetTime(data);
        
        $scope.$apply();
    }
    
    $scope.init();
});

rocketelo.filter('secondsToDisplayTime', [function() {
    return function(s) {
        var minutes = Math.floor(s / 60);
        var seconds = s % 60;
        return minutes.toString() + ":" + seconds.toString();
    };
}]);

rocketelo.controller('TeamGameOverViewController', function($scope, socketIOService, compositeStatsService) {
    $scope.minimumMeterPercentage = 15.0;
    $scope.initWithTeamIndex = function(teamIndex) {
        $scope.teamIndex = teamIndex;
        
        // Default Values -- Probably should pull in the real values here.
        $scope.allPlayers = [
            {champ: "/images/champions/Aatrox_Square_0.png", player: "/images/players/c9-balls.jpg", playerName: "Balls", stats:{kills: 10, deaths: 5, assists: 8, creeps: 200}, 
                percentages:{kills: 100.0, deaths: 10.0, assists: 20.0, creeps: 30.0}},
            {champ: "/images/champions/Aatrox_Square_0.png", player: "/images/players/c9-balls.jpg", playerName: "Balls", stats:{kills: 10, deaths: 5, assists: 8, creeps: 200 },
                percentages:{kills: 100.0, deaths: 10.0, assists: 20.0, creeps: 30.0}},
            {champ: "/images/champions/Aatrox_Square_0.png", player: "/images/players/c9-balls.jpg", playerName: "Balls", stats:{kills: 10, deaths: 5, assists: 8, creeps: 200 },
                percentages:{kills: 100.0, deaths: 10.0, assists: 20.0, creeps: 30.0}},
            {champ: "/images/champions/Aatrox_Square_0.png", player: "/images/players/c9-balls.jpg", playerName: "Balls", stats:{kills: 10, deaths: 5, assists: 8, creeps: 200 },
                percentages:{kills: 100.0, deaths: 10.0, assists: 20.0, creeps: 30.0}},
            {champ: "/images/champions/Aatrox_Square_0.png", player: "/images/players/c9-balls.jpg", playerName: "Balls", stats:{kills: 10, deaths: 5, assists: 8, creeps: 200 },
                percentages:{kills: 100.0, deaths: 10.0, assists: 20.0, creeps: 30.0}}
        ];
        
        $scope.teamStats = {
            gold: 25200,
            towers: 3,
            dragons: 1,
            barons: 0,
            kills: 30,
            series: 0,
            team: "/images/teams/cloud9.png"
        };
        
        compositeStatsService.RegisterListener($scope.ReceiveNewCompositeData);
        socketIOService.RegisterHandler($scope.ReceiveData);
    }
    
    $scope.ReceiveNewCompositeData = function(composite) {
        $scope.overallStats = composite;
        for (var i = 0; i < 5; ++i) {
            $scope.allPlayers[i].percentages.kills = $scope.allPlayers[i].stats.kills / composite.kills * (100.0 - $scope.minimumMeterPercentage) + $scope.minimumMeterPercentage;
            $scope.allPlayers[i].percentages.deaths = $scope.allPlayers[i].stats.deaths / composite.deaths * (100.0 - $scope.minimumMeterPercentage) + $scope.minimumMeterPercentage;
            $scope.allPlayers[i].percentages.assists = $scope.allPlayers[i].stats.assists / composite.assists * (100.0 - $scope.minimumMeterPercentage) + $scope.minimumMeterPercentage;
            $scope.allPlayers[i].percentages.creeps = $scope.allPlayers[i].stats.creeps / composite.creeps * (100.0 - $scope.minimumMeterPercentage) + $scope.minimumMeterPercentage;
        }
        $scope.$apply();
    }
    
    $scope.ReceiveData = function(data) {
        // Strip data to the relevant bits and organize it such that the UI can use it.
        var myTeam = LiveStats.GetTeam(data, $scope.teamIndex);
        
        // Player Stats (Kills, Deaths, Assits, Creeps)
        for (var i = 0; i < 5; ++i) {
            var player = LiveStats.GetPlayerFromTeam(myTeam, i);
            
            var uiPlayerObject = {};
            if ($scope.allPlayers.length > i) {
                uiPlayerObject = $scope.allPlayers[i];
            } else {
                $scope.allPlayers.push(uiPlayerObject);   
            }
            
            uiPlayerObject.champ = LiveStatsUtility.GetChampionProfilePicture(LiveStats.GetPlayerChampionName(player));
            uiPlayerObject.playerName = LiveStats.GetPlayerName(player);
            uiPlayerObject.player = LiveStatsUtility.GetPlayerProfilePicture(uiPlayerObject.playerName);
            uiPlayerObject.stats = LiveStats.GetPlayerOverviewStats(player);
            
            // Cleanup Stats (make sure they're not negative...)
            uiPlayerObject.stats.kills = Math.max(uiPlayerObject.stats.kills, 0);
            uiPlayerObject.stats.deaths = Math.max(uiPlayerObject.stats.deaths, 0);
            uiPlayerObject.stats.assists = Math.max(uiPlayerObject.stats.assists, 0);
            uiPlayerObject.stats.creeps = Math.max(uiPlayerObject.stats.creeps, 0);
            
            compositeStatsService.SetStatsForPlayerOnTeam(i, $scope.teamIndex, uiPlayerObject.stats);
        }
        
        // Compute composite stats (aka combine stats from both teams to get meaningful data).
        // This is necessary because we want to set the "max" value of the bars for kills, deaths, etc. to be based off the 
        // maximum value existing in the game at the current moment in time to allow clients to make meaningful visual comparisons.
        compositeStatsService.ComputeCompositeStats();
        
        // Team Stats (Gold, Towers, Dragons, Barons, Kills)
        $scope.teamStats = {
            gold: Math.max(LiveStats.GetTeamGold(myTeam), 0),
            towers: Math.max(LiveStats.GetTeamTowers(myTeam), 0),
            dragons: Math.max(LiveStats.GetTeamCurrentDragons(myTeam), 0),
            barons: Math.max(LiveStats.GetTeamBarons(myTeam), 0),
            kills: Math.max(LiveStats.GetTeamKills(myTeam), 0),
            series: Math.max(LiveStats.GetTeamSeriesWins(myTeam), 0),
            team: LiveStatsUtility.GetTeamImage(myTeam)
        };
        
        $scope.$apply();
    }
});

rocketelo.directive('ngPlayerStats', function() {
    return {
        restrict: 'A',
        templateUrl: '/partials/match/templates/playerStats.html',
        scope: {
            player: '='
        }
    }
});

rocketelo.directive('ngStatBar', function() {
    return {
        restrict: 'A',
        templateUrl: '/partials/match/templates/statBar.html',
        scope: {
            barclass: '@',
            value: '@',
            min: '@',
            max: '@',
            percentage: '@',
            label: '@'
        }
    }
});