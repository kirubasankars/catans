package main

type DefaultMap struct {
}

func (defaultMap DefaultMap) GetTileConfig() string {
	return `
-,s,s?5,s,so4,-
s?0,t,p,m,s,-
s,h,t,m,f,sb3
sl0,h,p,f,s,-
-,s,sw1,s,sg2,-
	`
}

func (defaultMap DefaultMap) GetChits() []int {
	return []int{10, 2, 9, 12, 6, 4, 10, 9, 11, 3, 8, 8, 3, 4, 5, 5, 6, 11}
}