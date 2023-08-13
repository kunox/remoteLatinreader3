package  tablecreation 

import (
	Ending "latindictionary/endinglist"
	"latindictionary/fulldict"
	S "latindictionary/structs"
	"strconv"
	"strings"
)

func GetVtable(exdc S.ExdicContent) [][]string {
	sa := strings.Split(exdc.Dic.PARTS, ",") // PARTS : V11 X,am,am,amav,amat
	sa2 := strings.Split(sa[0], " ")         // sa2 : V11 X (このXの部分にDEPなどが入る事あり）
	spnote := ""
	if len(sa2) > 1 {
		spnote = sa2[1]
	}
	tname := sa2[0]
	keys := setupkeys(exdc.Dic)
	ret := tableCommon(tname, keys)
	// 例外的な者たち
	if exdc.Dic.WORD == "mĕmĭni" {
		return exceptionMemini(exceptionalPERFDEF(ret))
	}
	if exdc.Dic.WORD == "absum" {
		return exceptionalABSUM(ret)
	}
	if strings.HasSuffix(exdc.Dic.WORD, "facio") || strings.HasSuffix(exdc.Dic.WORD, "făcĭo") {
		stem := exdc.Dic.WORD[0 : len(exdc.Dic.WORD)-5]
		return exceptionFacio(ret, stem)
	}
	if tname == "V62" {
		return exceptionV62(ret)
	}
	switch {
	case spnote == "DEP":
		return exceptionalDEP(ret)
	case spnote == "SEMIDEP":
		return exceptionalSEM(ret)
	case spnote == "PERFDEF":
		return exceptionalPERFDEF(ret)

	}
	return ret
}

//Vの現在、完了、未来分詞についての処理
// 基本的には形容詞と同じだが、テーブルサーチが異なる
// まずサーチするTableは例えばV21（permoti, permoveo)ならばVPAR20とVPAR00を捜す
//
// VPARxx, VPAR00は夫々、ACTIVEとPassiveで２つ分のテーブルがある。
// 現在分詞はVPARxのOffset０から
// 未来分詞受動（動形容詞）は同じテーブルをOffset６
// 能動未来分詞はVPAR00をOffset 0
// 完了分詞ははVPAR00をOffset 6
// exdcのattrNumには”PPL" + ｘの値が入っており、０：現在分詞、２：未来分詞能動、４：完了分詞　５：動形容詞
func GetPPLTable(exdc S.ExdicContent) [][]string {
	tname := ""
	num := strings.Split(exdc.AttrNum, " ")[1]
	if num == "0" || num == "5" {
		part := exdc.Dic.PARTS // V11 Trans, など
		tname = strings.Split(part, " ")[0]
		// switch {										// V1x, V2x, V3xは夫々VPAR10,20,30を使う。それ以外V５2などはそのままVPAR52
		// case tname[1]== '1' || tname[1] == '2' || tname[1] == '3' : tname = "VPAR" + string(tname[1]) + "0"
		// case tname[1] == '7' : tname = "VPAR62"
		// default : tname = "VPAR" + tname[1:]
		// }
		switch tname[1] { // V1x, V2x, V3xは夫々VPAR10,20,30を使う。それ以外V５2などはそのままVPAR52
		case '1', '2', '3':
			tname = "VPAR" + string(tname[1]) + "0"
		case '7':
			tname = "VPAR62"
		default:
			tname = "VPAR" + tname[1:]
		}
	} else if num == "2" || num == "4" {
		tname = string("VPAR00")
	}
	keys := setupkeys(exdc.Dic)
	offset := 0
	if num == "5" || num == "4" {
		offset = 6
	}
	return ppttablecreation(tname, keys, offset)
}

// xxfacioの場合 xxfioも調べる
func SpecialFacio(exdc S.ExdicContent) [][]string {
	dc := exdc.Dic
	table1 := GetVtable(exdc)
	wd := ""
	if len(dc.WORD) > 5 {
		wd = dc.WORD[0 : len(dc.WORD)-5]
	}
	dc2 := fulldict.Lookup(wd+"fio", false)
	if len(dc2) > 0 {
		exdc2 := S.ExdicContent{}
		exdc2.Dic = dc2[0] // ２つ以上はない
		exdc2.Attr = ""
		exdc2.AttrNum = ""
		table2 := GetVtable(exdc2)
		for y := 0; y < 3; y++ {
			for x := 0; x < 6; x++ {
				table1[x][y+6] = table2[x][y]
			}
		}
		for y := 9; y < 11; y++ {
			for x := 0; x < 6; x++ {
				table1[x][y+4] = table2[x][y]
			}
		}
		for y := 15; y < 17; y++ {
			for x := 0; x < 6; x++ {
				table1[x][y+2] = table2[x][y]
			}
		}
		table1[3][19] = table2[0][19]
	}
	return table1
}

