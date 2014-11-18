function getMatchShorthand(location) {
  var eventShorthand = $location.path().split("/")[2] || "none";
  return eventShorthand;
}

rocketelo.config(function($routeProvider) {
  $routeProvider.when('/current', {
    templateUrl : '/partials/event/currentMatches.html',
    controller: 'CurrentMatchesController'
  }).when('/upcoming', {
    templateUrl : '/partials/event/upcomingMatches.html',
    controller: 'UpcomingMatchesController'
  }).when('/recent', {
    templateUrl : '/partials/event/recentMatches.html',
    controller: 'PreviousMatchesController'
  });
});

rocketelo.controller('CurrentMatchesController', function($scope, $location) {
  $scope.init = function() {
      $scope.refreshMatches();
  }

  $scope.refreshMatches = function() {
  }

  $scope.init();
});

rocketelo.controller('PreviousMatchesController', function($scope, $location) {
  $scope.init = function() {
      $scope.refreshMatches();
  }

  $scope.refreshMatches = function() {
  }

  $scope.init();
});

rocketelo.controller('UpcomingMatchesController', function($scope, $location) {
  $scope.init = function() {
      $scope.refreshMatches();
  }

  $scope.refreshMatches = function() {
  }

  $scope.init();
});


