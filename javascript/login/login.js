rocketelo.controller('LoginController', function($scope, $http, $window) {
  $scope.loginSucceed = true;
  $scope.usernameError = false;

  // LOGIN FUNCTIONALITY
  $scope.login = function(user) {
    $http.post('/login', {"username": user.username, "password": user.password}).
      success(function(data, status, headers, config) {
        // Successful login so just do a redirect to wherever the server tells us to go to 
        $window.location.href = data.RedirectUrl;
      }).
      
      error(function(data, status, headers, config) {
        $scope.loginSucceed = true;
        retObj = angular.fromJson(data);
        // Notify user that login failed. Only way login fails is if either something we awry on the server or the username/password don't match.
        if (retObj.ErrorCode == 6) {
          // Unknown error on the server :(
        } else {
          $scope.loginSucceed = false; 
        }
      });
  };
  
  
  // REGISTER FUNCTIONALITY
  $scope.register = function(user) {
    $http.post('/register', {"username": user.username, "password": user.password, "email": user.email}).
      success(function(data, status, headers, config) {
        // Successful register so just do a redirect to wherever the server tells us to go to 
        $window.location.href = data.RedirectUrl;
      }).
      
      error(function(data, status, headers, config) {
        $scope.usernameError = false;
        console.log("on error register " + data);
        retObj = angular.fromJson(data);
        if (retObj.ErrorCode == 6) {
          // Unknown error on the server :(
        } else if (retObj.ErrorCode == 1){
          $scope.usernameError = true; 
        }
        console.log(retObj.ErrorCode);
      });
  };

});