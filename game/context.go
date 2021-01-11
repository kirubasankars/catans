package game

import (
	"catans/board"
	"errors"
	"fmt"
	"time"
)

type GameContext struct {
	tiles   []Tile
	board   board.Board

	GameSetting
	GameState
}

func (gc *GameContext) startPhase2() {
	gc.Phase = Phase2
	gc.CurrentPlayer = 0
	gc.scheduleAction(ActionPlaceSettlement)
}

func (gc *GameContext) startPhase3() {
	gc.Phase = Phase3
	gc.scheduleAction(ActionPlaceSettlement)
}

func (gc *GameContext) startPhase4() {
	gc.Phase = Phase4
	gc.CurrentPlayer = 0
	gc.scheduleAction(ActionRollDice)
}

func (gc GameContext) getCurrentPlayer() *Player {
	return gc.Players[gc.CurrentPlayer]
}

func (gc GameContext) getGamePhase() string {
	return gc.Phase
}

func (gc *GameContext) putSettlement(settlement Settlement) {
	gc.getCurrentPlayer().putSettlement(settlement)
}

func (gc *GameContext) putRoad(road [2]int) {
	gc.getCurrentPlayer().putRoad(road)
}

func (gc *GameContext) HandOverCards(player *Player, cardType int, count int) {
	player.Cards[cardType] = player.Cards[cardType] + count
}

func (gc *GameContext) updateGameSetting(gs GameSetting) error {
	if gc.GameState.Phase != Phase1 {
		return errors.New("invalid action")
	}
	gc.GameSetting = gs
	for i := 0; i < gs.NumberOfPlayers; i++ {
		player := NewPlayer()
		player.Id = i
		gc.Players = append(gc.Players, player)
	}
	gc.tiles = GenerateTiles("10MO,2PA,9FO,12FI,6HI,4PA,10HI,9FI,11FO,0DE,3FO,8MO,8FO,3MO,4FI,5PA,5HI,6FI,11PA")
	gc.board = board.NewBoard(gs.Map)
	return nil
}

func (gc *GameContext) isInitialSettlementDone() bool {
	settlementCount := 0
	for _, player := range gc.Players {
		settlementCount = settlementCount + len(player.Settlements)
	}
	return settlementCount == (gc.GameSetting.NumberOfPlayers * 2)
}

func (gc GameContext) getAction() *Action {
	return &gc.Action
}

func (gc *GameContext) scheduleAction(action string) {
	gc.Action = Action{Name: action, Timeout: time.Now()}
}

func (gc *GameContext) endAction() error {

	fmt.Println("END",gc.getAction().Name, gc.CurrentPlayer)

	NumberOfPlayers := gc.GameSetting.NumberOfPlayers - 1

	if Phase4 == gc.Phase {
		lastAction := gc.getAction().Name

		if lastAction == ActionDiscardCards {
			gc.scheduleAction(ActionPlaceRobber)
			return nil
		}

		if lastAction == ActionPlaceRobber {
			gc.scheduleAction(ActionSelectToRob)
			return nil
		}

		if lastAction == ActionSelectToRob || lastAction == ActionRollDice {
			gc.scheduleAction(ActionTurn)
			return nil
		}

		if lastAction == ActionTurn {
			if gc.CurrentPlayer < NumberOfPlayers {
				gc.CurrentPlayer++
			} else {
				gc.CurrentPlayer = 0
			}
			gc.scheduleAction(ActionRollDice)
		}
		return nil
	}

	if Phase2 == gc.Phase {
		nextAction := Phase2GetNextAction(*gc)
		if nextAction == "" && gc.CurrentPlayer < NumberOfPlayers {
			gc.CurrentPlayer++
			nextAction = Phase2GetNextAction(*gc)
		}
		if nextAction == "" && gc.CurrentPlayer == NumberOfPlayers {
			gc.startPhase3()
		} else {
			gc.scheduleAction(nextAction)
		}
	}

	if Phase3 == gc.Phase {
		nextAction := Phase3GetNextAction(*gc)
		if nextAction == "" && gc.CurrentPlayer > 0 {
			gc.CurrentPlayer--
			nextAction = Phase3GetNextAction(*gc)
		}
		if nextAction == "" && gc.CurrentPlayer == 0 {
			gc.startPhase4()
		} else {
			gc.scheduleAction(nextAction)
		}
	}

	return nil
}

func NewGameContext() GameContext {
	gc := new(GameContext)
	gc.Phase = Phase1
	gc.Bank = NewBank()
	return *gc
}
