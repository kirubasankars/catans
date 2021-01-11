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

func (game *Game) start() {
	game.startLoop()
	game.gameContext.startPhase2()
}

func (game *Game) stop() {
	game.stopLoop()
}

func (game *Game) UpdateGameSetting(gs GameSetting) error {
	return game.gameContext.updateGameSetting(gs)
}

func NewGame() *Game {
	game := new(Game)
	game.ticker = time.NewTicker(500 * time.Millisecond)
	game.tickerDone = make(chan bool)
	game.gameContext = NewGameContext()
	return game
}
