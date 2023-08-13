package latdicmain

import (
	// Ending "latindic/endinglist"
	Full "latindictionary/fulldict"
	// Stemlist "latindic/stemlist"
	S "latindictionary/structs"
	// Infl "latindic/inflection"
	"strings"
	_"strconv"
)

	//
	// facioに対する特別処理
	// facioの受動態はfioである.　またfioはSemi－Depの自動詞でもある。
	// fioでなにか見つかると、それはfacioの 受動である可能性がある。
	// facioを通常のV31で処理すると例えばfaciorなどが受動として見つかる可能性があるが、それは除去する必要がある。
	//      inflAnswerの中にbasicFormが3文字以上で最後がfioになっているものはないか確認
	//      もしあれば、fioであればbasicForm = Facioでattributeの能動を受動にして追加
	//      もし　~fioであれば~facioが存在するかどうか確認して、存在すれば同じ処理。
	//      もしbasicFormが~facioで、その属性が受動であるものは取り除く
	//
func faciofioCheck(dicfind []S.Findings, orgAnswer []S.Findings ) []S.Findings {
	for _, f := range dicfind {
		s := f.BasicForm
		if len(s) >= 3 && strings.HasSuffix(s, "fio") {
			orgAnswer = append(orgAnswer,f)	///////???????///////
			newbasic := ""
			if len(s) > 3 {
				newbasic = s[0: len(s) -3]
				dic :=Full.Lookup(newbasic + "facio", false)
				if len(dic) == 0 {
					continue
				}
			}
			// 新しくfacioを追加
			newf := S.Findings{}
			newf.BasicForm = newbasic + "facio"
			newf.Parts = "V"
			newf.Group = "V31"
			attrNew := []S.Cordinate{}
			for _, attrorg := range f.Attribute {
				nc, ok := toPassive(attrorg)
				if ok {
					attrNew = append(attrNew, nc)
				}
			}
			newf.Attribute = attrNew
			orgAnswer = append(orgAnswer, newf)
		} else if len(s) >= 5 && strings.HasSuffix(s, "facio") {
			for _, attr := range f.Attribute {
				if isActive(attr) {
					orgAnswer = append(orgAnswer, f)	// fecereは能動不定詞であるが、同時に受動命令形でもあるので。
					break
				}
			}
		} else {
			orgAnswer= append(orgAnswer, f)
		}
	}
	return orgAnswer
}


func toPassive(cin S.Cordinate) (S.Cordinate, bool) {
	switch {
	case cin.X == 0 && cin.Y == 19 : return S.Cordinate{X:3,Y:19},true
	case cin.Y >=0 && cin.Y <= 2 : return S.Cordinate{X:cin.X, Y: cin.Y +6}, true
	case cin.Y == 8 || cin.Y == 9: return S.Cordinate{X:cin.X, Y: cin.Y + 4}, true
	case cin.Y == 15 || cin.Y == 16: return S.Cordinate{X:cin.X, Y:cin.Y + 2}, true
	default: return S.Cordinate{}, false
	}
}

func isActive(cin S.Cordinate) bool {
	switch {
	case cin.Y <= 8: return true
	case cin.Y >=9 && cin.Y <= 12: return true
	case cin.Y == 15 || cin.Y == 16: return true
	case cin.X == 0 && cin.Y ==19: return true	 // 不定詞は受動命令形
	default: return false
	}
}
