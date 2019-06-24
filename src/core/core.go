package core

import (
	"fmt"
	"log"
	"github.com/f26401004/Lifegamer-Diep-backend/src/game"
)

// define the App struct
type App struct {
	Configuration *Configuration
	Games []*game.Game
}
// define the Run function in App
func (app App) Run() {
	fmt.Println("core run...")
	app.Configuration = &Configuration{}
	err := app.Configuration.loadFromFile()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	app.Games = append(app.Games, game.NewGame("main_game") )
	app.runServer()
}
