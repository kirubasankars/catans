package game

import (
	"fmt"
	"time"
)

type Game struct {
	gameContext *GameContext

	ticker     *time.Ticker
	tickerDone chan bool
}

func (game *Game) startLoop() {
	go func() {
		for {
			select {
			case <-game.tickerDone:
				return
			case t := <-game.ticker.C:
				_ = t
				game.actionLoop()
			}
		}
	}()
}

func (game *Game) stopLoop() {
	game.tickerDone <- false
	game.ticker.Stop()
}

func (game *Game) actionLoop() {
	//time out, run the next action

	context := game.gameContext
	if context.getAction() == nil || game.isActionTimeout() {
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

	err := game.endAction()
	if err != nil {
		fmt.Println(err)
	}
}

func (game *Game) Start() {
	game.startLoop()
	game.startPhase2()
}

func (game *Game) Stop() {
	game.stopLoop()
}

func (game *Game) UpdateGameSetting(gs GameSetting) error {
	return game.gameContext.updateGameSetting(gs)
}

func (game *Game) startPhase2() {
	game.gameContext.Phase = Phase2
	game.gameContext.CurrentPlayer = 0
	game.scheduleAction(ActionPlaceSettlement)
}

func (game *Game) startPhase3() {
	game.gameContext.Phase = Phase3
	game.scheduleAction(ActionPlaceSettlement)
}

func (game *Game) startPhase4() {
	game.gameContext.Phase = Phase4
	game.gameContext.CurrentPlayer = 0
	game.scheduleAction(ActionRollDice)
}

func (game *Game) phase2GetNextAction() string {
	context := game.gameContext
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

func (game *Game) phase3GetNextAction() string {
	context := game.gameContext
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

func (game *Game) endAction() error {
	context := game.gameContext
	fmt.Println("END", context.getAction().Name, context.CurrentPlayer)

	NumberOfPlayers := context.GameSetting.NumberOfPlayers - 1

	if Phase4 == context.Phase {
		lastAction := context.getAction().Name

		if lastAction == ActionDiscardCards {
			game.scheduleAction(ActionPlaceRobber)
			return nil
		}

		if lastAction == ActionPlaceRobber {
			game.scheduleAction(ActionSelectToRob)
			return nil
		}

		if lastAction == ActionSelectToRob || lastAction == ActionRollDice {
			game.scheduleAction(ActionTurn)
			return nil
		}

		if lastAction == ActionTurn {
			if context.CurrentPlayer < NumberOfPlayers {
				context.CurrentPlayer++
			} else {
				context.CurrentPlayer = 0
			}
			game.scheduleAction(ActionRollDice)
		}
		return nil
	}

	if Phase2 == context.Phase {
		nextAction := game.phase2GetNextAction()
		if nextAction == "" && context.CurrentPlayer < NumberOfPlayers {
			context.CurrentPlayer++
			nextAction = game.phase2GetNextAction()
		}
		if nextAction == "" && context.CurrentPlayer == NumberOfPlayers {
			game.startPhase3()
		} else {
			context.scheduleAction(nextAction)
		}
	}

	if Phase3 == context.Phase {
		nextAction := game.phase3GetNextAction()
		if nextAction == "" && context.CurrentPlayer > 0 {
			context.CurrentPlayer--
			nextAction = game.phase3GetNextAction()
		}
		if nextAction == "" && context.CurrentPlayer == 0 {
			game.startPhase4()
		} else {
			game.scheduleAction(nextAction)
		}
	}

	return nil
}

func (game *Game) isActionTimeout() bool {
	context := game.gameContext
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
	if int(durationSeconds) < timeout {
		return false
	}
	return true
}

func (game *Game) scheduleAction(action string) {
	game.gameContext.scheduleAction(action)
}

func (game *Game) RollDice() error {
	dice1, dice2 := RandomDice()
	sum := dice1 + dice2
	if sum == 7 {
		game.scheduleAction(ActionDiscardCards)
		return nil
	}
	return HandleDice(game.gameContext, sum)
}

func NewGame() *Game {
	game := new(Game)
	game.ticker = time.NewTicker(500 * time.Millisecond)
	game.tickerDone = make(chan bool)
	game.gameContext = NewGameContext()
	return game
}

type Action struct {
	Name    string
	Timeout time.Time
}