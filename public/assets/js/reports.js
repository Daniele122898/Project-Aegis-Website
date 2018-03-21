$(document).ready(function () {
    SetPreview();
    SetSubmit();
});

function ShowError(msg) {
    var al = $(".alertPos");
    al.empty();
    al.append('<div class="alert alert-danger alert-dismissible">' +
        '<a href="#" class="close" data-dismiss="alert" aria-label="close">&times;</a>' +
        '<strong>Error!</strong>   '+msg +
        '</div>');
}

function ShowSuccess(msg){
    var al = $(".alertPos");
    al.empty();
    al.append('<div class="alert alert-success alert-dismissible">' +
        '<a href="#" class="close" data-dismiss="alert" aria-label="close">&times;</a>' +
        '<strong>Success!</strong>   '+msg +
        '</div>');
}

function SetSubmit() {
    var btn = $("#sub");
    var clicked = false;
    btn.click(function () {
        if(clicked){
            return;
        }
        clicked = true;
        var reason = $("#reasonText").val();
        var guildId = $("#inputBox").val();
        if (guildId === null || guildId ===0 ){
            ShowError("Enter a Guild Id!");
            clicked = false;
            return;
        }
        if (reason ===null || reason ===""){
            ShowError("Enter a valid Reason!");
            clicked = false;
            return;
        }
        //post request
        $.ajax({
            url: "http://project-aegis.pw/api/guildReport/"+guildId,
            type: "POST",
            data: JSON.stringify({reason: reason}),
            statusCode: {
                200: function (response) {
                    ShowSuccess("Submitted Report.");
                    $("#reasonText").val("");
                    $("#inputBox").val(null);
                    clicked = false;
                },
                400: function (response) {
                    ShowError(response.responseJSON.error);
                    clicked = false;
                },
                500: function (response) {
                    ShowError(response.responseJSON.error);
                    clicked = false;
                }
            }, success: function () {
            }
        });

    });
}
function SetPreview() {
    var prev = $("#previewWindow");
    var btn = $("#prev");
    btn.click(function () {
        var reason = $("#reasonText");
        var markdown = "";
        console.log("SEND: "+reason.val());
        //Send and receive markdown
        $.post("http://project-aegis.pw/api/generateMarkdown", reason.val(), function (data) {
            markdown = data;
            console.log("GOT DATA: "+markdown);
            prev.append('<div class="container">' +
                '<div class="contentBox">' +
                '<div id="close">' +
                '<i class="icon ion-close-round"></i>' +
                '</div>' + markdown+
                '</div>' +
                '</div>');
            prev.show();
            $("#close").click(function () {
                prev.empty();
                prev.hide();
            })
        });
    })
}