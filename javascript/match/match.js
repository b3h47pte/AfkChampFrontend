// Socket IO Initialization
var socket = io();

rocketelo.config(function($routeProvider) {
  $routeProvider.when('/draft', {
    templateUrl : '/partials/match/draft.html',
    controller: 'DraftBanController'
  }).when('/match', {
    templateUrl : '/partials/match/match.html',
    controller: 'MatchController'
  }).when('/post', {
    templateUrl : '/partials/match/post.html',
    controller: 'PostMatchController'
  }).otherwise({
    template: " ",
    controller: 'DefaultMatchController'
  });
});

rocketelo.controller('DraftBanController', function($scope, $location) {
    
});
                     
rocketelo.controller('MatchController', function($scope, $location) {
    
});

rocketelo.controller('PostMatchController', function($scope, $location) {
    
});

rocketelo.controller('DefaultMatchController', function($scope, $location) {
  $scope.init = function() {
    $location.path('/draft');
  }

  $scope.init();
});
