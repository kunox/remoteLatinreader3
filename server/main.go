package main

import (
	_ "embed"
	_ "encoding/json"
	"fmt"
	_ "go/format"
	"io/ioutil"
	Dictionary "latindictionary/latdicmainF"
	S "latindictionary/structs"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	 "path/filepath"
	"github.com/gin-gonic/gin"
)

//go:embed assets/roprefix.txt
var roprefix string

// var bookindex indexstruct
var registered []string =[]string {"senatus", "kunox"}

func main() {
	// var word string
	gin.SetMode(gin.ReleaseMode)
	if !Dictionary.Init(){
		fmt.Println("Error")
		return
	}
	ejinit()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	// router.Static("/images", "./images")
	router.Static("assets", "./assets")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "bookselection.html", nil)
	})


	router.POST("/setUsername", func(ctx *gin.Context){
		uname, _ := ctx.GetPostForm("uname")
		found := false
		for _,name := range registered {
			if name == uname {
				ctx.String(200, name)
				found = true
				break;
			}
		}
		if !found{
			ctx.String(200,"")
		}
	})

	router.POST("/getBookIndex", func(ctx *gin.Context){
		folder, _ := ctx.GetPostForm("fname")
		user, _ := ctx.GetPostForm("uname")
		indextext, bindex := createBookIndex(folder)
		if len(user) > 0 {
			indextext = checkUserBooks(indextext, user,bindex)
		}
		ctx.String(200, indextext)
	})



	router.POST("/showpage", func(ctx *gin.Context) {
		vol,_ := ctx.GetPostForm("vol")
		chap, _ := ctx.GetPostForm("fname")
		foldername, _ := ctx.GetPostForm("bookname")
		ext, _ := ctx.GetPostForm("ext")
		uname, _ := ctx.GetPostForm("uname")
		bookshelf,_ := ctx.GetPostForm("cell")		// bookshelf="" はbooksから、
		if bookshelf == ""{
			bookshelf = "books"
		}
		fpath := filepath.Join("./"+bookshelf, foldername,vol, chap+ext)
		// fname := "./books/Vol" + vol + "/Chap " +chap + ".htm"
		data, _ := os.ReadFile(fpath)
		fpath2 := filepath.Join("./"+uname, foldername,vol, chap+ext)
		// fpath2はhtmlに記録されるPath 登録ユーザーが新しいページを開いたときはbooksから
		//		ロードするが、格納はuser folderにしなければならない
		//		unameが""の場合はSaveできないようにする！！
		if len(uname) == 0 {
			fpath2 = ""					// fpath2が””のときはSaveできない
		}
		data2 := extractbody(string(data), fpath2);		//
		ctx.Header("Content-Type", "text/html")
		ctx.String(200,data2)
	})

	router.POST("/upload2", func(c *gin.Context){
		data, _ := c.GetPostForm("data")
		fname, _ := c.GetPostForm("fname")

		fmt.Print(fname)
		c.Header("Content-Type", "text/html")
		data2 := cretaeRO(data)
		ioutil.WriteFile(fname,[]byte(data2), os.ModePerm)
		// c.String(200,string(data2))
		// fmt.Print(string(data2))
	})

	router.POST("/upload3", func(c *gin.Context){
		data, _ := c.GetPostForm("data")
		flag, _ := c.GetPostForm("flag")

		fmt.Print(flag)
		c.Header("Content-Type", "text/html")
		data2 :=""
		if flag == "0"  {
			data2 = extractbody(data, "")
		}else {
			temp := createlr(data,flag)
			fmt.Print(string(temp))
			data2 = extractbody(temp, "")
		}
		c.String(200,string(data2))
	})


	router.POST("/uploadjson", func(c *gin.Context){
		data, _ := c.GetPostForm("data")
		fname, _ := c.GetPostForm("fname")
		// c.Header("Content-Type", "text/html")
		data2 := cretaeRO(data)
		err := os.WriteFile(fname,[]byte(data2), os.ModePerm)
		if err != nil {
			c.JSON(http.StatusOK, nil)		// errorのときはnilを返す
		}else {
			c.String(http.StatusOK, "OK")
		}
	})

	router.POST("/download", func(c *gin.Context){
		data, _ := c.GetPostForm("data")
		fname, _ := c.GetPostForm("fname")
		data2 := cretaeRO(data)
		c.Header("Content-Disposition", "attachment; filename=" + fname +".html")
		c.Header("Content-Type", "text/html")
		c.String(200,string(data2))

	})


	router.POST("/postjson2", func(c *gin.Context){
		var result []S.ExdicContent
		word, _ := c.GetPostForm("word")
		result, ok := Dictionary.Findword(word)
		if !ok {
			// fmt.Println("notthing found")
			result = nil
		}
		message := word 
		fmt.Print(message)
		// c.JSONはstructをうまくrenderして、httpで送り返す。
		// 受け側ではJSONデータなのでstringfyすること
		c.JSON(http.StatusOK, result)
	})

	router.POST("/englishDic", func(c *gin.Context) {
		eword, _ := c.GetPostForm("eword")
		result := ejtranlate(eword)
		c.String(200, result)
	})

	router.POST("/ifexists", func(ctx *gin.Context){
		vol,_ := ctx.GetPostForm("vol")
		chap, _ := ctx.GetPostForm("fname")
		foldername, _ := ctx.GetPostForm("bookname")
		ext, _ := ctx.GetPostForm("ext")
		uname, _ := ctx.GetPostForm("uname")
		fpath := filepath.Join("./"+uname, foldername,vol, chap+ext)
		if f, err :=os.Stat(fpath) ; os.IsNotExist(err) || f.IsDir() {
			ctx.String(200,"no")
		}else {
			ctx.String(200,"yes")
		}
	})

	router.POST("/tblcreation", func(ctx *gin.Context){
		var inexdc S.ExdicContent
		ctx.BindJSON(&inexdc)
		table := Dictionary.CreateTable(inexdc)
		if  len(table) == 0 {
			ctx.String(200, "")
		}
		var tname string
		// if len(table) == 2 {
		// 	tname = "NounTable.htm"
		// } else if len(table) <= 7 {
		// 	tname ="AdjTable.htm"
		// } else {
		// 	tname = "VerbTable.htm"
		// }
		switch {
		case len(table) == 2 : tname = "NounTable.htm"
		case len(table[0]) <= 7 : tname ="AdjTable.htm"
		default: 	tname = "VerbTable.htm"
		}
		htmlstring := tablecommon(tname,  table)
		// fmt.Println(htmlstring)
		ctx.String(200,htmlstring)
	})

	router.Run(":8080")
}




