rocketelo.factory('socketIOService', function($rootScope) {
    var socketService = {}
    socketService.isInit = false;
    socketService.Initialize = function(url, matchId) {
        if (socketService.isInit) {
            return;
        }
        socketService.isInit = true;
        socketService.socket = io.connect(url);
        socketService.socket.emit('identify', {match: matchId.toString()});
    }
    
    return socketService;
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

rocketelo.controller('TeamDraftController', function ($scope) {
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

rocketelo.controller('TeamGameOverViewController', function($scope) {
    $scope.init = function() {
        //TEMPORARY
        $scope.allPlayers = [
            {champ: "/images/champions/Aatrox_Square_0.png", player: "/images/players/c9-balls.jpg", playerName: "Balls", stats:{kills: 10, deaths: 5, assists: 8, creeps: 200 }},
            {champ: "/images/champions/Aatrox_Square_0.png", player: "/images/players/c9-balls.jpg", playerName: "Balls", stats:{kills: 10, deaths: 5, assists: 8, creeps: 200 }},
            {champ: "/images/champions/Aatrox_Square_0.png", player: "/images/players/c9-balls.jpg", playerName: "Balls", stats:{kills: 10, deaths: 5, assists: 8, creeps: 200 }},
            {champ: "/images/champions/Aatrox_Square_0.png", player: "/images/players/c9-balls.jpg", playerName: "Balls", stats:{kills: 10, deaths: 5, assists: 8, creeps: 200 }},
            {champ: "/images/champions/Aatrox_Square_0.png", player: "/images/players/c9-balls.jpg", playerName: "Balls", stats:{kills: 10, deaths: 5, assists: 8, creeps: 200 }}
        ];
    }
    
    $scope.init();
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