// main.js - basic image review behaviour
// zepto: for when you need some API action, but can not be bothered
// to deal with angular.

var currentImage;
$(function() {
    loadImage();
    function loadImage() {
        $.getJSON("/api/images/random", function(data) {
            console.log(data);
            $("#image").attr("src", data.data.url);
            $("#input_form").attr("action", "/api/images/" + data.data.id + "/numbers/new");
            currentImage = data.data;
        });
    }

    $("#go").click(function() {
        $.post("/api/images/" + currentImage.id + "/numbers/new?number=" + $("#number").val());
        loadImage();
    });
});
