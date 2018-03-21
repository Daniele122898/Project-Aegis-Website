$(document).ready(function(){
    //smooth scrolling on anchor links
    $(document).on('click', 'a[href^="#"]', function (event) {
        event.preventDefault();

        $('html, body').animate({
            scrollTop: $($.attr(this, 'href')).offset().top
        }, 500);
    });

    //change arrows on FAQ open
    $(".btn").on("click" , function () {
        var icon = $(this).find('i');
        //RESET ALL OTHERS
        if(icon.hasClass("ion-ios-arrow-down")){
            //clicked on an down arrow
            icon.removeClass("ion-ios-arrow-down");
            icon.addClass("ion-ios-arrow-right");
            return;
        }
        var btns = $(".btn").find('i');
        btns.removeClass("ion-ios-arrow-down");
        btns.removeClass("ion-ios-arrow-right");
        btns.addClass("ion-ios-arrow-right");
        icon.removeClass("ion-ios-arrow-right");
        icon.addClass("ion-ios-arrow-down");
        //$(this).find('i').toggleClass("ion-ios-arrow-right");
        //$(this).find('i').toggleClass("ion-ios-arrow-down");
    });

});