package fulldict

import (
	_"io/ioutil"
	"strings"
	S "latindictionary/structs"
	_"embed"
)

var Primaryword []string
var Pwordmacron [] string
var Notes [] string
var Meanings [] string
var Parts [] string

var start [26]int
var length [26]int

//go:embed newfulldict
var data string


func Init() bool {
	// data, err := ioutil.ReadFile("fulldict/newfulldict")
	// if err != nil{
	// 	return false
	// }

	datastring := string(data)
	lines := strings.Split(datastring, "\r\n")

	var prevchar byte= 'a'
	count := 1
	seq := 0
	start[0] =0

	for  i :=0; i < len(lines)-1; i++ {
		trimmed := strings.Trim(lines[i], " ")
		if trimmed[0] != '/' {
			items := strings.Split(trimmed, "#")
			if len(items) != 6 {
				return false
			} 
			macronw := strings.TrimSpace(items[5])
			if len(macronw) > 0 {
				 Pwordmacron = append(Pwordmacron, strings.Split(macronw, " ")[0])
			}else {
				Pwordmacron = append(Pwordmacron, "")
			}
			Primaryword = append(Primaryword, items[1])	// 大文字、小文字は区別して読み込み、照合時に辻褄を合わせる
			Notes = append(Notes,items[2] )
			Meanings = append(Meanings, items[3])
			Parts = append(Parts, items[4])	// V11など分類と,で区切られたstem（４個）V11,am,am,amav,amat
			
			//副詞の比較級などのための特別処理
			notesa := strings.Split(items[2], " ")
			if len(notesa) == 4 && notesa[2] == "ADV" {
				Primaryword = append(Primaryword, strings.Trim(notesa[0],","))
				Pwordmacron = append(Pwordmacron, "")
				Notes = append(Notes, " ADV COMP")
				Meanings = append(Meanings, items[1])
				Parts = append(Parts, "3")
				Primaryword = append(Primaryword,  strings.Trim(notesa[1], ","))
				Pwordmacron = append(Pwordmacron, "")
				Notes = append(Notes, " ADV SUPER")
				Meanings = append(Meanings, items[1])
				Parts = append(Parts, "3")
				seq += 3
				count += 3
				continue
			}
			t := items[1][0]
			if t < 'a' {
				t = t + 0x20
			}
			if t != byte(prevchar) {
				n := t - 'a'
				start[n] = seq
				length[n-1] = count		// 'a'の長さ、ここへ来るのは'b'から
				count = 1
				seq +=1
				prevchar = t
			} else {
				seq +=1
				count +=1
			}
		}
	}
	length[25] = count
	return true
}



func Lookup(instring string, macronflag bool) []S.DicContent {
	instr := strings.ToLower(instring)
	var ret []S.DicContent 
	if len(instr) == 0 {
		return ret		// len==0
	}
	n := instr[0] - 'a'
	index := start[n]
	lenhere := length[n]
	// count :=0
	found := false

	for i:=0; i <lenhere ; i++ {
		if strings.ToLower(Primaryword[i + index]) == instr {
			var tword string
			if !macronflag {
				tword = Primaryword[i + index]
			} else {
				tword = Pwordmacron[i + index]
			}
			dc := S.DicContent {
				WORD : tword,
				NOTES:Notes[i + index], 
				MEANINGS :Meanings[i + index],
				 PARTS :Parts[i +index],
			}
			ret = append(ret, dc)
			found = true
		} else {
			if found {
				break
			}
		}
	}
	return ret
}


