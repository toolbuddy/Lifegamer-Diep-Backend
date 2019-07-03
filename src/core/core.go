package core

import (
	"fmt"
	"log"
	"github.com/f26401004/Lifegamer-Diep-backend/src/game"
	"sync"
)

type ServerStatus struct {
	PlayerNumber int
	RequestNumber int
}

/**
 * App:
 * The struct to present the app configuration.
 *
 * @property {*Configuration} Configuration 								- the configuration struct of the app
 * @property {[]*game.Game} Games														- the slice of the games of the app
 * @property {chan *game.Game} CreateChannel								- the channel of create game
 */
type App struct {
	Configuration *Configuration
	Status ServerStatus
	Games []*game.Game
	ControlLock sync.Mutex
}

/**
 * <*App>.Run:
 * The function in App to run the app.
 *
 * @return {nil}
 */
func (app *App) Run() {
	fmt.Println("core run...")
	app.Configuration = &Configuration{}
	err := app.Configuration.loadFromFile()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	app.ControlLock.Lock()
	app.Games = append(app.Games, game.NewGame("playground", 8192, 8192) )
	app.ControlLock.Unlock()
	app.runServer()
}
