// main.js - cause we all love some good JS meat
// by Harrison Shoebridge

var app = angular.module('wnip', ['ngRoute']);

app.config(function($routeProvider, $locationProvider) {
    $routeProvider
    .when('/submit', {
        templateUrl: 'templates/submit.html',
        controller: 'SubmitCtrl'
    }).otherwise({
        templateUrl: 'templates/view_image.html',
        controller: 'ImageCtrl',
    });
});

app.controller('MainCtrl', function($scope) {
    $scope.hello = "Hello, world!";
});

app.controller('ImageCtrl', function($scope, $http) {
    $scope.image = {};
    $scope.numbers = [];

    $scope.loadImage = function()  {
        $http.get("/api/images/random").success(function(data) {
            console.log(data);
            $scope.image = data.data;
            console.log("loaded: ", $scope.image.id);

            $http.get("/api/images/" + $scope.image.id + "/numbers")
                .success(function(d) {
                    $scope.numbers = d.data;
                });
        });
    };

    $scope.go = function() {
        $http.post("/api/images/" + $scope.image.id + "/numbers/new?number=" + $scope.number)
            .success(function(d) {
                $scope.numbers.push(d.data);
            });
    };

    $scope.loadImage();
});
