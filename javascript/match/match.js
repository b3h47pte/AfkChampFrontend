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
});
                     
rocketelo.controller('MatchController', function($scope, $location, socketIOService) {
    socketIOService.Initialize($scope.apiUrl, $scope.matchId);
});

rocketelo.controller('PostMatchController', function($scope, $location, socketIOService) {
    socketIOService.Initialize($scope.apiUrl, $scope.matchId);
});

rocketelo.controller('DefaultMatchController', function($scope, $location) {
  $scope.init = function() {
    $location.path('/draft');
  }

  $scope.init();
});
