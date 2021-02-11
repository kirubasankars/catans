package main

type DefaultMap struct {
}

func (defaultMap DefaultMap) GetTileConfig() string {
	return `
-,-,sw2,s,s?3,s,-
-,s,o,w,l,sg3,-
-,s?1,g,b,w,b,s
s,g,l,d,l,o,sl4
-,s?1,l,b,g,w,s
-,s,o,g,w,s?5,-
-,-,so0,s,sb5,s,-
	`
}

func (defaultMap DefaultMap) GetChits() []int {
	return []int{10, 2, 9, 12, 6, 4, 10, 9, 11, 3, 8, 8, 3, 4, 5, 5, 6, 11}
}