package main

type SmallMap struct {
}

func (smallMap SmallMap) GetTileConfig() string {
	return `
-,s,s?2,s,so3,-
s?1,t,p,m,s,-
s,h,t,m,f,sb4
sl1,h,p,f,s,-
-,s,sw0,s,sg5,-
`
}

func (smallMap SmallMap) GetChits() []int {
	return []int{9, 10, 8, 12, 11, 4, 3, 5, 2, 6}
}
