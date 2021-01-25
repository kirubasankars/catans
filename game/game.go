package game

import (
	"fmt"
	"time"
)

type Game struct {
	context *gameContext

	ticker     *time.Ticker
	tickerDone chan bool
}

func (game *Game) UpdateGameSetting(gs GameSetting) error {
	return game.context.updateGameSetting(gs)
}

func (game *Game) Start() error {
	if err := game.context.startPhase2(); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-game.tickerDone:
				return
			case t := <-game.ticker.C:
				_ = t
				game.run()
			}
		}
	}()

	return nil
}

func (game *Game) Stop() {
	game.tickerDone <- false
	game.ticker.Stop()
}

func (game *Game) RollDice() error {
	dice1, dice2 := RandomDice()
	sum := dice1 + dice2
	return game.context.handleDice(sum)
}

func (game Game) BankTrade(gives [][2]int, wants [][2]int) error {
	return game.context.bankTrade(gives, wants)
}

func (game Game) SetupTrade(gives [][2]int, wants [][2]int) error {
	return game.context.setupTrade(gives, wants)
}

func (game Game) OverrideTrade(tradeID int, gives [][2]int, wants [][2]int) error {
	return game.context.overrideTrade(tradeID, gives, wants)
}

func (game Game) AcceptTrade(tradeID int) error {
	return game.context.acceptTrade(tradeID)
}

func (game Game) RejectTrade(tradeID int) error {
	return game.context.rejectTrade(tradeID)
}

func (game Game) CompleteTrade(tradeID int) error {
	return game.context.completeTrade(tradeID)
}

func (game Game) CurrentPlayer() int {
	return game.context.CurrentPlayer
}

func (game Game) GetPossibleSettlementLocations() ([]int, error) {
	return game.context.getPossibleSettlementLocations()
}

func (game Game) GetPossibleRoads() ([][2]int, error) {
	return game.context.getPossibleRoads()
}

func (game Game) PutRoad(road [2]int) error {
	return game.context.putRoad(true, road)
}

func (game Game) PutSettlement(intersection int) error {
	return game.context.putSettlement(true, intersection)
}

func (game *Game) run() {
	//time out, run the next action

	context := game.context
	if context == nil || context.getAction() == nil || !context.isActionTimeout() {
		return
	}

	fmt.Println(context.Phase, context.Action.Name, context.getCurrentPlayer().ID)

	playerAction := context.getAction()
	if context.Phase == Phase2 || context.Phase == Phase3 {
		switch playerAction.Name {
		case ActionPlaceSettlement:
			{
				context.randomPlaceInitialSettlement()
			}
		case ActionPlaceRoad:
			{
				context.randomPlaceInitialRoad()
			}
		}
	}

	switch playerAction.Name {
	case ActionDiscardCards:
		{
			context.randomDiscardCards()
		}
	case ActionRollDice:
		{
			game.RollDice()
		}
	}
}

func NewGame() *Game {
	game := new(Game)
	game.ticker = time.NewTicker(500 * time.Millisecond)
	game.tickerDone = make(chan bool)
	game.context = NewGameContext()
	return game
}
