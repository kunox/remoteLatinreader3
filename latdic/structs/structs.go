package structs

import (
	_ "fmt"
)

type DicContent struct {
	WORD, NOTES, MEANINGS, PARTS string
}

type ExdicContent struct {
	Dic DicContent
	Attr string
	AttrNum string
}

// func (r DicContent) String() string {
// 	return fmt.Sprintf()
// }

type StemListData struct {
	Group string;
	Key int16;
	Additional string;
}


type Cordinate struct {
	X int16
	Y int16
}

func (c Cordinate) IsEqual(a Cordinate) bool {
	if c.X == a.X && c.Y == a.Y {
		return true
	}else {
		return false
	}
}



type Findings struct {
	Parts string
	BasicForm string
	Note string
	Attribute []Cordinate
	Group string
	AttrNum string
}