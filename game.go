package main

import (
	"fmt"
	"strings"
	"time"
)

type Game struct {
	context *GameContext

	ticker     *time.Ticker
	tickerDone chan bool
}

func (game *Game) UpdateGameSetting(gs GameSetting) error {
	return game.context.updateGameSetting(gs)
}

func (game *Game) getPlayer(playerID int) Player {
	return *game.context.Players[playerID]
}

func (game *Game) Start() error {
	if err := game.context.startPhase2(); err != nil {
		return err
	}
	if !game.context.TurnTimeOut {
		return nil
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

func (game *Game) UI() string {
	board := game.context.board
	tiles := game.context.Tiles
	var nodes []string
	for idx, h := range board.grid.nodes {
		if h.resource == "-" || h.resource == "s" {
			continue
		}
		nodes = append(nodes, fmt.Sprintf("{id:%d,x:%f,y:%f,r:%f,rs:'%s'}", h.index, h.x, h.y, h.r, convertCardTypeToTerrain(tiles[idx][0])))
	}
	var intersections []string
	for _, ins := range board.grid.intersections {
		var rs = ""
		if ins.hasPort {
			for _, n := range ins.nodes {
				if n.port {
					rs = convertCardTypeToTerrain(tiles[n.index][0])
					break
				}
			}
		}
		intersections = append(intersections, fmt.Sprintf("{id:%d,x:%f,y:%f,r:%f,port:%t,rs:'%s'}", ins.index, ins.x, ins.y, ins.r,ins.hasPort, rs))
	}
	return "{'nodes': [" + strings.Join(nodes, ",") + "], 'intersections':["+ strings.Join(intersections, ",")+"]}"
}

func (game Game) BankTrade(gives [][2]int, wants [][2]int) error {
	return game.context.bankTrade(gives, wants)
}

func (game Game) SetupTrade(gives [][2]int, wants [][2]int) error {
	return game.context.setupTrade(gives, wants)
}

func (game Game) OverrideTrade(playerID, tradeID int, gives [][2]int, wants [][2]int) error {
	return game.context.overrideTrade(playerID, tradeID, gives, wants)
}

func (game Game) AcceptTrade(playerID, tradeID int) error {
	return game.context.acceptTrade(playerID, tradeID)
}

func (game Game) RejectTrade(playerID, tradeID int) error {
	return game.context.rejectTrade(playerID, tradeID)
}

func (game Game) CompleteTrade(tradeID int) error {
	return game.context.completeTrade(tradeID)
}

func (game Game) CurrentPlayer() int {
	return game.context.CurrentPlayerID
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

func (game Game) UpgradeSettlement(intersection int) error {
	return game.context.upgradeSettlement(intersection)
}

func (game Game) BuyDevelopmentCard() error {
	return game.context.buyDevelopmentCard()
}

func (game *Game) run() {
	//time out, run the next Action

	context := game.context
	if context == nil || context.getAction() == nil || !context.isActionTimeout() {
		return
	}

	fmt.Println(context.phase, context.Action.Name, context.getCurrentPlayer().ID)

	playerAction := context.getAction()
	if context.phase == Phase2 || context.phase == Phase3 {
		switch playerAction.Name {
		case ActionPlaceSettlement:
			{
				_ = context.randomPlaceInitialSettlement()
			}
		case ActionPlaceRoad:
			{
				_ = context.randomPlaceInitialRoad()
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
			_ = game.RollDice()
		}
	case ActionPlaceRobber:
		{
			context.randomPlaceRobber()
		}
	case ActionSelectToSteal:
		{
			context.randomSelectPlayerToSteal()
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
