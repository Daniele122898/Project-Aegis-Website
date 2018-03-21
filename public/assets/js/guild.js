$(document).ready(function () {
    getGuild();
});

function getGuild(){
    var id = getGuildIdFromPath();

    //Get Guild info
    $.ajax({
        url: "http://project-aegis.pw/api/guild/"+id,
        type: "GET",
        statusCode: {
            200: function (response) {
                var data = JSON.parse(response);
                //set website
                $("#guildId").text(data.id);
                $("#reportAmount").text(data.reportCount);
                //Sec level
                var seclvl = $("#secLevel");

                switch(data.securityLevel) {
                    case 0:
                        seclvl.text("Good");
                        seclvl.addClass("green");
                        break;
                    case 1:
                        seclvl.text("Suspicious");
                        seclvl.addClass("suspicious");
                        break;
                    case 2:
                        seclvl.text("Moderate");
                        seclvl.addClass("moderate");
                        break;
                    case 3:
                        seclvl.text("High");
                        seclvl.addClass("high");
                        break;
                    case 4:
                        seclvl.text("Extreme");
                        seclvl.addClass("extreme");
                        break;
                    default:
                }

                var chekd = $("#checked");
                if( data.checked){
                    chekd.text("Checked: Yes");
                    chekd.addClass("greenText");
                } else {
                    chekd.text("Checked: No");
                    chekd.addClass("redText");
                }

                //comments
                postComments(data);
            },
            400: function (response) {
                ShowError(response.responseJSON.error);
            }
        }, success: function () {
        }
    });
}

function postComments(data) {
    var row = $("#commentRow");
    var html = "";
    for (var i=0; i<data.reports.length; i++){
        var report = data.reports[i];
        var classToAdd ="";
        if (report.closed){
            classToAdd = "closedReport";
        }else{
            classToAdd = "openReport";
        }
        var comments ="";
        var split = report.text.split("-- EDIT --");
        for (var j=0; j<split.length; j++){
            var s = split[j];
            comments += '<div class="commentBlock '+classToAdd+'">' +s+'</div>';
        }
        html+= '<div class="comment">' +
            '<img src="'+report.user.avatar+'" alt="Avatar"><br>' +
            '<div class="commentContent">' +
            '<h5>'+report.user.username+'<small>#'+report.user.discrim+' ~ <span>'+getDate(report.date)+'</span></small></h5>' +
            //'<div class="commentBlock '+classToAdd+'">' +report.text+
            comments+
            '</div>' +
            '</div>'
    }

    row.html(html);
}

function getDate(unix) {
    // Create a new JavaScript Date object based on the timestamp
// multiplied by 1000 so that the argument is in milliseconds, not seconds.
    var date = new Date(unix*1000);
    return date.toLocaleDateString();
/* Hours part from the timestamp
    var hours = date.getHours();
// Minutes part from the timestamp
    var minutes = "0" + date.getMinutes();
// Seconds part from the timestamp
    var seconds = "0" + date.getSeconds();

// Will display time in 10:30:23 format
    return hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);*/
}

function getGuildIdFromPath(){
    var path = window.location.pathname;
    // /guild/1231241
    return path.substring(path.indexOf("/", 2) + 1, path.length);
}