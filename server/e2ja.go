package main

import (
	_ "embed"
	"regexp"
	"strconv"
	"strings"
)

//go:embed ejdic.txt
var ejtext string


var wordarray []string
var meanarray []string

func ejinit() {
	reg := "\r\n|\n"
	lines := regexp.MustCompile(reg).Split(ejtext,-1)
	for _, line  := range lines {
		items := strings.Split(line, "\t")
		if len(items) == 2 {
		wordarray = append(wordarray, items[0])
		meanarray = append(meanarray, items[1])
		}else {
			println(line)
		}
	}
}

func ejtranlate(text string) string {
	// sb := strings.Builder{}
	 cut := ";|:|,| "
	 text2 := strings.Trim(text,cut)
	for i, word := range wordarray {
		if word == text2 {
			mean := strings.Split(meanarray[i], "/")
			disp := ""
			for j:=0; j < len(mean);j++ {
				disp += "[" + strconv.Itoa(j+1) +"]" + mean[j] + "<br>"
			}
			return disp
		}
	}
	return ""
}