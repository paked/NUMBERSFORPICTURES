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
    });
});
