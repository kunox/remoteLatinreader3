package latdicmain

import (
	Ending "latindictionary/endinglist"
	Full "latindictionary/fulldict"
	Stemlist "latindictionary/stemlist"
	S "latindictionary/structs"
	Infl "latindictionary/inflection"
	TC "latindictionary/tablecreationF"
	"strings"
	_"strconv"

)


func Init() bool {
	if	Full.Init() && Stemlist.Init() && Ending.Init() {
		return true
	} else {
		return false
	}
}

func Findword(word string) ([]S.ExdicContent, bool){
	ret := []S.ExdicContent {}
	dcAnswer := []S.DicContent{}
	inflAnswer := []S.Findings{}

    // 最初にFulldictを検索。変化しない語、例外語はここに出る。
	wordx := parseString(word)
	if len(wordx) == 0 {
		return ret, false
	}
	temp := Full.Lookup(wordx,false)
	if len(temp) > 0 {
		for _,dcnow := range temp{
			if  dcnow.PARTS == "0" || dcnow.PARTS== "1" {
				inflAnswer = append(inflAnswer, convertToFinding(dcnow))
			} else {
				dcAnswer = append(dcAnswer, dcnow)
			}
		}
	}
	// すべての語で曲用を探す。
	// FIndingsのbasicFormとAttribute(格座標）がセットされる
	//変化しない語（副詞など）、また最後のPartsが０，１の例外の語は、falseが帰る
	finds, ok := Infl.GetInfllection(wordx)
	if ok {
		inflAnswer = faciofioCheck(finds, inflAnswer)
	}
	if len(dcAnswer) == 0 && len(inflAnswer) == 0 {
		return	[]S.ExdicContent{}, false		// どちらにも出ないのは答えがない
	}
	for _ ,dc :=range dcAnswer {
		sa := strings.Split(dc.PARTS, " ")
		part := sa[0]
		if !IsDecline(part) {		 // 変化しない語はinflの方にはでない。（副詞など）
			exdc := S.ExdicContent {
				Dic : dc,
				Attr : "",
			}
			ret = append(ret,exdc)
		}
	}

	// 基数以外はfulldiscにないので、PARTSにstemを作って入れておかないと、一覧表示出来ない。
	// inflectの最後にNUM00 X,nominalbasic,,,と作成してあるのでそれをコピーすれば良い
	// またattrNumには”１”、"2","3","4"が入っていて夫々、基数、序数、配分数、数副詞と判定できる
	for _,fd := range inflAnswer {
		if strings.HasPrefix(fd.Parts, "NUM") && fd.AttrNum != "1" {
			dcNum := S.DicContent {
				WORD: fd.BasicForm,
				PARTS: fd.Parts,
				MEANINGS: fd.Note,
				NOTES: "",		/////////////////////////////////
			}
			exdicNum := S.ExdicContent {
				Dic: dcNum,
				AttrNum: fd.AttrNum,
				Attr: "",		////////////////////////////////////
			}
			ret = append(ret, exdicNum)
			continue
		}
		dcx := Full.Lookup(fd.BasicForm, true)
		if len(dcx) > 0 {
			for _,d := range dcx {
				sa := strings.Split(d.PARTS, " ")	// V11 X,am,am,amav,amat, N11 FT,ros,ros,,
				if sa[0] == fd.Group {
					exdc := S.ExdicContent{}
					exdc.Dic = d
					if fd.Parts == "V" {
						if fd.Note == "TRANS" || fd.Note == "INTRANS" {
							sa2 := strings.Split(sa[1], ",")	// V31 TRANS,--,--,--の後半を更に分解
							if sa2[0] != fd.Note {
								continue
							}
						}
						dicnote := strings.Split(exdc.Dic.NOTES, " ")
						if dicnote[len(dicnote)-2] == fd.Note {		// goの場合は、Splitで最後に""がでて一つ多い結果となる
							fd.Note = ""	// dic.NOTESの最後がDEP、DAT、Tramsなどainfoの最初と同じだとダブルので入れない
						}
						exdc.Dic.NOTES += " " + fd.Note
					} else if strings.HasPrefix(fd.Group, "PRON1") {
						exdc.Dic.NOTES += " " + strings.Split(sa[1], ",")[0]	// ADJECT, INDEFなどがfulldic内に応じて入る
					} else if strings.HasPrefix(fd.Group, "PRON2") {
						exdc.Dic.NOTES += "PRON" + strings.Split(sa[1], ",")[0]	// PRON2Xは”PRON"が入っていない
					} else if fd.Parts == "N" {
						sagen := strings.Split(sa[1], ",")	// N11 FT,pen,pen,の後半をSplit
						if fd.Note != sagen[0] {
							continue
						}
					} else {
						exdc.Dic.NOTES += " " + fd.Note
					}
					exdc.AttrNum = fd.AttrNum	// 配分数などのための数字
					exdc.Attr = allAttr2string(fd.Attribute, exdc)
					ret  = append(ret, exdc)
				}
			}
		}
	}
	if len(ret) > 0 {
		return ret,true
	} else {
		return ret, false
	}
}

