package latdicmain

import (
	"fmt"
	S "latindictionary/structs"
	"strconv"
	"strings"
)

// DicContentの中にFindings 相当が入っている
// WORD : vis
// NOTE : V62 1,0 PRES ACTIVE IND 2 S
// Mean : volo:be willing; wish
// PART : 0
// V7で変更、座標をいれる　pres active Ind 2 s => 1,0
// 複数あるときは＄で区切る

func convertToFinding( dc S.DicContent) S.Findings {
	ret := S.Findings{}
	notesa := strings.Split(dc.NOTES, " ")
	if len(notesa) < 2 {
		return ret	// 空のret
	}
	ret.Parts = notesa[0] 	// V 6 2が入る V62に変更されている
	cordinates := strings.Split(notesa[1], "$")
	for n := 0; n <len(cordinates) ; n++ {
		cord := string2cordinate(cordinates[n])
		// if !ok {
		// 	return ret
		// }
		ret.Attribute = append(ret.Attribute, cord)	
	}
	ret.Group = dc.PARTS		///////////////////////////////////C# ?
	ret.BasicForm = toBasciForm(dc.MEANINGS)
	ret.Note = ""
	return ret
}

func string2cordinate(sin string) S.Cordinate {
	sa := strings.Split(sin, ",")
	 x, _ := strconv.Atoi(sa[0])
	y, _ := strconv.Atoi(sa[1])
	return  S.Cordinate {X: int16(x), Y: int16(y)}
}

// Mean : volo:be willing; wish
func toBasciForm(meaning string) string {
	sa := strings.Split(meaning, ":")
	if sa[0][0] < 'A' {
		sa[0] = sa[0][1:]
	}
	return strings.TrimSpace(sa[0])
}


func allAttr2string(attrs []S.Cordinate, exdc S.ExdicContent) string {
	if len(attrs) > 0 {
		switch {
		case exdc.Dic.PARTS[0] == 'V' && !strings.HasPrefix(exdc.AttrNum, "PPL") : return caseV(attrs)
		case exdc.Dic.PARTS[0] == 'N' && exdc.Dic.PARTS[1] != 'U': return caseNoun(attrs)
		case strings.HasPrefix(exdc.AttrNum, "PPL") : return casePPL(attrs, exdc)
		default: return caseADJ(attrs)
		}
	}
	return ""
}

func caseV(attrs []S.Cordinate) string{
	var ret []string
	for _, cord := range attrs {
		if cord.Y <19 {
			ret = append(ret, fmt.Sprintf("(%s %s)", vindexX[cord.X], vindexY[cord.Y]))
		} else {
			ret = append(ret, over19(cord))
		}
	}
	return strings.Join(ret, " ")
}

var vindexX [6]string  = [6] string {"1人称 単数", "2人称 単数","3人称 単数", "1人称 複数", "2人称 複数", "3人称 複数"}
var vindexY [19]string = [...] string {"現在 能動態 直説法",
														"未完了 能動態 直説法",
														"未来 能動態 直説法",
														"完了 能動態 直説法",
														"過去完了 能動態 直説法",
														"未来完了 能動態 直説法",
														"現在 受動態 直説法",
														"未完了 受動態 直説法",
														"未来 受動態 直説法",
														"現在 能動形 接続法",
														"未完了 能動形 接続法",
														"完了 能動形 接続法",
														"過去完了 能動形 接続法",
														"現在 受動態 接続法",
														"未完了 受動態 接続法",
														"現在 能動 命令法",
														"未来 能動 命令法",
														"現在 受動 命令法",
														"未来 受動 命令法",
}

func over19(cord S.Cordinate) string {
	if cord.Y == 19 {
		switch cord.X {
		case 0 : return "不定詞"
		case 1 : return "完了不定詞"
		case 3 : return "受動 不定詞"
		default : return ""
		}
	}
	if cord.Y ==20 {
		switch cord.X {
		case 0 : return "現在分詞"
		case 2 : return "未来分詞 能動"
		case 4 : return "完了分詞 受動"
		case 5 : return "未来分詞 受動"
		default: return ""
		}
	}
	return ""
}

func caseNoun(attrs []S.Cordinate) string {
	ret := make([]string, len(attrs))
	for i, cord := range attrs {
		sp := "単数"
		if cord.X == 1 {
			sp = "複数"
		}
		casestring := getCasestring(cord.Y)
		ret[i] = fmt.Sprintf("(%s %s)", sp, casestring)
	}
	return strings.Join(ret, " ")
}

func getCasestring(y int16) string {
	// caseindex := []string {"主格","対格","属格","与格","奪格","呼格"}
	return caseindex[y]
}

var caseindex [7]string = [...]string {"主格","対格","属格","与格","奪格","呼格","処格"}
var indexX [6]string = [...]string {"M 単数", "F 単数", "N 単数", "M 複数","F 複数", "N 複数"}

func caseADJ(attrs []S.Cordinate) string {
	ret := make([]string, len(attrs))
	for i, cord := range attrs {
		ret[i] = fmt.Sprintf("(%s %s)", indexX[cord.X], caseindex[cord.Y])
	}
	return strings.Join(ret, " ")
}

func casePPL(attrs []S.Cordinate, exdc S.ExdicContent) string {
	parName := participles(exdc.AttrNum)
	infl := attrToStringPar(attrs)
	s1 := fmt.Sprintf("【%s %s】", parName, infl)
	s2 :=""
	if parName == "動形容詞" {
		gerundcase := gerund(attrs)
		if len(gerundcase) > 0{
			s2 = fmt.Sprintf("【%s %s】", "動名詞",gerundcase)
		}
	}
	return s1 + s2
}

func participles(ppl string) string {
	lastchar := ppl[len(ppl) -1]
	switch lastchar {
	case '0' : return "現在分詞"
	case '2' : return "未来分詞（能動）"
	case '4' : return "完了分詞（受動）"
	case '5' : return "動形容詞"
	default : return ""
	}
}

func attrToStringPar(attrs  []S.Cordinate) string {
	ret := make([]string, len(attrs))
	for i, cord := range attrs {
		if cord.Y > 5 {
			cord.Y = cord.Y -6
		}
		ret[i] = fmt.Sprintf("(%s %s)", indexX[cord.X], caseindex[cord.Y])
	}
	return strings.Join(ret, " ")	
}

func gerund(attrs []S.Cordinate) string {
	ret :=[]string {}
	for _, cord := range attrs {
		if cord.X == 2 && cord.Y >= 1 && cord.Y <= 4{
			ret = append(ret, fmt.Sprintf("(%s)",caseindex[cord.Y]))
		}
	}
	return strings.Join(ret, " ")
}