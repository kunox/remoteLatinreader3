package endinglist

import (
	_"io/ioutil"
	"strings"
	_ "strconv"
	_"fmt"
	_"embed"
	// S "latindic/structs"
)


//go:embed endingTable
var data string
 var Etable  map[string][]string 

func Init() bool {
	// data, err :=  ioutil.ReadFile("endinglist/endingtable")
	// if err != nil {
	// 	return false
	// }
	groups := strings.Split(string(data), "#")
	if Etable == nil {
		Etable = make(map[string][]string)
	}
	for _,s := range groups {
		if len(s) > 0 {
			glines := strings.Split(s, "\r\n")
			mkey := glines[0]
			var mval []string
			for i :=1; i < len(glines)-1 ; i++ {
				if len(glines[i]) == 0 || glines[i] [0]== '/' {
					continue
				}
				// s := strings.TrimRight(glines[i], "\t ")
				mval = append(mval, glines[i])
			}
			Etable[mkey] = mval
		}
	}
	// testv, ok := etable["PRON41"]
	// if ok {
	// 	for _, s := range testv {
	// 		fmt.Println(s)
	// 	}
	// }
	return true
}