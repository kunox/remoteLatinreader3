package inflection

import (
	Ending "latindictionary/endinglist"
	_ "latindictionary/fulldict"
	_ "latindictionary/stemlist"
	S "latindictionary/structs"
	"strings"
	"strconv"
)

// 指定のSLDのグループテーブルから、指定のKeyに合う、項目とendingwordの一致するものを返す。
// テーブルのendingwordはMacron付きなので外して比較
// Attributesは合致した場合のテーブル上の位置（例えば (o,0)とか（１，０）
//（０，０）の語尾が入るが、テーブルのKey==0 の場合は、基本形そのものが入る。基本形はbasicFormに格納
// PRON54 Suiの場合はGen　Sが基本形として返る（!sui)
// Finding2にはAttributesとBasicFormがセットされるがその他はセットされない

func Findmatch(sld S.StemListData, endingword string, sword1 string) (S.Findings, bool) {
	sword := strings.ToLower(sword1)
	// ainfo := strings.Split(sld.Additional, ",")
	// tname := sld.Group		// tableサーチ用
	// bname := ""				// basicformを探すため
	// finds := S.Findings{}
	// enclitic := ""			//encliticがある場合のbasicFormのため

	switch true{
	case sld.Group[0] == 'N' && sld.Group[1] != 'U' : 
		return caseNoun(sld, endingword, sword)
	case strings.HasPrefix(sld.Group, "ADJ") : 
		return caseADJ(sld, endingword, sword)
	case strings.HasPrefix(sld.Group, "NUM") :
		return caseNum(sld, endingword, sword)
	case strings.HasPrefix(sld.Group, "PRON") :
		return casePron(sld, endingword, sword)
	case strings.HasPrefix(sld.Group, "V") :
		return caseV(sld, endingword, sword)
	default : return S.Findings{}, false
	}
}



func caseNoun(sld S.StemListData, endingword string, sword string) (S.Findings, bool){
	ainfo, tname := getparam(sld)
	if len(ainfo) < 1{
		return S.Findings{}, false
	}
	gender := ainfo[1]
	if gender== "M" || gender == "F" || gender =="X" {
		gender = "C"
	}
	tname = tname + " " + gender
	finds, ok := tmatch(sld, tname, endingword,sword, "")
	if ok {
		finds.Parts = "N"
		return finds, true
	} else {
		return	finds, false
	}
}

func caseADJ(sld S.StemListData, endingword string, sword string) (S.Findings, bool){
	_, tname := getparam(sld)
	bname := ""
	if sld.Group[3] != '0'{
		bname = tname		// tnameは変わる可能性がある（比較級など）
	}
	if sld.Key ==3 {
		tname = "ADJ00 200"
	}
	if sld.Key == 4 {
		tname = "ADJ00 100"
	}
	finds, ok := tmatch(sld, tname, endingword,sword, bname)
	if ok {
		if sld.Key == 3 {
			finds.Note = "比較級"
			finds.AttrNum = "200"
		} else if sld.Key == 4 {
			finds.Note = "最上級"
			finds.AttrNum = "100"
		}
		finds.Parts ="ADJ"
		return finds,true
	} else {
		return finds,false
	}

}

// NUMは基本的にfulldictにない
// key ==1 基数、key==4 数副詞（これはstemlistにしかないのでFindingsを作成する）
// key==2序数はNUM00 2, key == 3配分数はNUM00 3を使用
func caseNum(sld S.StemListData, endingword string, sword string) (S.Findings, bool){
	ainfo, tname := getparam(sld)
	if sld.Key == 4 {
		ret := S.Findings{}
		ret.BasicForm = sword
		ret.Note = "数副詞 " + ainfo[1]
		ret.AttrNum = "4"
		ret.Parts = "NUM"
		if (tname == "NUM20" || tname =="NUM14") {
			if endingword != "ies" {
				return ret, false
			} else {
					ret.BasicForm = sword + "ies"		//Stemlistにもsemel,bis, ter以外ではiesを除いたStemしかない
			}
		}else if len(endingword) > 0{		// semel, bis, terはendingがないはず
				return ret, false
		}
		return ret, true
	
	} else if sld.Key == 2{		// 序数
		tname = "NUM00 2"
	} else if sld.Key == 3 {	// 配分数
		tname = "NUM00 3"
	}
	// 基数の場合はtnameはそのまま
	finds, ok := tmatch(sld, tname, endingword,sword,"")
	if ok {
		if sld.Key == 0 || sld.Key == 1 {
			finds.Note = ainfo[1]
			finds.AttrNum = "1"	// 基数
		} else if sld.Key == 2 {
			finds.Note = "序数 " + ainfo[1]
			finds.Parts = "NUM00 X," + sword+ "," +sword + ",,"	// fullDictにないのでここで作っておく
		} else if sld.Key== 3 {
			finds.Note = "配分数 " + ainfo[1]
			finds.AttrNum = "3"
			finds.Parts = "NUM00 X," + sword+ "," +sword + "," + sword + ","
		}
		return finds, true
	}
	return finds, false
}

