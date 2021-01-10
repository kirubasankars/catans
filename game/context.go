package game

import (
	"catans/board"
	"fmt"
	"time"
)

type GameContext struct {
	setting GameSetting
	tiles   []board.Tile
	board   board.Board

	GameState
}

func (gc *GameContext) StartPhase2() {
	gc.Phase = Phase2
	gc.CurrentPlayer = 0
	gc.ScheduleAction(ActionPlaceSettlement)
}

func (gc *GameContext) StartPhase3() {
	gc.Phase = Phase3
	gc.ScheduleAction(ActionPlaceSettlement)
}

func (gc *GameContext) StartPhase4() {
	gc.Phase = Phase4
	gc.CurrentPlayer = 0
	gc.ScheduleAction(ActionRollDice)
}

func (gc GameContext) GetCurrentPlayer() *Player {
	return gc.Players[gc.CurrentPlayer]
}

func (gc GameContext) GetGamePhase() string {
	return gc.Phase
}

func (gc *GameContext) PutSettlement(settlement Settlement) {
	gc.Settlements = append(gc.Settlements, settlement)
}

func (gc *GameContext) PutRoad(road Road) {
	gc.Roads = append(gc.Roads, road)
}

func (gc *GameContext) HandOverCards(player *Player, cards []string) {
	player.Cards = append(player.Cards, cards...)
}

func (gc *GameContext) UpdateGameSetting(gs GameSetting) {
	gc.setting = gs
	for i := 0; i < gs.NumberOfPlayers; i++ {
		player := NewPlayer()
		player.Id = i
		gc.Players = append(gc.Players, player)
	}
	gc.tiles = gc.board.GenerateTiles("10MO,2PA,9FO,12FI,6HI,4PA,10HI,9FI,11FO,0DE,3FO,8MO,8FO,3MO,4FI,5PA,5HI,6FI,11PA")
	gc.board = board.NewBoard()
}

func (gc *GameContext) IsInitialSettlementDone() bool {
	settlementCount := 0
	for _, player := range gc.Players {
		for _, owner := range gc.Settlements {
			if player.Id == owner.player.Id {
				settlementCount++
			}
		}
	}
	return settlementCount == (gc.setting.NumberOfPlayers * 2)
}

func (gc GameContext) GetAction() *PlayerAction {
	return &gc.Action
}

func (gc *GameContext) ScheduleAction(action string) {
	gc.Action = PlayerAction{Name: action, Timeout: time.Now()}
}

func (gc *GameContext) EndAction() error {

	fmt.Println("END",gc.GetAction().Name, gc.CurrentPlayer)

	NumberOfPlayers := gc.setting.NumberOfPlayers - 1

	if Phase4 == gc.Phase {
		lastAction := gc.GetAction().Name

		if lastAction == ActionDiscardCards {
			gc.ScheduleAction(ActionPlaceRobber)
			return nil
		}

		if lastAction == ActionPlaceRobber {
			gc.ScheduleAction(ActionSelectToRob)
			return nil
		}

		if lastAction == ActionSelectToRob || lastAction == ActionRollDice {
			gc.ScheduleAction(ActionTurn)
			return nil
		}

		if lastAction == ActionTurn {
			if gc.CurrentPlayer < NumberOfPlayers {
				gc.CurrentPlayer++
			} else {
				gc.CurrentPlayer = 0
			}
			gc.ScheduleAction(ActionRollDice)
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
			gc.StartPhase3()
		} else {
			gc.ScheduleAction(nextAction)
		}
	}

	if Phase3 == gc.Phase {
		nextAction := Phase3GetNextAction(*gc)
		if nextAction == "" && gc.CurrentPlayer > 0 {
			gc.CurrentPlayer--
			nextAction = Phase3GetNextAction(*gc)
		}
		if nextAction == "" && gc.CurrentPlayer == 0 {
			gc.StartPhase4()
		} else {
			gc.ScheduleAction(nextAction)
		}
	}

	return nil
}

func NewGameContext() GameContext {
	gc := new(GameContext)
	gc.Phase = Phase1
	return *gc
}
