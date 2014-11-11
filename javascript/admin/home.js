var rocketelo = angular.module('rocketelo', ['ui.utils', 'ngRoute']).
  config(function($locationProvider) {
    $locationProvider.html5Mode(true);
  });

angular.element("a").prop("target", "_self");

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
