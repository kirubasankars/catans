package game

import (
	"catans/board"
	"errors"
	"time"
)

type GameContext struct {
	tiles        []Tile
	board        board.Board
	tradeCounter int
	trades       []GameTrade

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

func (gc GameContext) getActionString() string {
	return gc.Action.Name
}

func (gc *GameContext) scheduleAction(action string) {
	gc.Action = Action{Name: action, Timeout: time.Now()}
}

func (gc *GameContext) setupTrade(gives [][2]int, wants [][2]int) error {

	currentPlayer := gc.getCurrentPlayer()
	if !gc.isSafeTrade(gives, wants) || !gc.isPlayerHasAllCards(currentPlayer.ID, gives) {
		return errors.New("invalid operation")
	}

	if len(gives) == 1 && len(wants) == 1 {
		wantTradeCount := wants[0][1]
		if wantTradeCount == 1 {
			wantCardType := wants[0][0]
			giveCardType := gives[0][0]
			giveTradeCount := gives[0][1]
			if giveTradeCount > 1 {
				if currentPlayer.has21 || currentPlayer.has31 {
					if currentPlayer.has21 && currentPlayer.cards21[giveCardType] == 1 && giveTradeCount == 2 {
						if _, err := gc.Bank.Give(wantCardType, wantTradeCount); err != nil {
							currentPlayer.Cards[giveCardType] = -2
							currentPlayer.Cards[wantCardType]++
						} else {
							return errors.New("invalid operation")
						}

					} else if currentPlayer.has31 && giveTradeCount == 3 {
						if _, err := gc.Bank.Give(wantCardType, wantTradeCount); err != nil {
							currentPlayer.Cards[giveCardType] = -3
							currentPlayer.Cards[wantCardType]++
						} else {
							return errors.New("invalid operation")
						}
					}
				} else {
					if giveTradeCount == 4 {
						if _, err := gc.Bank.Give(wantCardType, wantTradeCount); err != nil {
							currentPlayer.Cards[giveCardType] = -4
							currentPlayer.Cards[wantCardType]++
						} else {
							return errors.New("invalid operation")
						}
					}
				}
			}
		}
	} else {
		for _, otherPlayer := range gc.Players {
			if otherPlayer.ID != currentPlayer.ID {

				hasAllCards := false
				if gc.isPlayerHasAllCards(otherPlayer.ID, wants) {
					hasAllCards = true
				}

				var trade GameTrade
				trade.ID = gc.tradeCounter
				trade.PlayerID = otherPlayer.ID
				trade.Gives = gives
				trade.Wants = wants
				trade.HasAllCards = hasAllCards
				gc.trades = append(gc.trades, trade)

				//race condition
				gc.tradeCounter++
			}
		}
	}

	return nil
}

func (gc *GameContext) overrideTrade(tradeID int, gives [][2]int, wants [][2]int) error {
	var trade *GameTrade
	for _, t := range gc.trades {
		if t.ID == tradeID {
			trade = &t
		}
	}
	if trade == nil || !gc.isSafeTrade(gives, wants) || !gc.isPlayerHasAllCards(trade.PlayerID, gives) {
		return errors.New("invalid operation")
	}

	currentPlayer := gc.getCurrentPlayer()

	hasAllCards := false
	if gc.isPlayerHasAllCards(currentPlayer.ID, wants) {
		hasAllCards = true
	}

	trade.PlayerID = gc.CurrentPlayer
	trade.Gives = gives
	trade.Wants = wants
	trade.HasAllCards = hasAllCards
	trade.OK = false

	return nil
}

func (gc GameContext) isSafeTrade(gives [][2]int, wants [][2]int) bool {
	if gc.Phase != Phase4 {
		return false
	}
	gl := len(gives)
	wl := len(wants)
	if gl == 0 || wl == 0 || wl > 5 || gl > 5 {
		return false
	}
	return true
}

func (gc GameContext) isPlayerHasAllCards(playerID int, cards [][2]int) bool {
	player := gc.Players[playerID]
	for _, giveCard := range cards {
		giveCardType := giveCard[0]
		giveTradeCount := giveCard[1]
		if giveTradeCount <= 0 || player.Cards[giveCardType] <= giveTradeCount {
			return false
		}
	}
	return true
}

func NewGameContext() *GameContext {
	gc := new(GameContext)
	gc.Phase = Phase1
	gc.Bank = NewBank()
	return gc
}
