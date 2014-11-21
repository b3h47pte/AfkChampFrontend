rocketelo.controller('MatchFilterController', function($scope, $location, $document, matchDisplayService) {
  $scope.currentDate = moment().format("MMMM D")
  $scope.showFilter = false;

  $scope.toggleFilter = function() {
    var filter = $document.find("#match-filter");
    var link = $document.find("#filter-link");
    link.toggleClass("down");
    $scope.showFilter = !$scope.showFilter;
    if ($scope.showFilter) {
      filter.show();
    } else {
      filter.hide();
    }
  }
});

rocketelo.directive('jqdatepicker', function() {
  return {
    restrict: 'A',
    require: 'ngModel',
    link: function(scope, element, attrs, ngModelCtrl) {
      element.datepicker({
      }); 
    }
  };
});

rocketelo.directive('jqtimepicker', function() {
  return {
    restrict: 'A',
    require: 'ngModel',
    link: function(scope, element, attrs, ngModelCtrl) {
      element.timepicker({
      }); 
    }
  };
});