func parseString(instr string) string {
	instrx := []rune(instr)
	outdata := []rune {}
	for i :=0 ; i <len(instrx) ; i++ {
		rx := toPlain(instrx[i])
		 if !alpha(rx) {
			 continue
		 } 
		 outdata = append(outdata,rx)
		 i += 1
		 for j := i ; j< len(instrx) ;j++ {
			rx = toPlain(instrx[j])
			if !alpha(rx) {
				// i = len(instrx)
				break
			} else {
				outdata = append(outdata, rx)
			}
		}
		break
	}
	return string(outdata)
}

func alpha(c rune) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
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

func IsDecline(part string) bool {
	if strings.HasPrefix(part, "CONJ") || strings.HasPrefix(part, "ADV") || 	
	strings.HasPrefix(part, "PREP") 	||	strings.HasPrefix(part, "INTERJ") || strings.HasPrefix(part, "3") {
			return false		 // "3"は追加したADV COMP、SUPERのため
	} else if strings.HasPrefix(part, "NUM") && part[3] != '1' {
		return false
	} else if strings.HasSuffix(part, "98") || strings.HasSuffix(part, "99") {
		return false
	} else {
		return  true
	}
}

func CreateTable(exdc S.ExdicContent) [][]string{
	part := exdc.Dic.PARTS
	if strings.HasPrefix(part, "V") {
		if strings.HasPrefix(exdc.AttrNum, "PPL") {
			return TC.GetPPLTable(exdc)
		}else if strings.HasSuffix(exdc.Dic.WORD, "facio") || strings.HasSuffix(exdc.Dic.WORD,"făcĭo") {
			return TC.SpecialFacio(exdc)
		} else {
			return TC.GetVtable(exdc)
		}
	}
	if strings.HasPrefix(part,"NUM"){
		if exdc.AttrNum == "1" {
			return TC.GetADJTable(exdc.Dic)	// 基数は通常の形容詞と同じ（unus, duo tresの3つだけ）
		} else if exdc.AttrNum=="4" {
			return nil			// 数副詞は変化しない
		} else if exdc.AttrNum=="" {	//NUM20の変化しない数詞はIsDeclineでattrNum=""になっている
			return nil								// 58行目あたり
		} else {
			return TC.NUM00Table(exdc)
		}
	}
	if strings.HasPrefix(part,"ADJ") {
		if len(exdc.AttrNum) == 3 {
			return TC.GetADJ00Table(exdc)
		}else {
			return TC.GetADJTable(exdc.Dic)
		}
	}
	if strings.HasPrefix(part, "N"){
		return TC.GetNtable(exdc.Dic)
	}
	if strings.HasPrefix(part, "PRON"){
		return TC.GetPRONTable(exdc.Dic)
	}
	return nil
}
