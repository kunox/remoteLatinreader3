
(function($) {
    $.wordbox=function(el) {
        var xpos, ypos;     //  左上座標
        var opened;
        var plugin = this;
        var $wordbox = el;
    

         plugin.init = function(){
            plugin.el = el;
            xpos = ($(window).width() - el.outerWidth(true))/2;
            ypos =  ($(window).height() - el.outerHeight(true)) / 2;
            // $wordbox.offset({top:ypos, left:xpos});
            $wordbox.children(".wb-header").mousedown(function(ex){
                var my=ex.pageY - $wordbox.offset().top;
                var mx = ex.pageX - $wordbox.offset().left;

                $wordbox.on("mousemove", function(e){
                    $wordbox.offset({top:e.pageY-my, left:e.pageX-mx});
                    return false;
                });
                $(document).one("mouseup",function(e){
                    ypos = $wordbox.offset().top - $(document).scrollTop();
                    xpos =  $wordbox.offset().left;
                    $wordbox.off("mousemove");
                    $wordbox.css("position", "fixed");
                });
            });
        }

        plugin.openbox = function(){
            // plugin.el.offset({top:ypos, left:xpos});
            $wordbox.css({ top: ypos, left: xpos });
            $wordbox.show();
        }

        plugin.closebox = function(){
            $wordbox.hide();
        }

        // init();
    }
})(jQuery);