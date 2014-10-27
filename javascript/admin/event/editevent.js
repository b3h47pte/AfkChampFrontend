rocketelo.controller('EditEventController', function($scope) {
  // CREATE/UPDATE EVENT FUNCTIONALITY
  $scope.save = function(event) {
    $http.post('/admin/event', {}).
      success(function(data, status, headers, config) {
        // Redirect back to game list
        $window.location.href = "/admin/event";
      }).
      
      error(function(data, status, headers, config) {
        // Display a message telling us what went wrong
        retObj = angular.fromJson(data);
        if (retObj.ErrorCode == 1) {
          $scope.setFlashMessage("Some error happened.");
        } else {
          $scope.setFlashMessage("Unknown error. Please try again.");
        }
      });
  };
});
