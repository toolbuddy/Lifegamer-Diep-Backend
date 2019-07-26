package main

import (
	"github.com/f26401004/Lifegamer-Diep-backend/src/core"
	"os"
	"io/ioutil"
	"path"
)

func main() {
	// remove all logs if not in the production mode
	if (os.Getenv("APP_ENV") != "production") {
		dir, _ := ioutil.ReadDir("./logs")
    for _, d := range dir {
        os.RemoveAll(path.Join([]string{"./logs", d.Name()}...))
    }
	}
	var app = core.App{}
	app.Run()
}