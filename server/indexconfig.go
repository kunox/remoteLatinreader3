package main

type indexstruct struct{
	number int
	bookname string
	extension string
	vols []volstruct
}

type volstruct struct {
	name string
	chapsec string
	number int
}