package main

func main() {
	//b := board.NewBoard(0)
	//b.GetUINodes()
	//lobby := game.NewLobby()
	//
	//gs := *new(game.GameSetting)
	//gs.NumberOfPlayers = 3
	//
	//gameId, _ := lobby.CreateGame(gs)
	//game := lobby.GetGame(gameId)
	//game.Start()
	//
	//game.RollDice()
	//
	//webserver.StartWebServer()

	//m := `
	//		-,-,t,t,t,-,-
	//	 	 -,t,t,t,t,-
    //        t,t,t,t,t
	//		 -,t,t,t,t,-
	//	    -,-,t,t,t,-,-
	//	 `
	//rx, _ := regexp.Compile("[[:blank:]]")
	//o := rx.ReplaceAll([]byte(m), []byte(""))
	//
	//var output [][]string
	//segments := strings.Split(string(o), "\n")
	//i := 0
	//for _, seg := range segments {
	//	row := strings.Split(seg, ",")
	//	if len(row) > 1 {
	//		var newRow []string
	//		for _, v := range row {
	//			if len(v) == 0 {
	//				continue
	//			}
	//			if v == "-" {
	//				newRow = append(newRow, v)
	//			} else {
	//				newRow = append(newRow, strconv.Itoa(i))
	//				i++
	//			}
	//
	//		}
	//		output = append(output, newRow)
	//	}
	//}
	//
	//for i := 0; i < len(output); i++ {
	//	var (
	//		currRow []string
	//		prevRow []string
	//		nextRow []string
	//	)
	//
	//	currRow = output[i]
	//	if i != 0 {
	//		prevRow = output[i-1]
	//	}
	//	if i < len(output) {
	//		nextRow = output[i+1]
	//	}
	//
	//	for j := 0; j < len(currRow); j++ {
	//		if currRow[j] == "-" {
	//			continue
	//		}
	//
	//		var sides [][2]string
	//		addSide := func(r []string, idx int, s int) {
	//			if r != nil && len(r) >= idx && r[idx] != "-" {
	//				sides = append(sides, [2]string{strconv.Itoa(s), r[idx]})
	//			}
	//		}
	//		addSide(prevRow, j+1, 0)
	//		addSide(currRow, j+1, 1)
	//		addSide(nextRow, j, 2)
	//		addSide(nextRow, j-1, 3)
	//		addSide(currRow, j-1, 4)
	//		addSide(prevRow, j, 5)
	//
	//		fmt.Println(currRow[j], sides)
	//	}
	//}
}
