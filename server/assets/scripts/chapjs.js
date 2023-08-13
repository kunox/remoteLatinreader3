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

    // space barにボタン追加
    insertButtons();

    //spnote translationをeditmodeに
    $("div.innerdiv").parent().on("dblclick", editingmode);

    window.onbeforeunload = calledwhenunload

    // word クリックで辞書を動かす部分
    $('[class ^="note"]').on("click", xeditNote3);
    $("#wbtn").on("click", postjson21);              // #wbtnは画面右にある「送信」ボタン

    if(!$("#saveoption").length){
        $("#header").append(spacebutton);   //buttonがすでにある可能性もある
    }
    $("#bodyportion").after(wboxhtml);       //bodyportionのあとに追加しておく
    $("#saveoption").on("click", uploadsave);
    $("#downloadsave").on("click", downloadsave)        // Local版ではdisplay:None

    // wordbox関連の初期化
    initializeWordbox();

    $("#msgbox").dialog({
        modal:true,
        autoOpen:false,
        show: "slide",
        title:"message",
        width:300,
        height:250,
        draggable:false,
        buttons: {
            "OK":function(){
                $(this).dialog("close");
            }
        }
    });

    $("div.latin").removeAttr('style');
    $("div.innerdiv").removeAttr('style');
// 活用表の表示
    // $("#disp").on("contextmenu", ".ditems", function (e) {
    //     e.preventDefault();
    //     var msg = "tcreation#" + $(this).children(".exdc").text();
    //     window.chrome.webview.postMessage(msg);    // send back exdc comparables

    //     //alert($(this).children(".exdc").text());
    //     return false;           // defaultの右クリックメニュが出ないように
    // });
}

var spacebutton = `
<div class="boxright">
<input id="saveoption" class = "rbtn1" type="button" value="保存する"
    style="font-size: 12px; font-family: meiryo; " />
<input id="downloadsave" class="rbtn1" type = "button" value="ダウンロード"
    style="font-size: 12px; font-family: meiryo;" />
</div>`
var wboxhtml = `
<div class="wordbox" id="wordbox">
<div class="wb-header">
  <button type="button" id="nbold" class="sbutton">b</button>
  <button type="button" id="nunder" class="sbutton">u</button>
  <button type="button" id="nitalic" class="sbutton">i</button>
  <button type="button" id="nstrike" class="sbutton">s</button>

  <button type="button" id ="closebtn" class="save-close"> x </button>
  <button type="button" id="savebtn" class="save-close">save</button>
</div>
<div>
  <input type="text" class ="ta" id="word"/>
</div>
<div class="wb-content" id=meanings>
  <p>ダイアログです</p>
</div>
</div>

<div id="msgbox" style="display:none;">
<p id="msgtxt">Dialog</p>
</div>
`


var tempdiv = `
    <div class="ditems">
        <p class="exdc">%s0</p>
        <span class="dw1">%s1</span>:<span class="dw2">%s2</span>
        <p class="dw3">%s3</p>
        <p class="dw4">%s4</p>
    </div>
`
// msgbox 
function msgbox(msg){
    $("#msgtxt").html(msg);
    $("#msgbox").dialog("open");
}


//*************** 単語クリック 関連 ********************************** */
// 単語のクリック
var $wordnote   // クリックされた語を覚えておく
var popbox;
function initializeWordbox(){
    popbox = new $.wordbox($("#wordbox"));
    popbox.init();
    $("#closebtn").on("click",wcancel);
    $("#savebtn").on("click", wordsave);

    $("#nbold").on("click", function() { if(checkSelection()) {document.execCommand("bold", false, true);} });
    $("#nunder").on("click", function() { if(checkSelection()) {document.execCommand("underline", false);} });
    $("#nitalic").on("click", function() { if(checkSelection()) {document.execCommand("italic", false);} });
    $("#nstrike").on("click", function() { if(checkSelection()) {document.execCommand("strikeThrough", false);} });
}

function xeditNote3() {
    $wordnote = $(this);
    var txt = $(this).children().html();
    var txt2 = $(this).html();
    var start = txt2.indexOf("<span");
    var w = txt2.substring(0, start);
    
    $("#keyin").val(w)

    postjson2();

    popbox.openbox();
    $("#word").val(w);
    $("#meanings").empty().append("<p>" + txt + "</p>");
    $("#meanings").attr("contentEditable", true);
    $("#meanings").focus();

    // $("#wbox").dialog("option","title",w);
    // $("#wboxtext").html(txt);
    // $("#wbox").dialog("open");

    // 直接 ditemsにリスナーをつけるとタイミングにより発火しない。
    // すでにある要素に対してリスナーを定義しておく
    $("#disp").on("click", ".ditems",function(){
        var wx1 = $(this).children(".dw1").text();
        var wx2 = $(this).children(".dw2").text();
        var wx3 = $(this).children(".dw3").text();
        var wx4 = $(this).children(".dw4").text();

        $("#meanings").html("<b>"+wx1 +"</b>," + wx2 +"<br/>" + "<b>" + wx3 + "</b><br/>" + wx4);  //***************** */
    })
}

