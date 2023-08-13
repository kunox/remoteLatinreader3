package stemlist

import (
	_"io/ioutil"
	"strings"
	"strconv"
	S "latindictionary/structs"
	_"embed"
)


var Stems [] string
var Keys [] int16
var Group [] string
var Info [] string

var start [26]int
var length [26]int

//go:embed stemlist
var data string

func Init() bool {
	// data, err := ioutil.ReadFile("stemlist/stemlist")
	// if err != nil {
	// 	return false
	// }

	lines := strings.Split(string(data),"\r\n")
	var prevchar byte= 'a'
	count := 1
	seq := 0
	start[0] =0

	for  i :=0; i < len(lines)-1; i++ {
		trimmed := strings.Trim(lines[i], " ")
		if trimmed[0] != '/' {
			items := strings.Split(trimmed, "\t")
			if len(items) < 4 {
				continue
			}
			// if items[0] == "lob"{
			// 	println(items[0])
			// }
			Stems = append(Stems, items[0])
			Group = append(Group, items[1])
			var keyn int
			var infon int
			if items[1] == "PREP" || items[1] == "ADV"{
				keyn = 3
				infon = 2
			}else {
				keyn = 2
				infon = 3
			}
			x, _ := strconv.Atoi(items[keyn])
			Keys = append(Keys, int16(x))
			Info = append(Info, items[infon])

			t := items[0][0]
			if t < 'a' {
				t = t + 0x20
			}
			if t != prevchar {
				n := t-'a'
				start[n] = seq
				length[n-1] = count	// 'a'の長さ、ここへ来るのは'b'から
				count = 1
				seq += 1
				prevchar = t
			}else {
				seq +=1
				count +=1
			}
		}		
	}
	length[25] = count
	return true
}


func Lookup(instem string) []S.StemListData {
	trimmed := strings.ToLower(instem)
	var ret  []S.StemListData

	n := trimmed[0] - 'a'
	index := start[n]
	lenhere := length[n]
	found := false

	for i:=0; i < lenhere -1; i++ {
		if trimmed == strings.ToLower(Stems[i + index]) {
			sld := S.StemListData {
				Group: Group[i+index],
				Key: Keys[i+index],
				Additional: Info[i+index],
			}
			ret = append(ret, sld )
			found= true
		} else {
			if found {
				break
			}
		}
	}
	return ret
}