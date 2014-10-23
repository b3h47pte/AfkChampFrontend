rocketelo.controller('EditUserController', function($scope, $http, $window) {
  // CREATE/UPDATE USER FUNCTIONALITY
  $scope.save = function(user, isNew, userId) {
    if (isNew) {
      return;
    }
    submitUser = user
    submitUser.UserId = userId;
    submitUser.isadmin = (user.isadmin == "true");
    $http.post('/admin/user', {isNew: isNew, user: submitUser}).
      success(function(data, status, headers, config) {
        // Redirect back to game list
        $window.location.href = "/admin/user";
      }).
      
      error(function(data, status, headers, config) {
        // Display a message telling us what went wrong
        retObj = angular.fromJson(data);
        $scope.setFlashMessage("Unknown error. Please try again.");
      });
  };
});