func tablecommon(tname string, table [][]string) string{
	fpath := filepath.Join("./tables", tname)
	tdata, _ := os.ReadFile(fpath)
	tdata2 := string(tdata)
	for  x:=0; x < len(table); x++ {
		for y:=0; y < len(table[0]) ; y++ {
			snew := table[x][y]
			sold := fmt.Sprintf("{%d,%d}",x,y)
			tdata2 = strings.Replace(tdata2, sold, snew, 1)
		}
	}
	return tdata2
}


func createBookIndex(folder string) (string, indexstruct) {
	var data []byte
	if folder == "FabulaeFaciles" {
		data,_ = os.ReadFile(filepath.Join("./books",folder, "index.htm"))
		bindex := indexstruct{4, "FabulaeFaciles","htm",nil}
		bindex.vols = append(bindex.vols, volstruct{"Perseus", "p", 11})
		bindex.vols = append(bindex.vols, volstruct{"Hercules", "p", 44})
		bindex.vols = append(bindex.vols, volstruct{"Argonauts", "p", 23})
		bindex.vols = append(bindex.vols, volstruct{"Ulysses", "p", 19})
		// ctx.String(200, string(data))
		return string(data), bindex
	}else {

		data,_ = os.ReadFile(filepath.Join("./books",folder,"indexconfig.txt"))
		array := strings.Split(string(data), "\r\n")
		bindex := indexstruct{}
		n, _ := strconv.Atoi(array[0])
		bindex.number = n
		bindex.bookname = array[1]
		bindex.extension= array[2]
		for x:=0; x < n; x++ {
			varr := strings.Split(array[x + 3],"#")
			num,_ := strconv.Atoi(varr[2])
			bindex.vols = append(bindex.vols, volstruct{varr[0], varr[1], num})
		}

		var s strings.Builder
		s.WriteString("<div id=\"outer\">\r\n")
		s.WriteString("<div id='extension' style='display:none';>" + bindex.extension + "</div>\r\n")

		s.WriteString(`<div id="header">` + bindex.bookname + "</div>\r\n")
		for i:=0; i < len(bindex.vols); i++ {
			lnum := strconv.Itoa(i+1)
			s.WriteString(`<div id="l` + lnum + `" class="tabs">` + bindex.vols[i].name + "</div>\r\n")
			s.WriteString(`<div id="tab` + lnum +"\" class=\"tab\">\r\n")
			for j:=0; j < bindex.vols[i].number ; j++ {
				s.WriteString("<div>"+ bindex.vols[i].chapsec + strconv.Itoa(j+1)+ "</div>\r\n")
			}
			s.WriteString("</div>")
		}
		s.WriteString("</div>\r\n</div>\r\n")

		// ctx.String(200, s.String())
		return s.String(),bindex
	}
}


