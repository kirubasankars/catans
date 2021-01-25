package game

import (
	"regexp"
	"strconv"
	"strings"
)

type tile struct {
	Token   int
	Terrain int
}

func convertTerrainToCardType(terrain string) int {
	switch terrain {
	case "FO":
		return 0
	case "HI":
		return 1
	case "PA":
		return 2
	case "FI":
		return 3
	case "MO":
		return 4
	}
	return -1
}

func convertCardTypeToName(cardType int) string {
	switch cardType {
	case 0:
		return "Log"
	case 1:
		return "Brick"
	case 2:
		return "Wool"
	case 3:
		return "Grain"
	case 4:
		return "Ore"
	}
	return ""
}

func generateTiles(tileSettings string) []tile {
	r := regexp.MustCompile(`(?P<Token>\d+)(?P<Terrain>\w+)?`)
	segs := strings.Split(tileSettings, ",")
	tiles := make([]tile, len(segs))
	for idx, seg := range segs {
		rs := r.FindAllStringSubmatch(seg, -1)
		if len(rs) > 0 {
			tiles[idx].Token, _ = strconv.Atoi(rs[0][1])
			if len(rs[0]) > 1 {
				tiles[idx].Terrain = convertTerrainToCardType(rs[0][2])
			}
		}
	}
	return tiles
}
