package main

import (
	"errors"
	"time"
)

type GameContext struct {
	board        Board
	tradeCounter int
	trades       []*GameTrade

	GameSetting
	GameState
}

func (context GameContext) getCurrentPlayer() *Player {
	return context.Players[context.CurrentPlayerID]
}

func (context *GameContext) updateGameSetting(gs GameSetting) error {
	if context.GameState.phase != Phase1 || gs.NumberOfPlayers <= 1 || gs.Map < 0 || gs.Map > 1 {
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
	context.board = NewBoard(gs.Map)
	context.Tiles = context.board.GetTiles()
	return nil
}

func (context *GameContext) isInitialSettlementDone() bool {
	settlementCount := 0
	for _, player := range context.Players {
		settlementCount = settlementCount + len(player.Settlements)
	}
	return settlementCount == (context.GameSetting.NumberOfPlayers * 2)
}

func (context GameContext) getGamePhase() string {
	return context.phase
}

func (context *GameContext) startPhase2() error {
	if context.phase != Phase1 {
		return errors.New(ErrInvalidOperation)
	}
	context.phase = Phase2
	context.CurrentPlayerID = 0
	context.scheduleAction(ActionPlaceSettlement)
	return nil
}

func (context *GameContext) startPhase3() error {
	if context.phase != Phase2 {
		return errors.New(ErrInvalidOperation)
	}
	context.phase = Phase3
	context.scheduleAction(ActionPlaceSettlement)
	return nil
}

func (context *GameContext) startPhase4() error {
	if context.phase != Phase3 {
		return errors.New(ErrInvalidOperation)
	}
	context.phase = Phase4
	context.giveInitialFreeCards()
	context.CurrentPlayerID = 0
	context.scheduleAction(ActionRollDice)
	return nil
}

func (context *GameContext) phase2GetNextAction() string {
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

func (context *GameContext) phase3GetNextAction() string {
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

func (context *GameContext) endAction() error {
	//fmt.Println("END", context.getActionString(), context.CurrentPlayerID)

	if Phase4 == context.phase {
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

	if Phase2 == context.phase {
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

	if Phase3 == context.phase {
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
	if context.phase != Phase4 {
		return false
	}
	gl := len(gives)
	wl := len(wants)
	if gl == 0 || wl == 0 || wl > 5 || gl > 5 {
		return false
	}
	c := 0
	for _, w := range wants {
		c += w[1]
	}
	if c > 4 { // not allowed to trade more then 3 cards
		return false
	}
	return true
}

func (context GameContext) isPlayerHasAllCards(playerID int, cards [][2]int) bool {
	player := context.Players[playerID]
	for _, giveCard := range cards {
		giveCardType := giveCard[0]
		giveTradeCount := giveCard[1]
		if giveTradeCount <= 0 || player.Cards[giveCardType] < giveTradeCount {
			return false
		}
	}
	return true
}

func NewGameContext() *GameContext {
	gc := new(GameContext)
	gc.phase = Phase1
	gc.Bank = NewBank()
	return gc
}

type GameAction struct {
	Name    string
	Timeout time.Time
}

func (context GameContext) getAction() *GameAction {
	return &context.Action
}

func (context GameContext) getActionString() string {
	return context.Action.Name
}

func (context *GameContext) scheduleAction(action string) {
	context.Action = GameAction{Name: action, Timeout: time.Now()}
}

func (context *GameContext) isActionTimeout() bool {
	action := context.Action

	timeout := 0
	if context.GameSetting.Speed == 0 {
		switch action.Name {
		case ActionTurn:
			timeout = 30
		case ActionRollDice:
			timeout = 10
		case ActionPlaceSettlement:
			if context.phase == Phase2 || context.phase == Phase3 {
				timeout = 12
			}
		case ActionPlaceRoad:
			if context.phase == Phase2 || context.phase == Phase3 {
				timeout = 15
			}
		case ActionDiscardCards:
			timeout = 15
		case ActionPlaceRobber:
			timeout = 10
		case ActionSelectToSteal:
			timeout = 10
		}
	}
	durationSeconds := time.Now().Sub(action.Timeout).Seconds()
	if int(durationSeconds) > timeout {
		return false
	}
	return true
}