// User Booksを走査して存在するものをマークつけ
func checkUserBooks(plainindex string, uname string, bindex indexstruct) string {
	libernum := len(bindex.vols)
	bookfolder := filepath.Join("./", uname, bindex.bookname)
	if f, err :=os.Stat(bookfolder); os.IsNotExist(err) || !f.IsDir() {		// 指定Bookがunameの中にない場合
		createBookFolders(uname, bindex)
		return plainindex
	}
	ptextblock := strings.Split(plainindex,"<div id=\"tab")
	for i:=0; i < libernum; i++ {
		folderfullpath := filepath.Join("./", uname, bindex.bookname,bindex.vols[i].name)
		fnames, _ := dirwalk2(folderfullpath)
		for _,fname := range fnames {
			ptextblock[i+1] = strings.Replace(ptextblock[i+1], "<div>"+fname, "<div class = 'exist'>" + fname, 1)
		}
		// fmt.Print(folderfullpath)
	}
	fmt.Print(len(ptextblock))
	var s strings.Builder
	s.WriteString(ptextblock[0])
	for i:=1; i < len(ptextblock); i++ {
		 s.WriteString("<div id=\"tab" + ptextblock[i])
	}
	return s.String()
}


func dirwalk2(folder string) ([]string, int){
	direntries, err := os.ReadDir(folder)
	if err != nil {
		return nil, 0
	}
	retfname := []string {}
	retint := 0
	for _,dentry := range direntries {
		if dentry.IsDir() {
			continue
		} else {
			f := strings.Split(dentry.Name(), ".")
			if f[1]== "html" || f[1] == "htm"{
				retfname = append(retfname,f[0])
				retint++
			}
		}
	}
	return retfname, retint
}


// userBookにbindex	で記述されるfolderを新規に作成する
func createBookFolders(uname string, bindex indexstruct) error {
	if err := os.Mkdir(filepath.Join("./", uname, bindex.bookname),0777); err !=nil {
		return err
	}
	for i:=0; i < bindex.number; i++ {
		if err := os.Mkdir(filepath.Join("./",uname, bindex.bookname, bindex.vols[i].name),0777); err != nil {
			return err
		}
	}
	return nil
}


	// yourbooksのFoldernameの中を探索する
	// DIRは無視する Fileの数も返す
// func dirwalk(folder string) (string, int) {
// 	files, err := ioutil.ReadDir(folder)
// 	if err != nil{
// 		panic(err)
// 	}
// 	var paths strings.Builder
// 	retnum := 0
// 	for _,file := range files {
// 		if file.IsDir(){
// 			// paths.WriteString("%" + file.Name() + "#")
// 			// p := filepath.Join(folder, file.Name())
// 			// fs := dirwalk(p)
// 			// // fmt.Print(fs)
// 			// paths.WriteString(fs)
// 		}else {
// 			f := strings.Split(file.Name(), ".")
// 			if f[1] == "html" || f[1] == "htm"{
// 				paths.WriteString(f[0]+ "#")
// 				retnum++
// 			}
// 		}
// 	}
// 	return paths.String(), retnum
// }

