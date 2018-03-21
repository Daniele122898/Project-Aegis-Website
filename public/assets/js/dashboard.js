$(document).ready(function () {
    getSearch();
});

function getSearch(){
    var search = $("#inputBox");
    search.keypress(function (e) {
        var val = search.val();
        if (e.which == 13) {
            if (val === null || val <10000){
                return false;
            }
            window.location.href ="http://project-aegis.pw/guild/"+val;
            return false;
        }
    });

    $(".searchButton").click(function () {
        var val = search.val();
        if (val === null || val <10000){
            return;
        }
        window.location.href="http://project-aegis.pw/guild/"+val;
    });
}