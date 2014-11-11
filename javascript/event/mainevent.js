rocketelo.controller('MainEventPageController', function($scope, $location) {
  var eventShorthand = $location.path().split("/")[2] || "none";
  var matchesResponse = null;
  $scope.init = function() {
      $scope.refreshMatches();
  }

  $scope.refreshMatches = function() {
    console.log(eventShorthand);
  }

  $scope.init();
});
