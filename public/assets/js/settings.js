var lkajsdlkasj ="";

$(document).ready(function () {
    Token();
    Sync();
    GetInitialInfo();
});

function GetInitialInfo() {

    $.get("http://project-aegis.pw/api/getUserInfo", function (response) {
        var data = JSON.parse(response);
        //change info
        $("#userInfo").html('<h3>'+data.username+'<small>#'+data.discriminator+'</small></h3>' +
            '<p>ID: <span id="userId">'+data.id+'</span></p>');
        //change avatar
        $("#userImage").attr("src", data.avatar);
    });

}

function Sync(){
    var sync = $("#updateData");
    var clicked = false;
    sync.click(function () {
        if(clicked){
            ShowError("Can't Sync again!");
            return;
        }
        else{
            clicked = true;
        }
        //do post request to update info
        $.ajax({
            url: "http://project-aegis.pw/api/syncProfile",
            type: "GET",
            statusCode: {
                200: function (response) {
                    ShowSuccess("Synced Data.");
                },
                400: function (response) {
                    ShowError(response.responseJSON.error);
                },
                500: function (response) {
                    ShowError(response.responseJSON.error);
                }
            }, success: function () {
            }
        });
    });
}

function Token(){
    var tknH = $("#tokenHider");
    var showToken = $("#showToken");
    showToken.click(function () {
        //tknH.toggleClass("tokenHider")
        if(tknH.hasClass("tokenHider")){
            tknH.removeClass("tokenHider");
            showToken.text("Hide Token");
            //Request token
            if(lkajsdlkasj===""){
                //Get token
                $.get("http://project-aegis.pw/api/getToken", function (response) {
                    var data = JSON.parse(response);
                    //change
                    lkajsdlkasj = data.token;
                    $(".tokenContainer").text(data.token);
                });
            } else{
                $(".tokenContainer").text(lkajsdlkasj);
            }
        } else {
            tknH.addClass("tokenHider");
            showToken.text("Show Token");
        }
    });

    //Token Generation
    var tokenGen = $("#genToken");
    var clicked = false;
    tokenGen.click(function () {
        if(clicked){
            ShowError("Can't Generate again!");
            return;
        }
        else{
            clicked = true;
        }
        $.ajax({
            url: "http://project-aegis.pw/api/genNewToken",
            type: "GET",
            statusCode: {
                200: function (response) {
                    ShowSuccess("Generated new Token.");
                    var data = JSON.parse(response);
                    $(".tokenContainer").text(data.newToken);
                    lkajsdlkasj = data.newToken;
                },
                400: function (response) {
                    ShowError(response.responseJSON.error);
                },
                500: function (response) {
                    ShowError(response.responseJSON.error);
                }
            }, success: function () {
            }
        });
    })
}

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