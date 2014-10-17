var rocketelo = angular.module('rocketelo', ['ui.utils']);

rocketelo.controller('MainBodyController', function($scope) {
  var flashMessage = "";

  $scope.getFlashMessage = function() {
    return flashMessage;
  }
  
  $scope.setFlashMessage = function(newMessage) {
    flashMessage = newMessage;
  }

  $scope.clearFlashMessage = function() {
    flashMessage = "";
  }
});