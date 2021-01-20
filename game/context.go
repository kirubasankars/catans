package game

import (
	"catans/board"
	"errors"
	"time"
)

type GameContext struct {
	tiles   []Tile
	board   board.Board

	GameSetting
	GameState
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

func (gc *GameContext) handOverCards(player *Player, cardType int, count int) {
	player.Cards[cardType] = player.Cards[cardType] + count
}

func (gc *GameContext) updateGameSetting(gs GameSetting) error {
	if gc.GameState.Phase != Phase1 {
		return errors.New("invalid action")
	}
	gc.GameSetting = gs
	for i := 0; i < gs.NumberOfPlayers; i++ {
		player := NewPlayer()
		player.ID = i
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

func NewGameContext() *GameContext {
	gc := new(GameContext)
	gc.Phase = Phase1
	gc.Bank = NewBank()
	return gc
}