// Saveが押されたとき
function wordsave() {
    var tword = $("#word").val()
    if(tword.length < 1){
        popbox.closebox();
    }
    var meantext = $("#meanings").html();
    meantext = meantext.replace(/<div>/g,"");
    meantext = meantext.replace(/(<p>|<\/p>)/g, '<br>');  // ContentEditableをDIVで使うとCRの入力で＜P>が挿入されるので、それを取る
    $wordnote.html(tword + "<span>"+ meantext + "</span>");
    if(meantext.length >1){
        $wordnote.removeClass().addClass("note1")
    }
    else
        $wordnote.removeClass().addClass("note0");
        wcancel();
    }

    function wcancel(){
        $("#meanings").attr("contentEditable", false);
        popbox.closebox();
    }
    

function insertButtons(){
    // if(!$(".lbtn1").length){
    if(!$("button").hasClass("lbtn1")){
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

/***************************Editing mode ********************************************* */
// spnote, translation editmode
var $trans;         // editingmode 終了のときに使うので、Globalにしておく
var $transp;
var isEditing = false;
var changed = false;        // 変更があったかどうか

function editingmode(){
    if(isEditing == true)
        return;
    $transp = $(this); // transpはspnoteかtranslationが入る
    $trans = $(this).children(".innerdiv");

    $("#div.innerdic").parent().off();      // 他の編集が始まらないように
    attachemenubar($(this)) // $(this) = spnote/ translation
    // initColorPicker();
    initEditableFunc();
    $("#marker").on("click", drawMarker);
    $("#editclose").on("click",editFinish);
    $transp.on("keydown", macronx); // initeditableの中ではセットしない
    var what = $trans.html().trim();
    var what2 = $trans.text().trim();
    if (what == '<p></p>' || what2 == '') //innerdivの中身をチェック
        $trans.empty().append('<p><br></p>');
    $transp.css('background-color', '#fffacd');
    $transp.attr("contentEditable", true);
    isEditing = true;
}

function editFinish() {
    $("#menubarx").remove();
    $transp.css("background-color", "#fffff8");
    $trans.parent().attr("contentEditable", false);
    $("div.innerdiv").parent().on("dblclick", editingmode);
    // setActionFlag();
    isEditing = false;
    changed = true;
}

/************************data save************************************************ */

function uploadsave(){
    var fname = $("#pathname").text();
    if(fname.length == 0){
        // alert("このページはサーバーに保存できません。")
        msgbox("このページはサーバーにSaveできません。\r\nダウンロードしてください。")
    }else {
        datasave(fname);
    }
}

function datasave(fname){
    var data = getpagedata();
    if(data.length >0){
        // postForm("/upload2",{"data": data, "fname":fname});
        datauploadjson({"data": data, "fname":fname});
    }
    else{
        msgbox("編集中です。")
    }
}

// data = {"data": data, "fname":fname}
function datauploadjson(data){
    // var sdata = JSON.stringify(data,null,"");
    $.post("/uploadjson",data,null, "script")
    .done(function(rvalue) {
        if(rvalue != "null"){
            changed = false;
            msgbox("Saveしました。")
        }
        else {
            msgbox("Save中エラーがありました。")
        }
    })
}


// upload and change content readonly to download
function downloadsave(){
    var data = getpagedata();
    postForm("/download",{"data": data, "fname":"download"});
}

// saveなしで閉じようとしたとき
function calledwhenunload(e){
    if(changed){
        e.preventDefault();
        e.returnValue = " ";
    } 
}


function getpagedata() {
    if (isEditing == true){
        // alert("編集を終了してください。")
        return "";
    }
    // $(".lbtn1").remove();
    // $(".lbtn2").remove();
    // $(".boxright").remove(); // boxright(save, download)は残す。roprefixでhideする
    //$("#pathname").remove();    // ?????????? 残すか？
    // $(".wordbox").remove();
    // $("#msgbox").remove();
    

    var allcontent = $("#outer").html();
    var pre = `<body id="bodyportion">
        <div id="outer">
    `
    var trail = `</div>
    </body>
    <html>
    `
    return pre + allcontent + trail
}

function attachemenubar(obj) {
    var thtml = `
    <div id="menubarx" class="titlex" style="width:100%;">
    <div class="fleft" style="padding-left:12px;">編集操作</div>
    <div class="fleft" style="padding-left:24px;">
        <input id="boldb" class="cbutton" type="button" title="bold" value="B">
        <input id="ulineb" class="cbutton" type="button" title="UnderLine" value="U">
        <input id="italicb" class="cbutton" type="button" title="Italic" value="I">
        <input id="soutb" class="cbutton" type="button" title="StrikeOut" value="S">
    </div>
    <div class="fleft" style="font-size:10px;padding-left:20px;">背景色</div>
    <button id ="marker" style="height:22px">マーカーペン</button>
    <div style="float:right;padding-right:12px;line-height:25px;">
    <input id="editclose" type="button" value="終了" style="font-weight: bold;height:22px;border-width:1px;">
    </div>
`
    obj.before(thtml);
}

function  initEditableFunc() {
    $("#boldb").on("click", function() { if(checkSelection()) {ornamentText("b")} });
    $("#ulineb").on("click", function() { if(checkSelection()) {ornamentText("u")} });
    $("#italicb").on("click", function() { if(checkSelection()) {ornamentText("i")} });
    $("#soutb").on("click", function() { if(checkSelection()) {ornamentText("s")} });

}
function checkSelection() {
    var selecFlag = window.getSelection().toString().length > 0;
    if (!selecFlag)
        alert("文字列が選択されていません");
    return selecFlag;
}

function ornamentText(s){
    var selection = window.getSelection();
    let content = document.createElement(s);
    const selectedRanage = selection.getRangeAt(0);
    try {
      selectedRanage.surroundContents(content)
    }catch(e) {alert(e)}
  }

  function drawMarker(){
    const span = document.createElement("span");
    span.style.background = "yellow";
    const selection = window.getSelection();
    const range = selection.getRangeAt(0);
    try{
        range.surroundContents(span);
    }catch(e) { alert(e)}
  }


var postForm = function (url, data) {
    var $form = $('<form/>', { 'action': url, 'method': 'post' });
    for (var key in data) {
        $form.append($('<input/>', { 'type': 'hidden', 'name': key, 'value': data[key] }));
     }
    $form.appendTo(document.body);
    var test = $form.serialize();
    $form.submit();
};


function postFormFile(url, files_){
    var $form = $('<form/>',{'action':url, 'method':'post',  'enctype' :"multipart/form-data"});
    var $inp = $('<input/>',{'type':'file', 'name':'file'});
    $inp.files = files_;
    $form.append($inp);
    $form.appendTo(document.body);{}
    $form.submit()
}

function postjson21(){
    event.preventDefault();
    postjson2();
}

// formデータをserializeしたもの（word=xxxx&name=yyyの形式）を送る
function postjson2() {
    // event.preventDefault();
    var testdata = $("#form1").serialize();
    console.log(testdata);
    $.post("/postjson2", testdata, null, "json")
        .done(function (data) {
            if (data == null) {
                $("#disp").text("Nothing Found");
            }
            else {
                $("#disp").text(" ");
            }
            const data2 = JSON.stringify(data);
            // $("#message").text(data2);
            var datajson = $.parseJSON(data2);
            for (var i in datajson) {
                var result = tempdiv.replace("%s0",datajson[i]);        // 全体をdisplay:noneの場所に入れておく（CrreationTableのため）
                result = result.replace("%s1", datajson[i].Dic.WORD);
                result = result.replace("%s2", datajson[i].Dic.NOTES);
                result = result.replace("%s3", datajson[i].Dic.MEANINGS);
                result = result.replace("%s4", datajson[i].Attr);
                $("#disp").append(result);
            }
        })
}

  // macronをCTRL/ALTで、その直前の文字につける
// option: trueのときはEnter, leftarrow, rightarrow, uparrow, downarrowでカスタムEventを起こす
// ctrl/alt でその前のKey入力を長音にする

function macronx(eo) {
    var ctrl = eo.ctrlKey;
    var shift = eo.shiftKey;
    var alt = eo.altKey;
    if (alt && ctrl) {
        var selection = window.getSelection(); // rangeを得るためSelectionを作成
        var range = selection.getRangeAt(0); // 選択範囲は大概一つ。この場合カーソルのみで、collaspeした状態
        var sstart = range.startOffset; // 現在の選択範囲の先頭（カーソル位置）
        if (sstart > 0) {
            var endContainer = range.endContainer;
            var value = endContainer.nodeValue; // カーソルのある要素（text)の値
            if (value != null) {
                var ch = value.substr(sstart - 1, 1); // カーソルｍｐ前一文字
                var newch = checkprev(ch);
                if (newch != null) {
                    var str = null;
                    if (sstart == 1) {
                        str = newch + value.substr(sstart);
                    } else if (sstart == value.length)
                        str = value.substr(0, value.length - 1) + newch;
                    else {
                        str = value.substr(0, sstart - 1) + newch + value.substr(sstart, value.length - sstart);
                    }
                    endContainer.nodeValue = str; // 新しい値をセット
                    var range2 = document.createRange();
                    range2.setStart(endContainer, sstart);
                    range2.setEnd(endContainer, sstart);
                    selection.removeAllRanges();
                    selection.addRange(range2); // 新しいカーソル位置
                    return false;
                }
            }
        }
    }
}

function checkprev(c) {
    if (c == "a")
        return "ā";
    else if (c == "i")
        return "ī";
    else if (c == "u")
        return "ū";
    else if (c == "e")
        return "ē";
    else if (c == "o")
        return "ō";
    else if (c == "A")
        return "Ā";
    else if (c == "I")
        return "Ī";
    else if (c == "U")
        return "Ū";
    else if (c == "E")
        return "Ē";
    else if (c == "O")
        return "Ō";
    else
        return null;
}
