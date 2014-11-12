function getMatchShorthand(location) {
  var eventShorthand = $location.path().split("/")[2] || "none";
  return eventShorthand;
}

rocketelo.config(function($routeProvider) {
  $routeProvider.when('/current', {
    templateUrl : 'html/angular/currentMatches.html',
    controller: 'CurrentMatchesController'
  }).when('/upcoming', {
    templateUrl : 'html/angular/currentMatches.html',
    controller: 'CurrentMatchesController'
  }).when('/recent', {
    templateUrl : 'html/angular/currentMatches.html',
    controller: 'CurrentMatchesController'
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
