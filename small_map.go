package main

type SmallMap struct {
}

func (smallMap SmallMap) GetTileConfig() string {
	return `
-,s,s?2,s,so2,-
s?0,t,p,m,s,-
s,h,t,m,f,sb3
sl0,h,p,f,s,-
-,s,sw1,s,sg2,-
	`
}

func (smallMap SmallMap) GetChits() []int {
	return []int{9, 10, 8, 12, 11, 4, 3, 5, 2, 6}
}
