var rocketelo = angular.module('rocketelo', ['ui.utils']);

rocketelo.controller('LoginController', function($scope, $http) {
  // LOGIN FUNCTIONALITY
  $scope.login = function(user) {
    $http.post('/login', {"username": user.username, "password": user.password});
  };
  
  
  // REGISTER FUNCTIONALITY
  $scope.register = function(user) {
    $http.post('/register', {"username": user.username, "password": user.password, "email": user.email});
  }

});