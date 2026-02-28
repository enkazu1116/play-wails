package main

import (
	"context"
	"embed"
	"log"
	"play-wails/internal/db"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

var assets embed.FS

func main() {

	// TursoDBを起動
	db, err := db.NewTursoDB()
	if err != nil {
		log.Fatal("【ERROR】TursoDBの起動に失敗しました")
		return
	}

	app := NewApp(db)

	err = wails.Run(&options.App{
		Title:  "ToDo App",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind: []interface{}{
			app,
		},
		OnShutdown: func(ctx context.Context) {
			err := db.Close()
			if err != nil {
				log.Fatal("【ERROR】TursoDBのクローズに失敗しました")
				return
			}
		},
		OnStartup: app.startup,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
