package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"time"
)

type GameContext struct {
	board        Board
	tradeCounter int
	trades       []*GameTrade

	GameSetting
	GameState

	Users        []*User
}

func (context GameContext) getCurrentPlayer() *Player {
	return context.Players[context.CurrentPlayerID]
}

func (context *GameContext) updateGameSetting(gs GameSetting) error {
	if context.GameState.Phase != Phase1 || gs.NumberOfPlayers <= 1 || gs.Map < 0 || gs.Map > 1 {
		return errors.New(ErrInvalidOperation)
	}
	if context.GameSetting.DiscardCardLimit < 7 {
		context.GameSetting.DiscardCardLimit = 7
	}
	context.GameSetting = gs
	for i := 0; i < gs.NumberOfPlayers; i++ {
		player := NewPlayer()
		player.ID = i
		context.Players = append(context.Players, player)
	}
	//TODO: safe resizing
	//context.Users = make([]*User, gs.NumberOfPlayers)
	context.board = NewBoard(gs.Map)
	context.Tiles = context.board.GetTiles()
	return nil
}

func (context *GameContext) isInitialSettlementDone() bool {
	settlementCount := 0
	for _, player := range context.Players {
		settlementCount = settlementCount + len(player.settlements)
	}
	return settlementCount == (context.GameSetting.NumberOfPlayers * 2)
}

func (context GameContext) getGamePhase() string {
	return context.Phase
}

func (context *GameContext) startPhase2() error {
	if context.Phase != Phase1 {
		return errors.New(ErrInvalidOperation)
	}
	context.Phase = Phase2
	context.CurrentPlayerID = 0
	context.scheduleAction(ActionPlaceSettlement)
	return nil
}

func (context *GameContext) startPhase3() error {
	if context.Phase != Phase2 {
		return errors.New(ErrInvalidOperation)
	}
	context.Phase = Phase3
	context.scheduleAction(ActionPlaceSettlement)
	return nil
}

func (context *GameContext) startPhase4() error {
	if context.Phase != Phase3 {
		return errors.New(ErrInvalidOperation)
	}
	context.Phase = Phase4
	context.giveInitialFreeCards()
	context.CurrentPlayerID = 0
	context.scheduleAction(ActionRollDice)
	return nil
}

func (context *GameContext) phase2GetNextAction() string {
	currentPlayer := context.getCurrentPlayer()
	settlementCount := len(currentPlayer.settlements)
	roadCount := len(currentPlayer.roads)

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

func (context *GameContext) phase3GetNextAction() string {
	currentPlayer := context.getCurrentPlayer()
	settlementCount := len(currentPlayer.settlements)
	roadCount := len(currentPlayer.roads)

	if settlementCount == 1 && roadCount == 1 {
		return ActionPlaceSettlement
	}
	if settlementCount == 2 && roadCount == 1 {
		return ActionPlaceRoad
	}
	return ""
}

func (context *GameContext) endAction() error {
	//fmt.Println("END", context.getActionString(), context.CurrentPlayerID)

	if Phase4 == context.Phase {
		//clean up trades
		if len(context.trades) > 0 {
			context.trades = []*GameTrade{}
		}

		switch context.getActionString() {
		case ActionDiscardCards:
			{
				context.scheduleAction(ActionPlaceRobber)
			}
		case ActionPlaceRobber:
			{
				context.scheduleAction(ActionSelectToSteal)
			}
		case ActionSelectToSteal:
			{
				context.scheduleAction(ActionTurn)
			}
		case ActionRollDice:
			{
				context.scheduleAction(ActionTurn)
			}
		case ActionTurn:
			{
				NumberOfPlayers := context.GameSetting.NumberOfPlayers - 1
				if context.CurrentPlayerID < NumberOfPlayers {
					context.CurrentPlayerID++
				} else {
					context.CurrentPlayerID = 0
				}
				context.scheduleAction(ActionRollDice)
			}
		}
		return nil
	}

	if Phase2 == context.Phase {
		NumberOfPlayers := context.GameSetting.NumberOfPlayers - 1
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
		return nil
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
		return nil
	}

	return errors.New(ErrInvalidOperation)
}

func (context GameContext) isSafeTrade(gives [][2]int, wants [][2]int) bool {
	if context.Phase != Phase4 {
		return false
	}

	glen := len(gives)
	wlen := len(wants)
	if glen == 0 || wlen == 0 || wlen > 5 || glen > 5 {
		return false
	}

	wcount := 0
	for _, w := range wants {
		wcount += w[1]
		for _, g := range gives {
			if w[0] == g[0] {
				return false // same card type can't be traded.
			}
		}
	}
	if wcount > 4 { // not allowed to trade more then 3 cards
		return false
	}

	gcount := 0
	for _, c := range gives {
		gcount += c[1]
	}
	if gcount > 4 { // not allowed to trade more then 3 cards
		return false
	}
	return true
}

func (context GameContext) isPlayerHasAllCards(playerID int, cards [][2]int) bool {
	player := context.Players[playerID]
	for _, giveCard := range cards {
		giveCardType := giveCard[0]
		giveTradeCount := giveCard[1]
		if giveTradeCount <= 0 || player.cards[giveCardType] < giveTradeCount {
			return false
		}
	}
	return true
}

func (context GameContext) publishMessage() {
	for _, user := range context.Users {
		if user.Status == 0 {
			continue
		}
		user.LastPublishedEventID = context.EventID
		user.con.SetWriteDeadline(time.Now().Add(writeWait))
		user.con.WriteMessage(websocket.TextMessage, []byte(context.Events[context.EventID]))
	}
}

func NewGameContext() *GameContext {
	gc := new(GameContext)
	gc.Phase = Phase1
	gc.Bank = NewBank()
	gc.Users = []*User{}
	return gc
}