//Adj
func GetADJTable(dcx S.DicContent) [][]string {
	sa := strings.Split(dcx.PARTS, " ") // N11 FT,ros,ros,,
	keys := setupkeys(dcx)
	return tableCommon(sa[0], keys)
}

//adj00
// fulldictにADJ00として載っていないもの、keyの値で最上級などを判定している
// dcの内容は例えば　ADJ11(altissimus)などと成っているので注意
// ADJ00として載っているものも、活用を調べる時点でexdiccontentのattrNumに100、200が入っているので
//      結局はここに来る。従って、keystemが4個あるか1個しかないかの違いだけ！！
func GetADJ00Table(exdc S.ExdicContent) [][]string {
	tname := "ADJ00 100"
	if exdc.AttrNum == "200" {
		tname = "ADJ00 200"
	}
	return tableCommon(tname, setupkeys(exdc.Dic))
}

// 序数、配分数の一覧表、
// 使用するテーブルはNUM00のみ、序数はkey=2, 配分数はKey=3を使用
// exdc.attrNum ="2"は序数、”３”が配分数
// Partsには作成した、stemが最初だけ入っている（ADJ00のものと同じ）
func NUM00Table(exdc S.ExdicContent) [][]string {
	tname := "NUM00 2"
	if exdc.AttrNum == "3" {
		tname = "NUM00 3"
	}
	return tableCommon(tname, setupkeys(exdc.Dic))
}

//名詞
// dcにはmacron付きが入っている
func GetNtable(dc S.DicContent) [][]string {
	sa := strings.Split(dc.PARTS, " ") // PART : N11 FT,ros,ros,,
	keys := strings.Split(sa[1], ",")  // FTの最初 本当のKeyは1Originでkeys[1],[2],[3],[4]に入っている
	gender := string(keys[0][0])       // FTのF
	tname := sa[0] + " " + gender
	if _, ok := Ending.Etable[tname]; !ok {
		tname = sa[0] + " C" // tableCommonで再度チェックされる
	}
	keys = setupkeys(dc)
	ret := tableCommon(tname, keys)
	if ret == nil {
		return nil // N9xではnilが返る。
	}
	nounexception(ret, dc)
	if len(ret) == 1 { // 1カラムしかない場合、追加する。（jesusのみ）
		newret := make([][]string, 2)
		for i := 0; i < len(ret[0]); i++ {
			newret[0][i] = ret[0][i]
			newret[1][i] = ""
		}
	}
	if dc.WORD[0] < 'a' && len(ret) > 1 {
		for i := 0; i < len(ret[0]); i++ {
			ret[1][i] = ""
		}
	}
	return ret
}

//PRON
// PRON11 ~ 14
// PRON51 52 53はego, tu, vestra, nostraなど代名詞分、基本的には名詞と同様で、GenderをC（３）として捜す

func GetPRONTable(dc S.DicContent) [][]string {
	tname := strings.Split(dc.PARTS, " ")[0]
	keys := setupkeys(dc)
	enclitic := ""
	switch tname {
	case "PRON21":
		{
			tname = "PRON11"
			enclitic = getEnclitic(dc)
		}
	case "PRON22":
		{
			tname = "PRON21"
			enclitic = getEnclitic(dc)
		}
	case "PRON24", "PRON25", "PRON26":
		enclitic = getEnclitic(dc)
	}
	ret := tableCommon2(tname, keys, enclitic)
	if dc.WORD == "quisque" || dc.WORD == "quispiam" {
		ret[0][0] = strings.Replace(ret[0][0], "qui", "quis", -1)
		ret[3][1] = strings.Replace(ret[3][0], "qui", "quis", -1)
	}
	return ret
}

func tableCommon(tname string, keys []string) [][]string {
	return tableCommon2(tname, keys, "")
}

