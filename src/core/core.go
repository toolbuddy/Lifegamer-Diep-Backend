package core

import (
	"fmt"
	"log"
	"github.com/f26401004/Lifegamer-Diep-backend/src/game"
)

/**
 * App:
 * The struct to present the app configuration.
 *
 * @property {*Configuration} Configuration 								- the configuration struct of the app
 * @property {[]*game.Game} Games														- the slice of the games of the app
 */
type App struct {
	Configuration *Configuration
	Games []*game.Game
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

	app.Games = append(app.Games, game.NewGame("main_game", 8192, 8192) )
	app.runServer()
}