func casePron(sld S.StemListData, endingword string, sword string) (S.Findings, bool){
	if len(sword) == 0 {
		return S.Findings{}, false
	}
	_, tname := getparam(sld)
	if sld.Group[4] != '2' {// Pron2X以外
		finds, ok := tmatch(sld, tname, endingword, sword, "")
		if ok {
			if sld.Group == "PRON54" {
				finds.BasicForm = "sui"
			} else if sld.Group == "PRON53" {
				finds.BasicForm = sword + "os"	 // vos nos　PRON53は単数がないので正しくbasicformが戻ってこない
			}
			finds.Parts = "PRON"
			return finds, true
		} else {
			return finds,false
		}
	}
	// PRON2X
	enclitic := ""
	if sld.Group=="PRON21" || sld.Group == "PRON24" {	// 21はAdjective 24はIndef
		if sword + endingword == "quisque" || sword + endingword == "quispiam" {
			ret := S.Findings{}
			ret.BasicForm = sword + endingword
			ret.Parts = "PRON"
			ret.Attribute = append(ret.Attribute, S.Cordinate{X:0,Y:0})
			return ret,true
		}
		enclitic = getEnclitic21(sword + endingword)
		if sld.Group == "PRON21" {
			tname = "PRON11"
		}
	} else if sld.Group == "PRON22" {
		enclitic  = getEnclitic22(sword + endingword)
		tname = "PRON12"
	} else if sld.Group == "PRON25" || sld.Group == "PRON26" {
		enclitic = getEnclitic25(sword + endingword)
	}
	if len(enclitic) > 0 && len(endingword) >= len(enclitic) {
		endingword = endingword[0:len(endingword) - len(enclitic)]
		finds, ok := tmatch(sld, tname, endingword, sword, "")
		if ok {
			finds.BasicForm = finds.BasicForm + enclitic
			finds.Parts = "PRON"
			return finds, true
		} else {
			return finds, false
		}
	}
	return S.Findings{}, false
}

func caseV(sld S.StemListData, endingword string, sword string) (S.Findings, bool){
	_, tname := getparam(sld)
	ret, ok := tmatch(sld, tname, endingword, sword, "")
	if ok {
		ret.Parts = "V"
		return ret, true
	}
	// 分詞の場合
	bname := tname
	xkind := "0"
	c := tname[1]
	var newname string
	switch c {
	case '1', '2', '3' : newname = "VPAR" + string(c) + "0"
	case '7' : newname = "VPAR62" 	// V71, V72はVPAR62を使用する
	default : newname = "VPAR" + tname[1:]
	}
	ret2, ok := tmatch(sld, newname, endingword, sword, bname)
	if ok {
		if ret2.Attribute[0].Y > 5 {
			xkind = "5"			// 動形容詞
		}
		// それ以外ではxkind =0, 現在分詞
		ret2.AttrNum = "PPL " +xkind	//pplの後スペース
		ret2.Parts = "V"
		return ret2, true
	}
	ret2, ok = tmatch(sld, "VPAR00", endingword, sword, bname)
	if ok {
		if ret2.Attribute[0].Y > 5 {		// VPAR00 受動は完了分詞
			xkind = "4"
		}	else {
			xkind= "2"	// VPAR00 能動は未来分詞
		}
		ret2.AttrNum = "PPL " + xkind
		ret2.Parts = "V"
		return ret2, true
	}
	return ret, false
}

// ainfo, tname
func getparam(sld S.StemListData) ([]string, string) {
	return strings.Split(sld.Additional, ","), sld.Group 
}




// 指定のSLDのグループテーブル(tname)から、指定のKeyに合う、項目とendingwordの一致するものを返す。
// 一致するものがあれば、findings2のAttributesに一致した座標（０，０）などを入れて返す。
// 更には単語の基本形がセットされる。但し、これはbascinameTable(basicnameTable)の(０，０)の語尾をつけるだけの標準的なものである。
//      basicnameの例外的な処理は外側で行う。
// 単語の基本系を求めるためのテーブルがtableと異なる場合は、basicnametableに名前を求めるテーブル名が入っている。
// 例えば形容詞比較級の場合、変化はADJ00 200などのテーブルで調べるが、基本形は元のADJ11のテーブルで調べる

