$(document).ready(function () {
    createList();
});

function createList() {
    var box = $(".inputBox");
    $.ajax({
        url: "http://project-aegis.pw/api/guildlist",
        type: "GET",
        statusCode: {
            200: function (response) {
                var data = JSON.parse(response);
                for(var i=0; i<data.length; i++){
                    var d = data[i];

                    var secLev = "";
                    var secClass = "";
                    switch(d.securityLevel) {
                        case 0:
                            secLev= "Good";
                            secClass= "good";
                            break;
                        case 1:
                            secLev= "Suspicious";
                            secClass= "suspicious";
                            break;
                        case 2:
                            secLev= "Moderate";
                            secClass= "moderate";
                            break;
                        case 3:
                            secLev= "High";
                            secClass="high";
                            break;
                        case 4:
                            secLev= "Extreme";
                            secClass="extreme";
                            break;
                        default:
                    }


                    box.append(' <div class="guild">' +
                        '<h5 class="'+(d.checked ? "green": "red")+'">'+d.id+'</h5>' +
                        '<div class="row">' +
                        '<div class="col-4">' +
                        '<p style="text-align: left">Security: <span class="'+secClass+'">'+secLev+'</span></p>' +
                        '</div>' +
                        '<div class="col-4">' +
                        '<p style="text-align: center">Reports: <span>'+d.reportCount+'</span></p>' +
                        '</div>' +
                        '<div class="col-4">' +
                        '<a href="/guild/'+d.id+'" style="float: right"><i class="icon ion-eye"></i>View</a>' +
                        '</div>' +
                        '</div>' +
                        '</div>')
                }
            },
            400: function (response) {
                ShowError(response.responseJSON.error);
            }
        }, success: function () {
        }
    });
}

function ShowError(msg) {
    var al = $(".alertPos");
    al.empty();
    al.append('<div class="alert alert-danger alert-dismissible">' +
        '<a href="#" class="close" data-dismiss="alert" aria-label="close">&times;</a>' +
        '<strong>Error!</strong>   '+msg +
        '</div>');
}