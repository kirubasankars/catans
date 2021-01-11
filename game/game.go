package game

import (
	"time"
)

type Game struct {
	gameContext GameContext

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
	TimeoutHandler(&game.gameContext)
}

func (game *Game) Start() {
	game.startLoop()
	game.gameContext.startPhase2()
}

func (game *Game) Stop() {
	game.stopLoop()
}

func (game *Game) UpdateGameSetting(gs GameSetting) {
	game.gameContext.UpdateGameSetting(gs)
}



func NewGame() *Game {
	game := new(Game)
	game.ticker = time.NewTicker(500 * time.Millisecond)
	game.tickerDone = make(chan bool)
	game.gameContext = NewGameContext()
	return game
}
