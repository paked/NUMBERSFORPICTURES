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

app.controller('SubmitCtrl', function($scope, $http, $location) {
    $scope.url = "";

    $scope.submit = function() {
        if (!$scope.url) {
            return;
        }

        $http.post("/api/images/new?url=" + $scope.url)
            .success(function(d) {
                console.log(d);
                $location.path("/");
            });
    };
});

app.controller('ImageCtrl', function($scope, $http, $rootScope) {
    $scope.image = {};
    $rootScope.numbers = [];

    $scope.loadImage = function()  {
        $http.get("/api/images/random")
            .success(function(data) {
                console.log(data);
                $scope.image = data.data;
                console.log("loaded: ", $scope.image.id);
                $scope.$broadcast("SendID", $scope.image.id);
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

app.controller('StatisticsCtrl', function($scope, $rootScope, $http) {
    $scope.$on('SendID', function(event, id) {
        console.log("heep", id);

        $http.get("/api/images/" + id + "/numbers")
            .success(function(d) {
                console.log(d, "<-- d");
                if (d.data === null) {
                    return;
                }
                
                $scope.compute(d.data);
            });
    });

    $scope.compute = function(numbers) {
        console.log("in compute");
        function toNumbers(ns) {
            var i = 0;
            var res = [];
            while (i < ns.length) {
                res.push(ns[i].number);
                i++;
            }

            return res;
        }

        var res = toNumbers(numbers);

        $scope.mostCommon = 0;
        $scope.highestNumber = Math.max.apply(Math, res);
        $scope.lowestNumber = Math.min.apply(Math, res);
        $scope.numbers = numbers;
    };
});
