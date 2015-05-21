// main.js - cause we all love some good JS meat
// by Harrison Shoebridge

var app = angular.module('wnip', []);

app.controller('MainCtrl', function($scope) {
    $scope.hello = "Hello, world!";
});

app.controller('ImageCtrl', function($scope, $http) {
    $scope.image = {};
    $scope.loadImage = function()  {
        $http.get("/api/images/random").success(function(data) {
            console.log(data);
            $scope.image = data.data;
            console.log("loaded");
        });
    };

    $scope.go = function() {
        $http.post("/api/images/" + $scope.image.id + "/numbers/new?number=" + $scope.number);
        console.log("posted");
    };

    $scope.loadImage();
});
