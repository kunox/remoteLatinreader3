 // StoreしておくChapter用のscript
 // jqueryは使える前提（CDNでロード）
window.onload = function () {
    var touchdevice = ('ontouchstart' in window || navigator.msPointerEnabled) ? true : false;
    if (touchdevice == true) {
        var moevent = new Event('mouseover');

        var obj = document.getElementsByClassName('note1');
        for (var i = 0; i < obj.length; i++) {
            obj[i].addEventListener('touchstart', function (event) {
                event.target.dispatchEvent(moevent);
            }, false);
        }
    }

    if($("#saveoption").length){
        $("#saveoption").hide();
        $("#downloadsave").hide();
    }
    $("#pathname").hide();
    insertButtons2();
    $("#outer").css("width","96%");

}

function insertButtons2(){
    if(!$(".lbtn1")){
        $(".spacer").each(function(index, elem) {
            var target = $(elem);
            var tx1 = target.html();
            var color = "style='color:#5f9ea0;'";
            if (target.next().next().attr("class") != "spnote")
                color = "style='color:black;'";
            target.html(tx1 + "<button class='lbtn2'>Translation</button><button class='lbtn1'>Notes</button>");
        });
    }
    $(".lbtn1").on("click", function() {
        var target = $(this).parent().next().next();
        if (target.attr("class") != "spnote") {
            target = target.next();
            if (target.attr("class") != "spnote") // 編集中の場合
                return;
        }
        btoggle(target, $(this));
    });
    $(".lbtn2").on("click", function() {
        var target = $(this).parent().next().next().next();
        if (target.attr("class") != "translation") {
            target = target.next();
            if (target.attr("class") != "translation")
                return;
        }
        btoggle(target, $(this));
    });
}

function btoggle(obj, btn) {
    obj.toggle(function() {
        if ($(this).is(":visible"))
            btn.css("color", "black");
        else
            btn.css("color", "#5f9ea0");
    });
}
