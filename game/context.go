package game

// 0 - lumber
// 1 - brick
// 2 - wool
// 3 - grain
// 4 - ore

import (
	"catans/board"
	"catans/utils"
	"container/list"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type gameAction struct {
	Name    string
	Timeout time.Time
}

type gameTrade struct {
	ID          int
	From        int
	To          int
	Gives       [][2]int
	Wants       [][2]int
	HasAllCards bool
	OK          int
}

type gameContext struct {
	tiles        [][2]int
	board        board.Board
	tradeCounter int
	trades       []*gameTrade

	GameSetting
	GameState
}

func (context gameContext) getCurrentPlayer() *Player {
	return context.Players[context.CurrentPlayerID]
}

func (context gameContext) getGamePhase() string {
	return context.Phase
}

func (context *gameContext) putSettlement(validate bool, intersection int) error {
	if validate {
		availableIntersections, _ := context.getPossibleSettlementLocations()
		if !utils.Contains(availableIntersections, intersection) {
			return errors.New(utils.ErrInvalidOperation)
		}
	}
	currentPlayer := context.getCurrentPlayer()

	if Phase2 == context.Phase || Phase3 == context.Phase {
		if context.getActionString() != ActionPlaceSettlement {
			return errors.New(utils.ErrInvalidOperation)
		}
	}

	if Phase4 == context.Phase {
		cards := [][2]int{{0, 1}, {1, 1}, {2, 1}, {3, 1}}
		if !context.isPlayerHasAllCards(currentPlayer.ID, cards) {
			return errors.New(utils.ErrInvalidOperation)
		}

		context.Bank.Begin()
		for _, card := range cards {
			currentPlayer.Cards[card[0]] -= card[1]
			context.Bank.Return(card[0], card[1])
		}
		context.Bank.Commit()
	}

	var settlement Settlement
	tileIndices := context.board.GetTileIndices(intersection)
	tokens := make([]int, len(tileIndices))
	for i, idx := range tileIndices {
		tokens[i] = context.tiles[idx][1]
	}
	settlement = Settlement{Indices: tileIndices, Tokens: tokens, Intersection: intersection}
	currentPlayer.Settlements = append(currentPlayer.Settlements, settlement)

	if Phase2 == context.Phase || Phase3 == context.Phase {
		return context.endAction()
	}

	return nil
}

func(context *gameContext) upgradeSettlement(intersection int) error {
	currentPlayer := context.getCurrentPlayer()
	if Phase4 == context.Phase {
		var settlement *Settlement
		for _, s := range currentPlayer.Settlements {
			if s.Intersection == intersection {
				settlement = &s
				break
			}
		}
		if settlement == nil {
			return errors.New(utils.ErrInvalidOperation)
		}

		cards := [][2]int{{3, 2}, {4, 3}}
		if !context.isPlayerHasAllCards(currentPlayer.ID, cards) {
			return errors.New(utils.ErrInvalidOperation)
		}

		settlement.Upgraded = true
		context.Bank.Begin()
		for _, card := range cards {
			currentPlayer.Cards[card[0]] -= card[1]
			err := context.Bank.Return(card[0], card[1])
			if err != nil {
				return err
			}
		}
		context.Bank.Commit()

		return nil
	}
	return errors.New(utils.ErrInvalidOperation)
}

func (context *gameContext) buyDevelopmentCard() error {
	currentPlayer := context.getCurrentPlayer()
	if Phase4 == context.Phase {
		cards := [][2]int{{2, 1}, {3, 1}, {4, 1}}
		if !context.isPlayerHasAllCards(currentPlayer.ID, cards) {
			return errors.New(utils.ErrInvalidOperation)
		}

		context.Bank.Begin()

		for _, card := range cards {
			currentPlayer.Cards[card[0]] -= card[1]
			err := context.Bank.Return(card[0], card[1])
			if err != nil {
				return err
			}
		}
		card, err := context.Bank.BuyDevCard()
		if err != nil {
			return err
		}

		currentPlayer.DevCards = append(currentPlayer.DevCards, card)

		context.Bank.Commit()
	}
	return errors.New(utils.ErrInvalidOperation)
}

func (context *gameContext) getPossibleRoads() ([][2]int, error) {

	if Phase4 == context.Phase {
		currentPlayer := context.getCurrentPlayer()

		noOfCoords := len(currentPlayer.Settlements) + len(currentPlayer.Roads)
		var occupiedIns = make([]int, noOfCoords)
		i := 0
		for _, settlement := range currentPlayer.Settlements {
			occupiedIns[i] = settlement.Intersection
			i++
		}
		for _, road := range currentPlayer.Roads {
			if !utils.Contains(occupiedIns, road[0]) {
				occupiedIns[i] = road[0]
			}
			if !utils.Contains(occupiedIns, road[1]) {
				occupiedIns[i] = road[1]
			}
			i++
		}

		i = 0
		var roads = make([][2]int, noOfCoords*3)
		for _, ins := range occupiedIns {
			neighborIntersections := context.board.GetNeighborIntersections2(ins)
			for _, ni := range neighborIntersections {
				roads[i] = ni
				i++
			}
		}

		var allRoads [][2]int
		for _, player := range context.Players {
			allRoads = append(allRoads, player.Roads...)
		}

		var availableRoads [][2]int
		for _, newRoad := range roads {
			found := false
			for _, currentRoad := range allRoads {
				if newRoad[0] == currentRoad[0] && newRoad[1] == currentRoad[1] {
					found = true
					break
				}
			}
			if !found {
				availableRoads = append(availableRoads, newRoad)
			}
		}
		return availableRoads, nil
	}

	if Phase2 == context.Phase || Phase3 == context.Phase {
		getRoadsForIntersection := func(settlement *Settlement) [][2]int {
			var roads [][2]int
			if settlement != nil {
				neighborIntersections := context.board.GetNeighborIntersections2(settlement.Intersection)
				roads = append(roads, neighborIntersections...)
			}
			return roads
		}

		if Phase2 == context.Phase {
			if context.getActionString() != ActionPlaceRoad {
				return nil, errors.New(utils.ErrInvalidOperation)
			}
			currentPlayer := context.getCurrentPlayer()

			var firstSettlement *Settlement
			if len(currentPlayer.Settlements) > 0 {
				firstSettlement = &currentPlayer.Settlements[0]
			}
			return getRoadsForIntersection(firstSettlement), nil
		}

		if Phase3 == context.Phase {
			nextAction := context.getAction()
			if nextAction != nil && nextAction.Name != ActionPlaceRoad {
				return nil, errors.New(utils.ErrInvalidOperation)
			}

			currentPlayer := context.getCurrentPlayer()
			var (
				settlementCounter = 1
				secondSettlement  *Settlement
			)

			if len(currentPlayer.Settlements) > 1 {
				settlementCounter++
				secondSettlement = &currentPlayer.Settlements[1]
			}

			return getRoadsForIntersection(secondSettlement), nil
		}
	}

	return nil, errors.New(utils.ErrInvalidOperation)
}

func (context *gameContext) getPossibleSettlementLocations() ([]int, error) {
	if Phase4 == context.Phase {

		var occupiedIns []int
		for _, player := range context.Players {
			for _, settlement := range player.Settlements {
				occupiedIns = append(occupiedIns, settlement.Intersection)
			}
		}

		for _, v := range occupiedIns {
			neighborIntersections := context.board.GetNeighborIntersections1(v)
			for _, ni := range neighborIntersections {
				if !utils.Contains(occupiedIns, ni) {
					occupiedIns = append(occupiedIns, ni)
				}
			}
		}

		currentPlayer := context.getCurrentPlayer()
		var roadsIntersections []int
		for _, road := range currentPlayer.Roads {
			if !utils.Contains(occupiedIns, road[0]) {
				roadsIntersections = append(roadsIntersections, road[0])
			}
			if !utils.Contains(occupiedIns, road[1]) {
				roadsIntersections = append(roadsIntersections, road[1])
			}
		}

		var availableIntersections []int
		for _, ins := range roadsIntersections {
			if !utils.Contains(occupiedIns, ins) {
				availableIntersections = append(availableIntersections, ins)
			}
		}
		return availableIntersections, nil
	}

	if Phase2 == context.Phase || Phase3 == context.Phase {
		if context.getActionString() != ActionPlaceSettlement {
			return nil, errors.New(utils.ErrInvalidOperation)
		}
		occupied := make([]int, 0)
		for _, player := range context.Players {
			for _, s := range player.Settlements {
				occupied = append(occupied, s.Intersection)
			}
		}
		return context.board.GetAvailableIntersections(occupied), nil
	}

	return nil, errors.New(utils.ErrInvalidOperation)
}

func (context *gameContext) putRoad(validate bool, road [2]int) error {
	if road[0] > road[1] {
		s := road[1]
		road[1] = road[0]
		road[0] = s
	}
	if validate {
		availableRoads, _ := context.getPossibleRoads()
		matched := false
		for _, availableRoad := range availableRoads {
			if availableRoad == road {
				matched = true
				break
			}
		}
		if !matched {
			return errors.New(utils.ErrInvalidOperation)
		}
	}

	if Phase4 == context.Phase {
		currentPlayer := context.getCurrentPlayer()
		cards := [][2]int{{0, 1}, {1, 1}}
		if !context.isPlayerHasAllCards(currentPlayer.ID, cards) {
			return errors.New(utils.ErrInvalidOperation)
		}
		context.Bank.Begin()
		for _, card := range cards {
			currentPlayer.Cards[card[0]] -= card[1]
			context.Bank.Give(card[0], card[1])
		}
		context.Bank.Commit()
		currentPlayer.Roads = append(currentPlayer.Roads, road)
		return nil
	}

	if Phase2 == context.Phase || Phase3 == context.Phase {
		if context.getActionString() != ActionPlaceRoad {
			return errors.New(utils.ErrInvalidOperation)
		}
		currentPlayer := context.getCurrentPlayer()
		currentPlayer.Roads = append(currentPlayer.Roads, road)
		return context.endAction()
	}
	return nil
}

func (context *gameContext) handOverCards(player *Player, cardType int, count int) {
	player.Cards[cardType] = player.Cards[cardType] + count
}

func (context *gameContext) updateGameSetting(gs GameSetting) error {
	if context.GameState.Phase != Phase1 || gs.NumberOfPlayers <= 1 || gs.Map < 0 || gs.Map >= 1 {
		return errors.New(utils.ErrInvalidOperation)
	}
	context.GameSetting = gs
	for i := 0; i < gs.NumberOfPlayers; i++ {
		player := NewPlayer()
		player.ID = i
		context.Players = append(context.Players, player)
	}
	context.board = board.NewBoard(gs.Map)
	context.tiles = context.board.GetTiles()
	return nil
}

func (context *gameContext) isInitialSettlementDone() bool {
	settlementCount := 0
	for _, player := range context.Players {
		settlementCount = settlementCount + len(player.Settlements)
	}
	return settlementCount == (context.GameSetting.NumberOfPlayers * 2)
}

func (context gameContext) getAction() *gameAction {
	return &context.Action
}

func (context gameContext) getActionString() string {
	return context.Action.Name
}

func (context *gameContext) scheduleAction(action string) {
	context.Action = gameAction{Name: action, Timeout: time.Now()}
}

func (context *gameContext) bankTrade(gives, wants [][2]int) error {
	if context.Phase == Phase4 && len(gives) == 1 && len(wants) == 1 && wants[0][1] == 1 && gives[0][1] > 1 {
		currentPlayer := context.getCurrentPlayer()
		if !context.isSafeTrade(gives, wants) || !context.isPlayerHasAllCards(currentPlayer.ID, gives) {
			return errors.New(utils.ErrInvalidOperation)
		}

		bank := context.Bank
		bank.Begin()
		defer bank.Commit()

		wantCardType := wants[0][0]
		wantTradeCount := wants[0][1]
		giveCardType := gives[0][0]
		giveTradeCount := gives[0][1]

		if currentPlayer.ownPort21 || currentPlayer.ownPort31 {
			if currentPlayer.ownPort21 && currentPlayer.ports21[giveCardType] == 1 && giveTradeCount == 2 {
				if err := bank.Return(giveCardType, giveTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				if _, err := bank.Give(wantCardType, wantTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				currentPlayer.Cards[giveCardType] = -2
				currentPlayer.Cards[wantCardType]++
			} else if currentPlayer.ownPort31 && giveTradeCount == 3 {
				if err := bank.Return(giveCardType, giveTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				if _, err := bank.Give(wantCardType, wantTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				currentPlayer.Cards[giveCardType] = -3
				currentPlayer.Cards[wantCardType]++
			}
		} else {
			if giveTradeCount == 4 {
				if err := bank.Return(giveCardType, giveTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				if _, err := bank.Give(wantCardType, wantTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				currentPlayer.Cards[giveCardType] = -4
				currentPlayer.Cards[wantCardType]++
			}
		}

	}

	return errors.New(utils.ErrInvalidOperation)
}

func (context *gameContext) getTrade(tradeID int) *gameTrade {
	for _, t := range context.trades {
		if t.ID == tradeID {
			return t
		}
	}
	return nil
}

func (context *gameContext) setupTrade(gives [][2]int, wants [][2]int) error {

	currentPlayer := context.getCurrentPlayer()
	if !context.isSafeTrade(gives, wants) || !context.isPlayerHasAllCards(currentPlayer.ID, gives) {
		return errors.New(utils.ErrInvalidOperation)
	}

	for _, otherPlayer := range context.Players {
		if otherPlayer.ID != currentPlayer.ID {

			hasAllCards := false
			if context.isPlayerHasAllCards(otherPlayer.ID, wants) {
				hasAllCards = true
			}

			var trade = new(gameTrade)
			trade.ID = context.tradeCounter
			trade.From = currentPlayer.ID
			trade.To = otherPlayer.ID
			trade.Gives = gives
			trade.Wants = wants
			trade.HasAllCards = hasAllCards
			trade.OK = 0
			context.trades = append(context.trades, trade)

			//race condition
			context.tradeCounter++
		}
	}

	return nil
}

func (context *gameContext) overrideTrade(playerID, tradeID int, gives [][2]int, wants [][2]int) error {
	var trade = context.getTrade(tradeID)
	if trade == nil || !context.isSafeTrade(gives, wants) || playerID != trade.To || !context.isPlayerHasAllCards(trade.To, gives) {
		return errors.New(utils.ErrInvalidOperation)
	}

	currentPlayer := context.getCurrentPlayer()

	hasAllCards := false
	if context.isPlayerHasAllCards(currentPlayer.ID, wants) {
		hasAllCards = true
	}
	trade.From = trade.To
	trade.To = context.CurrentPlayerID
	trade.Gives = gives
	trade.Wants = wants
	trade.HasAllCards = hasAllCards
	trade.OK = 0

	return nil
}

func (context *gameContext) acceptTrade(playerID, tradeID int) error {
	trade := context.getTrade(tradeID)
	if trade == nil || trade.OK != 0 || playerID != trade.To {
		return errors.New(utils.ErrInvalidOperation)
	}
	trade.OK = 1
	return nil
}

func (context *gameContext) rejectTrade(playerID, tradeID int) error {
	trade := context.getTrade(tradeID)
	if trade == nil || trade.OK != 0 || playerID != trade.To {
		return errors.New(utils.ErrInvalidOperation)
	}
	trade.OK = -1
	return nil
}

func (context *gameContext) completeTrade(tradeID int) error {
	trade := context.getTrade(tradeID)
	if trade == nil || trade.OK != 1 {
		return errors.New(utils.ErrInvalidOperation)
	}

	trade.OK = 2

	player1 := context.Players[trade.From]
	player2 := context.Players[trade.To]

	for _, card := range trade.Gives {
		player1.Cards[card[0]] -= card[1]
		player2.Cards[card[0]] += card[1]
	}

	for _, card := range trade.Wants {
		player1.Cards[card[0]] += card[1]
		player2.Cards[card[0]] -= card[1]
	}

	return nil
}

func (context gameContext) isSafeTrade(gives [][2]int, wants [][2]int) bool {
	if context.Phase != Phase4 {
		return false
	}
	gl := len(gives)
	wl := len(wants)
	if gl == 0 || wl == 0 || wl > 5 || gl > 5 {
		return false
	}
	return true
}

func (context gameContext) isPlayerHasAllCards(playerID int, cards [][2]int) bool {
	player := context.Players[playerID]
	for _, giveCard := range cards {
		giveCardType := giveCard[0]
		giveTradeCount := giveCard[1]
		if giveTradeCount <= 0 || player.Cards[giveCardType] <= giveTradeCount {
			return false
		}
	}
	return true
}

func (context *gameContext) startPhase2() error {
	if context.Phase != Phase1 {
		return errors.New(utils.ErrInvalidOperation)
	}
	context.Phase = Phase2
	context.CurrentPlayerID = 0
	context.scheduleAction(ActionPlaceSettlement)
	return nil
}

func (context *gameContext) startPhase3() {
	context.Phase = Phase3
	context.scheduleAction(ActionPlaceSettlement)
}

func (context *gameContext) startPhase4() {
	context.Phase = Phase4
	context.CurrentPlayerID = 0
	context.scheduleAction(ActionRollDice)
}

func (context *gameContext) phase2GetNextAction() string {
	currentPlayer := context.getCurrentPlayer()
	settlementCount := len(currentPlayer.Settlements)
	roadCount := len(currentPlayer.Roads)

	if settlementCount == 0 && roadCount == 0 {
		return ActionPlaceSettlement
	}
	if settlementCount == 1 && roadCount == 0 {
		return ActionPlaceRoad
	}
	if settlementCount == 1 && roadCount == 1 {
		return ""
	}
	return ""
}

func (context *gameContext) phase3GetNextAction() string {
	currentPlayer := context.getCurrentPlayer()
	settlementCount := len(currentPlayer.Settlements)
	roadCount := len(currentPlayer.Roads)

	if settlementCount == 1 && roadCount == 1 {
		return ActionPlaceSettlement
	}
	if settlementCount == 2 && roadCount == 1 {
		return ActionPlaceRoad
	}
	return ""
}

func (context *gameContext) endAction() error {
	fmt.Println("END", context.getActionString(), context.CurrentPlayerID)

	NumberOfPlayers := context.GameSetting.NumberOfPlayers - 1

	if Phase4 == context.Phase {
		//clean up trades
		context.trades = []*gameTrade{}

		lastAction := context.getActionString()

		if lastAction == ActionDiscardCards {
			context.scheduleAction(ActionPlaceRobber)
			return nil
		}

		if lastAction == ActionPlaceRobber {
			context.scheduleAction(ActionSelectToRob)
			return nil
		}

		if lastAction == ActionSelectToRob || lastAction == ActionRollDice {
			context.scheduleAction(ActionTurn)
			return nil
		}

		if lastAction == ActionTurn {
			if context.CurrentPlayerID < NumberOfPlayers {
				context.CurrentPlayerID++
			} else {
				context.CurrentPlayerID = 0
			}
			context.scheduleAction(ActionRollDice)
		}
		return nil
	}

	if Phase2 == context.Phase {
		nextAction := context.phase2GetNextAction()
		if nextAction == "" && context.CurrentPlayerID < NumberOfPlayers {
			context.CurrentPlayerID++
			nextAction = context.phase2GetNextAction()
		}
		if nextAction == "" && context.CurrentPlayerID == NumberOfPlayers {
			context.startPhase3()
		} else {
			context.scheduleAction(nextAction)
		}
	}

	if Phase3 == context.Phase {
		nextAction := context.phase3GetNextAction()
		if nextAction == "" && context.CurrentPlayerID > 0 {
			context.CurrentPlayerID--
			nextAction = context.phase3GetNextAction()
		}
		if nextAction == "" && context.CurrentPlayerID == 0 {
			context.startPhase4()
		} else {
			context.scheduleAction(nextAction)
		}
	}
	return nil
}

func (context *gameContext) isActionTimeout() bool {
	action := context.Action

	timeout := 0
	if context.GameSetting.Speed == 0 {
		switch action.Name {
		case ActionTurn:
			timeout = 30
		case ActionRollDice:
			timeout = 10
		case ActionPlaceSettlement:
			if context.Phase == Phase2 || context.Phase == Phase3 {
				timeout = 12
			}
		case ActionPlaceRoad:
			if context.Phase == Phase2 || context.Phase == Phase3 {
				timeout = 15
			}
		case ActionDiscardCards:
			timeout = 15
		case ActionPlaceRobber:
			timeout = 10
		case ActionSelectToRob:
			timeout = 10
		}
	}
	durationSeconds := time.Now().Sub(action.Timeout).Seconds()
	if int(durationSeconds) > timeout {
		return false
	}
	return true
}

func (context *gameContext) handleDice(dice int) error {
	if dice == 7 {
		for _, player := range context.Players {
			if yes, _ := player.hasMoreCardsThen(context.DiscardCardLimit); yes {
				context.scheduleAction(ActionDiscardCards)
				return nil
			}
		}
		context.scheduleAction(ActionPlaceRobber)
		return nil
	}

	bank := context.Bank
	players := context.Players

	bank.Begin()
	defer bank.Commit()

	var cards [][3]int
	for _, player := range players {
		for _, settlement := range player.Settlements {
			for idx, token := range settlement.Tokens {
				if token == dice {
					tileIndex := settlement.Indices[idx]
					if tileIndex == context.RobberPlacement {
						continue
					}
					terrain := context.tiles[settlement.Indices[idx]][0]
					var count = 0
					var err error = nil
					if settlement.Upgraded {
						count, err = context.Bank.Give(terrain, 2)
					} else {
						count, err = context.Bank.Give(terrain, 1)
					}
					if err != nil {
						bank.Rollback()
						return err
					}
					cards = append(cards, [3]int{player.ID, terrain, count})
				}
			}
		}
	}

	for _, card := range cards {
		context.handOverCards(context.Players[card[0]], card[1], card[2])
	}

	return nil
}

func (context *gameContext) randomPlaceInitialSettlement() error {
	availableIntersections, _ := context.getPossibleSettlementLocations()
	selectedIntersection := availableIntersections[rand.Intn(len(availableIntersections))]
	return context.putSettlement(false, selectedIntersection)
}

func (context *gameContext) randomPlaceInitialRoad() error {
	availableRoads, _ := context.getPossibleRoads()
	selectedRoad := availableRoads[rand.Intn(len(availableRoads))]
	return context.putRoad(false, selectedRoad)
}

func (context gameContext) calculateLongestRoad(player Player, otherPlayersSettlements []int) int {
	roadNodes := player.uniqueRoadNodes()
	pending := list.New()
	longest := 0

	for _, node := range roadNodes {

		pending.PushBack(path{intersection: node, length: 0, visited: [][2]int{}})

		fmt.Println("FROM", node)

		for pending.Len() > 0 {

			var pathEnd = true

			el := pending.Front()
			pending.Remove(el)

			item := el.Value.(path)

			for _, road := range player.Roads {

				if road[0] == item.intersection || road[1] == item.intersection {
					////broken road check
					//if otherPlayersSettlements != nil && utils.Contains(otherPlayersSettlements, r1) {
					//	pathEnd = true
					//}

					p := -1
					if road[0] == item.intersection {
						p = road[1]
					} else {
						p = road[0]
					}

					visited := false
					for _, v := range item.visited {
						if v[0] == road[0] && v[1] == road[1] {
							visited = true
						}
					}
					if !visited {
						pathEnd = false
						item.visited = append(item.visited, road)
						pending.PushBack(path{intersection: p, length: item.length + 1, visited: item.visited})
					}

				}
			}

			if pathEnd {
				fmt.Println("Path End", item)
				if longest < item.length {
					longest = item.length
				}
			}
		}
	}

	return longest
}

func (context *gameContext) randomDiscardCards() {
	for _, player := range context.Players {
		if yes, cardCount := player.hasMoreCardsThen(context.DiscardCardLimit); yes {
			numberOfCardsRemove := cardCount / 2

			for {
				cardId := rand.Intn(5)
				card := player.Cards[cardId]
				if card == 0 {
					if numberOfCardsRemove <= 0 {
						break
					}
					continue
				}

				r := rand.Intn(card) + 1
				if numberOfCardsRemove-r >= 0 {
					player.Cards[cardId] = card - r
					numberOfCardsRemove -= r
				} else {
					player.Cards[cardId] = card - numberOfCardsRemove
					numberOfCardsRemove = 0
				}

				if numberOfCardsRemove <= 0 {
					break
				}
			}
		}
	}
}

func (context *gameContext) randomPlaceRobber() {
	var occupiedIns []int
	for _, player := range context.Players {
		if player.ID == context.CurrentPlayerID {
			continue
		}
		for _, settlement := range player.Settlements {
			occupiedIns = append(occupiedIns, settlement.Intersection)
		}
	}
	ins := rand.Intn(len(occupiedIns))
	tileIndices := context.board.GetTileIndices(ins)
	context.RobberPlacement = tileIndices[0]
	context.scheduleAction(ActionSelectToRob)
}

func (context *gameContext) randomSelectPlayerToRob() {
	var playerToRob int = -1
	for _, player := range context.Players {
		if player.ID == context.CurrentPlayerID {
			continue
		}
		for _, settlement := range player.Settlements {
			if settlement.Intersection == context.RobberPlacement {
				playerToRob = player.ID
			}
		}
	}

	if playerToRob > -1 {
		context.robAPlayer(playerToRob)
	}
}

func (context *gameContext) robAPlayer(playerID int) {
	currentPlayer := context.getCurrentPlayer()
	otherPlayer := context.Players[playerID]
	var availableCards []int
	for idx, card := range otherPlayer.Cards {
		if card == 0 {
			continue
		}
		availableCards = append(availableCards, idx)
	}
	r := rand.Intn(len(availableCards))
	randCardType := availableCards[r]

	otherPlayer.Cards[randCardType] -= 1
	currentPlayer.Cards[randCardType] += 1
}

func NewGameContext() *gameContext {
	gc := new(gameContext)
	gc.Phase = Phase1
	gc.Bank = NewBank()
	return gc
}
