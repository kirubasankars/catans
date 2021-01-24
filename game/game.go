package game

import (
	"fmt"
	"time"
)

type Game struct {
	context *GameContext

	ticker     *time.Ticker
	tickerDone chan bool
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
				RandomPlaceInitialSettlement(context)
			}
		case ActionPlaceRoad:
			{
				RandomPlaceInitialRoad(context)
			}
		}
	}

	switch playerAction.Name {
	case ActionDiscardCards:
		{
			RandomDiscardCards(context)
		}
	case ActionRollDice:
		{
			game.RollDice()
		}
	}

	err := context.endAction()
	if err != nil {
		fmt.Println(err)
	}
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

func (game *Game) UpdateGameSetting(gs GameSetting) error {
	return game.context.updateGameSetting(gs)
}

func (game *Game) RollDice() error {
	dice1, dice2 := RandomDice()
	sum := dice1 + dice2
	if sum == 7 {
		game.context.scheduleAction(ActionDiscardCards)
		return nil
	}
	return HandleDice(game.context, sum)
}

func NewGame() *Game {
	game := new(Game)
	game.ticker = time.NewTicker(500 * time.Millisecond)
	game.tickerDone = make(chan bool)
	game.context = NewGameContext()
	return game
}