func tmatch(sld S.StemListData, tname string, endingword string, sword string,bname string) (S.Findings, bool) {
	ret := new(S.Findings)
	etable, ok := Ending.Etable[tname]
	if ok {
		zeroflag := false
		line1 := strings.Split(etable[0], "\t")
		item00 := strings.Split(line1[0], " ")	// 最初の行の最初のアイテム「数字 語尾」 後で使用する
		if len(strings.Trim(line1[0], " ")) == 0 {
			item00 = strings.Split(line1[3]," ")		 // 複数のみ(duoなど)最初のItemがない場合は、複数１人称を使用する
		}
		comp := endingword
		if etable[0][0] == '0'{						// 最初の行の最初の文字
			if tname == "PRON52" && len(endingword) == 0 {
				return *ret, false			// tu と t + uの両方がOKになるのを防ぐ。tuのみが問題
			}
			comp = sword + endingword	// ０の場合はStem+EndingがTableにある
			zeroflag = true
		}
		maxCount := len(etable)
		if maxCount == 21 {
			maxCount -=1	 // VはVPAR,V00も調べるので最後の行PPLは調べなくて良い
		}
		for y:=0; y < maxCount; y++ {
			items := strings.Split(etable[y], "\t")
			for x:=0; x < len(items)  ; x++ {
				eachitem :=strings.Split(items[x], " ")	//各アイテムは、1 o のように数字＋語尾
				if len(eachitem) < 2 {
					continue
				}
				if sld.Key != 0 && eachitem[0] != "0" {
					if eachitem[0] != strconv.Itoa(int(sld.Key)){
						continue
					}
				}
				if sld.Key ==0 {	// sld.keyが０は１，２はOKだが３，４はだめ
					if eachitem[0] == "3" || eachitem[0] =="4" {
						continue
					}
				}
				temp := stripmacron(eachitem[1])
				temp2 := strings.Split(temp, ",")
				for _, entry := range temp2{
					if comp == entry {
						ret.Attribute = append(ret.Attribute, S.Cordinate{X:int16(x), Y:int16(y)})
						break
					}
				}
			}
		}
		if len(ret.Attribute) == 0 {
			return *ret, false		// 該当なし
		} else {
			// 該当ありのときだけbasicformを求める
			if len(bname) > 0 {					// bname用に別のGroup名が指定された場合
				etable, ok = Ending.Etable[bname]
				if !ok  {
					return *ret, false
				}
				// item00も再度取得
				first := strings.Split(etable[0], "\t")
				item00 = strings.Split(first[0], " ")	// 1行目の最初のアイテムのスペースで区切られた「数字＋語尾」
				if len(item00[0])==0 {
					item00 = strings.Split(first[3], " ")	// 複数形のみなど（０，０）がない場合、４番目の要素
				}
				if item00[0] == "0"{
					zeroflag = true
				}
			}
			if !zeroflag {
				basicstem := sword			// 通常はswordがstemになる
				stemb := strings.Split(sld.Additional, ",")	// sldのAdditionalに原形がある場合あり
				if len(stemb) == 3 && len(stemb[2]) > 0 {
					basicstem = stemb[2]
				}
				ret.BasicForm = basicstem + stripmacron(item00[1])
			} else {
				ret.BasicForm = stripmacron(item00[1])
			}
			return *ret, true
		}
	}
	return *ret, false
}


func stripmacron(sin string) string{
	runein := []rune(sin)
	outdata := make([]rune, len(runein))
	for pos, ch := range runein {
		outdata[pos] = toPlain(ch)	// outdataはruneの文字列
	}
	s := string(outdata)		// 	stringに変換して
	barray := []byte(s)			// そのstringをByte列に戻し
	return string(barray)	// Byte列をstringに戻す
}


func toPlain(ch rune) rune {
	switch ch {
	case 'ā': return 'a'
	case 'ē' : return 'e'
	case 'ī' : return 'i'
	case 'ō' : return 'o'
	case 'ū' : return 'u'
	case  'A' : return 'A'
	case 'Ē' : return 'E'
	case  'Ī' : return 'I'
	case 'Ō' : return 'O'
	case 'Ū' : return 'U'
	case 'ȳ' : return 'y'
	case  'Ȳ' : return 'Y'
	default : return ch
	}
}

func getEnclitic21(sin string) string {
	forp21 := []string {"libet", "lubet", "cumque", "cunque", "vis", "piam", "cum", "que", "nam" }
	return ecommon(sin, forp21)
}

func getEnclitic22(sin string) string {
	forp22 := []string {"quam", "nam"}
	return ecommon(sin, forp22)
}

func getEnclitic25(sin string) string {
	forp25 := []string{"dam"}
	return ecommon(sin, forp25)
}

func ecommon(sin string, selection []string) string {
	for _, s :=range selection {
		if strings.HasSuffix(sin, s) {
			return s
		}
	}
	return ""
}