func extractbody(indata string , fname string) string {
	content := indata
	
	startpos2 := strings.Index(content, `<div id="header"`)
	if startpos2 == -1 {
		startpos2 = search2nddiv(content)
	}
	
	endpos := strings.LastIndex(content, "</body")
	newcontent := content[startpos2:endpos]
	var prefix = `
	<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=EmulateIE11">
    <meta name='application-name' content='Editable HTML'>
    <meta name="viewport" content="width=device-width">
    <title>Editable</title>
    <link href="./assets/css/chapcss.css" rel="stylesheet" type="text/css">
	<link href="./assets/css/wordbox.css" rel="stylesheet" type="text/css">
    <link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/redmond/jquery-ui.css">
	<link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/spectrum/1.8.0/spectrum.min.css">

    <script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
    <script type="text/javascript" src="https://code.jquery.com/ui/1.12.0/jquery-ui.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/spectrum/1.8.0/spectrum.min.js"></script>
	
	<script src="./assets/scripts/chapjs.js" type="text/javascript"></script>
    <script src ="./assets/scripts/wordbox.js" type="text/javascript"></script>

		</head>
		<div id="pathname" style="display:none">%</div>
		<body id = "bodyportion">
		<div id="outer">
		<a name="begin"></a>
	`
	var postfix = `
	<div id="dict">
	<form id="form1">
		<input type="text" name="word" id="keyin">
		<button id="wbtn">送信</button>
	</form>
	<div id="disp">
		<div class="ditems" >
			<p class="exdc"></p>
			<span class="dw1"></span>  <span class="dw2"></span>
			<p class="dw3"></p>
			<p class="dw3"></p>
		</div>
	</div>
</div>
</body>
</html>
	`
	prefix2 := strings.Replace(prefix, "%", fname, 1)
	editabletext := prefix2 + newcontent + postfix
	// ioutil.WriteFile("./tempfile.txt",[]byte(editabletext),os.ModePerm)
	return editabletext
}

func cretaeRO(data string) string {
	return roprefix + data
}

func createlr(data string,title string) string {
	var prefix = `
<body id = "bodyportion">
<a name="begin"></a>
<div id="outer">
<div id="header">
<h1 style="color: rgb(0, 0, 0);">%s</h1>
</div>
	`
	var prefix2 = strings.Replace(prefix,"%s", title,1)
	var ret0  strings.Builder
	ret0.WriteString(prefix2)
	reg := "\r\n|\n"
	sarray := regexp.MustCompile(reg).Split(data,-1)
	for i, v:=range sarray {
		ret0.WriteString(latinconv(i, v))
	}
	ret0.WriteString(`
	<div id="footer">
	<a href="#begin"> ページの最初に </a>
	</div>
	</div>
	</body>
	</html>`)
	return ret0.String()
}

func latinconv(i int, c string)string {
	var ret strings.Builder
	ret.WriteString(fmt.Sprintf(`
	<div class="spacer">Line%d</div>
	`,i+1))
	ret.WriteString(`
	<div class="latin">
	<p class="pop" style="margin-left:10px;" name="latintext>
`)

	sarray := strings.Split(c, " ")
	for i:=0; i < len(sarray); i++{
		ret.WriteString(fmt.Sprintf(`
		<a class="note0">%s<span></span></a>`,sarray[i]))
	}
	ret.WriteString(`
</p>
</div>
<div class="spnote"  >
<div class="innerdiv" spellcheck="false">
<p></br></p>
</div>
</div>
<div class="translation">
<div class="innerdiv" spellcheck="false">
<p></br></p>
</div>
</div>
`)
	return ret.String()
}

func search2nddiv(content string) int {
	start1 := strings.Index(content, "<div id")
	newcontent := content[start1 +1:start1+100]
	start2 := strings.Index(newcontent, "<div id")
	return start1 + start2 + 1
}