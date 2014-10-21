rocketelo.controller('EditGameController', function($scope, $http, $window) {
  // CREATE/UPDATE GAME FUNCTIONALITY
  $scope.save = function(game, isnew, oldGameShorthand) {
    $http.post('/admin/game', {isNew: isnew, gameName: game.gameName, gameShorthand: game.gameShorthand, oldShorthand: oldGameShorthand}).
      success(function(data, status, headers, config) {
        // Redirect back to game list
        $window.location.href = "/admin/game";
      }).
      
      error(function(data, status, headers, config) {
        // Display a message telling us what went wrong
        retObj = angular.fromJson(data);
        // Notify user that login failed. Only way login fails is if either something we awry on the server or the username/password don't match.
        if (retObj.ErrorCode == 1 && isnew) {
          $scope.setFlashMessage("An existing game already exists. Please choose a new shortname.");
        } else if (retObj.ErrorCode == 1 && !isnew) {
          $scope.setFlashMessage("The game you are trying to modify does not exist anymore or the new game shortname is invalid. Please refresh, choose a new shortname and try again.");
        } else {
          $scope.setFlashMessage("Unknown error. Please try again.");
        }
      });
  };
});