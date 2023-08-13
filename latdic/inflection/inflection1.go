package inflection

import (
	S "latindictionary/structs"
	"strings"
	"latindictionary/stemlist"
)

func GetInfllection(word string) ([]S.Findings, bool) {
	var stem string
	var ending string
	finalAnswer := []S.Findings{}
	ret := []S.StemListData{}

	for x :=0 ; x <len(word)+1;x++ {
		if x == len(word){		// 本来はesse es estなどSUMの変化でstemがe のものがあればよいが
			stem = ""			// stem eはstemlistにない。追加しようとすると、sumの合成語(interestなど)
					 // との整合が取れないので、ここでstem=""はV51が見つかったことにする。keyは２
			ending = strings.ToLower(word)
			ret = append(ret, S.StemListData{Group: "V51", Key : 2, Additional : "To_BE,,s"})
		} else {
			stem = word[0: len(word) -x]
			ending = word[len(word)-x:]		// 最後まで
			ret = stemlist.Lookup(stem)
		}
		if len(ret) > 0 {
			for _, sld :=range ret {
				finds,ok := Findmatch(sld, ending,stem)
				if ok {
					finds.Group = sld.Group		// V11,ADJ11　など
					if finds.Parts == "V" {
						nominalStem := stem
						deponentflag := false	// Depの場合、能動不定詞を除くため
						ainfo := strings.Split(sld.Additional, ",")
						if len(ainfo[2]) != 0 && ainfo[2] != "zzz" {	 // stemlistの最後の項目で,で区切られた最後の項目はその後のStemである。
							nominalStem = ainfo[2]
						}
						// basicForm再確認
						if ainfo[0] == "DEP" {
							finds.BasicForm = nominalStem + endingDEP(sld.Group)	 // Vで　DEPの場合は必ず受身形でORで終了
							finds.Parts += "Deponent"
							if !strings.HasPrefix(finds.AttrNum, "PPL") {
								deponentflag = true
							}
						} else if ainfo[0] == "IMPERS" {
							finds.BasicForm = nominalStem + endingIMPERS(sld.Group)
							finds.Note += "Impersonal"
						}  else if ainfo[0] == "PERFDEF" {
							finds.BasicForm = nominalStem + "i"	// PERFDEFは完了形のみの欠如動詞　odi,coepi等がある
							finds.Note += "Perfect Defective"
						} else {
						//上記以外の場合で動詞のときにはTRANS,INTRANSが入る（またはX,例外的にTO_BEINGなど）
						// これはTRANS,INTRANSの両方がFulldictにある場合、重複して表示されるのを防ぐため使用する。
						// その他にはDEPやDATなどもありうる		
							finds.Note = ainfo[0]
						}
						if deponentflag {	// Deponentの場合、能動不定詞を取り除く
							newattribute := []S.Cordinate{}
							for i:=0; i < len(finds.Attribute) ; i++ {
								y := finds.Attribute[i].Y
								x := finds.Attribute[i].X
								if IsPassive(x,y) {
									newattribute = append(newattribute, finds.Attribute[i])
								}
							}
							finds.Attribute = newattribute
						}
					} else if finds.Parts == "N" {
						ainfo := strings.Split(sld.Additional, ",")
						if len(ainfo) >= 2 {		// P,N,popularisのように３つのときも多い
							finds.Note = ainfo[1] + ainfo[0]		// とにかく逆にして入れておく P,N -> NP
						}
					}
					if finds.BasicForm == "quipiam" || finds.BasicForm == "quique" {
						finds.BasicForm = "quis" + finds.BasicForm[3:len(finds.BasicForm) -3]
					}
					if len(finds.Attribute) > 0 || finds.AttrNum == "4" {	//attrnum=4 数副詞
						finalAnswer = append(finalAnswer, finds)
					}
				}
			}
		}
	}
	if len(finalAnswer) >0 {
		return finalAnswer,true
	} else {
		return finalAnswer, false
	}
}


func IsPassive(x,y int16) bool {
	switch {
	case y >= 6 && y<= 8 : return true
	case y==13 || y == 14 || y==17 || y==18 : return true
	case (y ==19 || y == 20) && (x >=3 && x <= 5) : return true
	default: return false
	}
}

func endingDEP(group string) string {
	if group[1] == '2' {
		return "eor"
	} else {
		return "or"
	}
}

func endingIMPERS(group string) string {
	switch {
	case group== "V11": return "at"
	case group == "V21" : return "et"
	case group == "V32" || group == "V34" || group == "v61" || 
		group == "V71" || group == "V72" : return "t"
	case group[1] == '5' : return "est"
	default : return "it"
	}
}