func tableCommon2(tname string, keys []string, enclitic string) [][]string {
	if table, ok := Ending.Etable[tname]; ok {
		xaxis := strings.Split(table[0], "\t")
		yaxis := len(table)
		ret := make([][]string, len(xaxis))
		for i := 0; i < len(xaxis); i++ {
			ret[i] = make([]string, yaxis)
		}
		for y := 0; y < yaxis; y++ {
			sa := strings.Split(table[y], "\t")
			for x := 0; x < len(sa); x++ {
				if len(sa[x]) == 0 || sa[x] == " " {
					continue // ブランクは読み飛ばす
				}
				saitem := strings.Split(sa[x], " ")
				if len(saitem) < 2 { // saitem : 1 aなど
					continue
				}
				keyseq, _ := strconv.Atoi(saitem[0])
				if keyseq == 0 {
					ret[x][y] = saitem[1] + enclitic
				} else if keys[keyseq] != "zzz" {
					ret[x][y] = keys[keyseq] + saitem[1] + enclitic
				}
			}
		}
		return ret
	} else {
		return nil
	}
}

// pplの２重テーブル用
func ppttablecreation(tname string, keys []string, offset int) [][]string {
	if table, ok := Ending.Etable[tname]; ok {
		ret := make([][]string, 6) // テーブルの長さに関係なく6x6
		for i := 0; i < 6; i++ {
			ret[i] = make([]string, 6)
		}
		for y := 0; y < 6; y++ {
			sa := strings.Split(table[y+offset], "\t")
			for x := 0; x < len(sa); x++ {
				if len(sa[x]) == 0 {
					continue // ブランクは読み飛ばす
				}
				itemsa := strings.Split(sa[x], " ")
				keyseq, _ := strconv.Atoi(itemsa[0])
				if keys[keyseq] != "zzz" {
					ret[x][y] = keys[keyseq] + itemsa[1]
				}
			}
		}
		return ret
	}
	return nil
}

func setupkeys(dc S.DicContent) []string {
	sa := strings.Split(dc.PARTS, " ") // N11 FT,ros,ros,,
	tname := sa[0]
	keys := strings.Split(sa[1], ",")
	if len(keys) < 5 { // qui,cuのあと,,が足りてないものもあり得る？
		for i := 0; i < 5-len(keys); i++ {
			keys = append(keys, "")
		}
	}
	for i := 1; i < 5; i++ {
		// V11の場合には、key 3, 4(完了、分詞）の最後から２番めの母音を長音かする ex. liberavi -> libera--vi
		if tname == "V11" && i >= 3 {
			newkey := []rune(keys[i])
			for n, c := range keys[i] {
				if n == len(keys[i])-2 {
					newkey[n] = macron(c)
				} else {
					newkey[n] = c
				}
			}
			keys[i] = string(newkey)
			// V61,V34のときはKey２では最後の文字、Key=3で最後から２つ目を長音化　（Coei, coeiv)
		} else if (tname == "V61" || tname == "V34") && (i == 3 || i == 2) {
			newkey := []rune(keys[i])
			for n, c := range keys[i] {
				if n == len(keys[i])-(i-1) {
					newkey[n] = macron(c)
				} else {
					newkey[n] = c
				}
			}
			keys[i] = string(newkey)
			// N31のKey 2の最後から２番めの母音がaかoのときは長音化
		} else if tname == "N31" && i == 2 {
			newkey := []rune(keys[i])
			c := newkey[len(newkey)-2]
			newkey[len(newkey)-2] = minimacron(c)
			keys[i] = string(newkey)
		}
	}
	if keys[0] == "SUPER" {
		keys[4] = keys[1] // ADJ00 SUPERとして載っているものはkey=4がなく１のみ
	}
	return keys
}

func macron(ch rune) rune {
	switch ch {
	case 'a':
		return 'ā'
	case 'e':
		return 'ē'
	case 'i':
		return 'ī'
	case 'o':
		return 'ō'
	case 'u':
		return 'ū'
	default:
		return ch
	}
}

func minimacron(ch rune) rune {
	switch ch {
	case 'a':
		return 'ā'
	case 'o':
		return 'ō'
	default:
		return ch
	}
}

func getEnclitic(dc S.DicContent) string {
	word := dc.WORD
	offset := 3
	if strings.HasPrefix(word, "quis") {
		offset = 3
	}
	return word[offset:]
}
